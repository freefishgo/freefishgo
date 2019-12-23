package fishSession

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"strconv"
	"sync"
	"time"
)

/*Session会话管理*/
type SessionMgr struct {
	mLock        sync.RWMutex //互斥(保证线程安全)
	mMaxLifeTime int64        //垃圾回收时间

	mSessions map[string]*Session //保存session的指针[sessionID] = session
}

func (mgr *SessionMgr) GetSessionKeyValue() (string, error) {
	return mgr.NewSessionID(), nil
}

//设置session里面的值
func (mgr *SessionMgr) SetSession(sessionID string, m map[interface{}]interface{}) error {
	mgr.mLock.Lock()
	defer mgr.mLock.Unlock()
	mgr.mSessions[sessionID] = &Session{mValues: m, mLastTimeAccessed: time.Now()}
	return nil
}

func (mgr *SessionMgr) RemoveBySessionID(sessionID string) error {
	mgr.mLock.RLock()
	defer mgr.mLock.RUnlock()
	delete(mgr.mSessions, sessionID)
	return nil
}

func (mgr *SessionMgr) Init(SessionAliveTime time.Duration) error {
	mgr.mMaxLifeTime = int64(SessionAliveTime)
	return nil
}

//创建会话管理器(cookieName:在浏览器中cookie的名字;maxLifeTime:最长生命周期)
func NewSessionMgr(maxLifeTime time.Duration) *SessionMgr {
	mgr := &SessionMgr{mMaxLifeTime: int64(maxLifeTime), mSessions: make(map[string]*Session)}

	//启动定时回收
	go mgr.GC()

	return mgr
}

//得到session里面的值
func (mgr *SessionMgr) GetSession(sessionID string) (map[interface{}]interface{}, error) {
	mgr.mLock.RLock()
	defer mgr.mLock.RUnlock()

	if session, ok := mgr.mSessions[sessionID]; ok {
		session.mLastTimeAccessed = time.Now()
		return session.mValues, nil
	}

	return nil, nil
}

//得到sessionID列表
func (mgr *SessionMgr) GetSessionIDList() []string {
	mgr.mLock.RLock()
	defer mgr.mLock.RUnlock()

	sessionIDList := make([]string, 0)

	for k, _ := range mgr.mSessions {
		sessionIDList = append(sessionIDList, k)
	}

	return sessionIDList[0:len(sessionIDList)]
}

//更新最后访问时间
func (mgr *SessionMgr) GetLastAccessTime(sessionID string) time.Time {
	mgr.mLock.RLock()
	defer mgr.mLock.RUnlock()

	if session, ok := mgr.mSessions[sessionID]; ok {
		return session.mLastTimeAccessed
	}

	return time.Now()
}

//GC回收
func (mgr *SessionMgr) GC() {
	mgr.mLock.Lock()
	defer mgr.mLock.Unlock()

	for sessionID, session := range mgr.mSessions {
		//删除超过时限的session
		if session.mLastTimeAccessed.Unix()+mgr.mMaxLifeTime < time.Now().Unix() {
			delete(mgr.mSessions, sessionID)
		}
	}

	//定时回收
	time.AfterFunc(time.Duration(mgr.mMaxLifeTime)*time.Second, func() { mgr.GC() })
}

//创建唯一ID
func (mgr *SessionMgr) NewSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		nano := time.Now().UnixNano() //微秒
		return strconv.FormatInt(nano, 10)
	}
	return base64.URLEncoding.EncodeToString(b)
}

//——————————————————————————
/*会话*/
type Session struct {
	mLastTimeAccessed time.Time                   //最后访问时间
	mValues           map[interface{}]interface{} //其它对应值(保存用户所对应的一些值，比如用户权限之类)
}
