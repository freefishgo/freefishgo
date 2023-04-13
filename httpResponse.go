// Copyright 2019 freefishgo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package freefishgo

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type IResponse interface {
	// WebSocket 升级为WebSocket服务 upgrades为空时采用默认的参数 为多个时只采用第一个作为WebSocket参数
	WebSocket(upgrades ...*websocket.Upgrader) (conn *websocket.Conn, err error)
	Hijack() (net.Conn, *bufio.ReadWriter, error)
	setISession(i ISession)
	getSessionKeyValue() (err error)
	// RemoveSession 移除Session
	RemoveSession()
	// GetSession 获取指定key的session值
	GetSession(key string) (interface{}, error)
	getSession() error
	// UpdateSession 更新session值
	UpdateSession() error
	// SetSession 设置session值
	SetSession(key string, val interface{}) error
	// SetCookie 设置Cookie
	SetCookie(c *http.Cookie)
	// SetCookieUseKeyValue 设置Cookie
	SetCookieUseKeyValue(key string, val string)
	// RemoveCookieByName 通过cookie名字移除Cookie
	RemoveCookieByName(name string)
	// WriteHeader 设置响应状态值
	WriteHeader(statusCode int)
	// ReadStatusCode 读取响应状态值
	ReadStatusCode() int
	// RemoveCookie 通过cookie移除Cookie
	RemoveCookie(ck *http.Cookie)
	// Write 写入前端的数据
	Write(b []byte) (int, error)
	// GetWriteCache 获取写入前端的缓存
	GetWriteCache() []byte
	// ClearWriteCache 清除写入前端的缓存
	ClearWriteCache()
	// WriteJson 写入前端的json数据
	WriteJson(i interface{}) error
	// Redirect 重定向路径
	Redirect(redirectPath string)
	getYourself() *Response
	// GetStarted 是否已经向前端写入数据了，默认是开启的，且GetMaxResponseCacheLen()设置很小
	GetStarted() bool
	// GetIsWriteInCache 获取当前请求是否临时缓存数据进入缓存中
	GetIsWriteInCache() bool
	// SetIsWriteInCache 设置是否延迟把写入前端的数据写入前端，即使设置了延迟写入前端，
	// 但当数据长度超过了配置文件设置的 MaxResponseCacheLen是依然会写入前端，
	// 判断是否已经写入前端 调用 GetStarted()进行判断
	SetIsWriteInCache(bool)
	Header() http.Header
	// GetMaxResponseCacheLen 查看延迟写入前端的数据的最大值
	GetMaxResponseCacheLen() int
	// SetMaxResponseCacheLen 设置延迟写入前端的数据的最大值，GetIsWriteInCache()为True时生效
	SetMaxResponseCacheLen(int)
	// GetMsgData 获取传送数据
	GetMsgData() map[string]interface{}
	// SetMsgData 设置传送数据
	SetMsgData(map[string]interface{})
	// Stack 500 错误的堆栈信息,其他状态为空
	Stack() string
	// Error 500 错误的信息,其他状态为空
	Error() interface{}
	SetStack(string)
	SetError(interface{})
}

type Response struct {
	http.ResponseWriter
	req *http.Request
	// 是否调用过Write
	Started bool
	// 回复状态
	status     int
	Gzip       *gzip.Writer
	IsOpenGzip bool
	isGzip     bool
	msgData    map[string]interface{}

	isWriteInCache      bool
	writeCache          []byte
	maxResponseCacheLen int
	//Cookies []*http.Cookie
	sessionFunc        ISession
	session            map[interface{}]interface{}
	isGetSession       bool
	SessionId          string
	SessionCookieName  string
	SessionAliveTime   time.Duration
	isUpdateSessionKey bool
	sessionIsUpdate    bool

	err   interface{}
	stack string
}

// GetMsgData 获取传递的信息
func (r *Response) GetMsgData() map[string]interface{} {
	return r.msgData
}

// SetMsgData 设置传递的信心
func (r *Response) SetMsgData(data map[string]interface{}) {
	r.msgData = data
}

// SetStack 设置错误的堆栈信息
func (r *Response) SetStack(str string) {
	r.stack = str
}

// SetError 设置错误的信息
func (r *Response) SetError(err interface{}) {
	r.err = err
}

// Stack 错误的堆栈信息
func (r *Response) Stack() string {
	return r.stack
}

// Error 错误的信息
func (r *Response) Error() interface{} {
	return r.err
}

func (r *Response) SetMaxResponseCacheLen(b int) {
	r.maxResponseCacheLen = b
}

func (r *Response) GetMaxResponseCacheLen() int {
	return r.maxResponseCacheLen
}

func (r *Response) SetIsWriteInCache(b bool) {
	r.isWriteInCache = b
}

func (r *Response) GetIsWriteInCache() bool {
	return r.isWriteInCache
}

func (r *Response) GetStarted() bool {
	return r.Started
}

func (r *Response) getYourself() *Response {
	return r
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocket 升级为WebSocket服务 upgrades为空时采用默认的参数 为多个时只采用第一个作为WebSocket参数
func (r *Response) WebSocket(upgrades ...*websocket.Upgrader) (conn *websocket.Conn, err error) {
	if upgrades == nil {
		conn, err = upgrade.Upgrade(r, r.req, r.Header())
	} else {
		conn, err = upgrades[0].Upgrade(r, r.req, r.Header())
	}
	if err == nil {
		r.Started = true
	}
	return
}

func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}
func (r *Response) setISession(i ISession) {
	r.sessionFunc = i
}

func (r *Response) getSessionKeyValue() (err error) {
	if r.SessionId == "" {
		r.isUpdateSessionKey = true
		r.SessionId, err = r.sessionFunc.GetSessionKeyValue()
	}
	return
}

func (r *Response) RemoveSession() {
	r.session = nil
	r.isGetSession = false
	r.sessionIsUpdate = false
	r.sessionFunc.RemoveBySessionID(r.SessionId)
}

func (r *Response) GetSession(key string) (interface{}, error) {
	if r.SessionId == "" {
		return nil, nil
	}
	if !r.isGetSession {
		var err error
		if err = r.getSession(); err == nil {
			r.isGetSession = true
			if r.session == nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	v, _ := r.session[key]
	return v, nil
}

func (r *Response) getSession() error {
	if r.SessionId == "" {
		return errors.New("没有设置Session")
	}
	var err error
	if !r.isGetSession {
		if r.session, err = r.sessionFunc.GetSession(r.SessionId); err == nil {
			if r.session != nil {
				r.isGetSession = true
			}
		}
	}
	return err
}

func (r *Response) UpdateSession() error {
	if r.SessionId == "" {
		return nil
	}
	if r.sessionIsUpdate {
		r.sessionIsUpdate = false
		return r.sessionFunc.SetSession(r.SessionId, r.session)
	}
	return nil
}

func (r *Response) SetSession(key string, val interface{}) error {
	r.sessionIsUpdate = true
	r.getSession()
	if err := r.getSessionKeyValue(); err != nil {
		return err
	}
	if r.session == nil {
		r.session = map[interface{}]interface{}{}
	}
	r.session[key] = val
	return r.UpdateSession()
}

// SetCookie 设置Cookie
func (r *Response) SetCookie(c *http.Cookie) {
	http.SetCookie(r, c)
}

// SetCookieUseKeyValue 设置Cookie
func (r *Response) SetCookieUseKeyValue(key string, val string) {
	http.SetCookie(r.ResponseWriter, &http.Cookie{Name: key, Value: val, Path: "/"})
}

// RemoveCookieByName 通过cookie名字移除Cookie
func (r *Response) RemoveCookieByName(name string) {
	if ck, err := r.req.Cookie(name); err != http.ErrNoCookie {
		ck.Expires = time.Now()
		http.SetCookie(r, ck)
	}
}
func (r *Response) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *Response) ReadStatusCode() int {
	return r.status
}

// RemoveCookie 通过cookie移除Cookie
func (r *Response) RemoveCookie(ck *http.Cookie) {
	ck.Expires = time.Now()
	http.SetCookie(r, ck)
}

// Write 写入前端的数据
func (r *Response) Write(b []byte) (int, error) {
	if !r.Started && r.isWriteInCache && len(r.writeCache) < r.maxResponseCacheLen {
		r.writeCache = append(r.writeCache, b...)
		return len(b), nil
	}
	defer func() {
		r.Started = true
	}()
	if r.isGzip || (r.IsOpenGzip && !r.Started) {
		if !r.Started {
			if r.SessionId != "" && r.isUpdateSessionKey {
				r.SetCookieUseKeyValue(r.SessionCookieName, r.SessionId)
			}
			r.isGzip = true
			if r.ResponseWriter.Header().Get("Content-Type") == "" {
				if r.writeCache == nil {
					r.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(b))
				} else {
					r.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(r.writeCache))
				}
			}
			r.ResponseWriter.WriteHeader(r.status)
		}
		if r.Gzip == nil {
			r.Gzip = gzip.NewWriter(r.ResponseWriter)
		}
		r.Gzip.Write(r.writeCache)
		r.writeCache = nil
		return r.Gzip.Write(b)
	}
	if !r.Started {
		if r.SessionId != "" && r.isUpdateSessionKey {
			r.SetCookieUseKeyValue(r.SessionCookieName, r.SessionId)
		}
		r.ResponseWriter.WriteHeader(r.status)
	}
	r.ResponseWriter.Write(r.writeCache)
	r.writeCache = nil
	return r.ResponseWriter.Write(b)
}

// GetWriteCache 获取写入前端的缓存
func (r *Response) GetWriteCache() []byte {
	return r.writeCache
}

// ClearWriteCache 清除写入前端的缓存
func (r *Response) ClearWriteCache() {
	r.writeCache = nil
}

// WriteJson 写入前端的json数据
func (r *Response) WriteJson(i interface{}) error {
	if b, err := json.Marshal(i); err == nil {
		r.ResponseWriter.Header().Set("Content-Type", "application/json")
		_, err = r.Write(b)
		return err
	} else {
		return err
	}
}

func (r *Response) Redirect(redirectPath string) {
	http.Redirect(r, r.req, redirectPath, 302)
}
