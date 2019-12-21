package httpContext

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"regexp"
	"time"
)

type Response struct {
	http.ResponseWriter
	req *http.Request
	// 是否调用过Write
	Started bool
	// 回复状态
	status      int
	Gzip        *gzip.Writer
	IsOpenGzip  bool
	NeedGzipLen int
	isGzip      bool
	MsgData     map[string]interface{}
	//Cookies []*http.Cookie

	sessionFunc      ISession
	session          map[string]interface{}
	isGetSession     bool
	sessionName      string
	SessionAliveTime time.Duration
}

// Session接口
type ISession interface {
	getSession(KeyValue string) (map[string]interface{}, error)
	getSessionKeyValue() (string, error)
	setSession(SessionName string, m map[string]interface{}, duration time.Duration) error
	removeSession(KeyValue string)
}

func (r *Response) RemoveSession() {
	r.session = nil
	r.isGetSession = true
	if r.sessionName != "" {
		r.sessionFunc.removeSession(r.sessionName)
	}
}

func (r *Response) GetSession() (map[string]interface{}, error) {
	if r.isGetSession {
		return r.session, nil
	} else {
		var err error
		if r.sessionName == "" {
			return r.session, err
		} else {
			if r.session, err = r.sessionFunc.getSession(r.sessionName); err == nil {
				r.isGetSession = true
			}
		}
		return r.session, err
	}
}

func (r *Response) SetSession(key string, val interface{}) error {
	if r.isGetSession {
		if r.session == nil {
			r.session = map[string]interface{}{}
		}
		r.session[key] = val
		if r.sessionName == "" {
			if key, err := r.sessionFunc.getSessionKeyValue(); err == nil {
				return r.sessionFunc.setSession(key, r.session, r.SessionAliveTime)
			} else {
				return err
			}
		} else {
			return r.sessionFunc.setSession(key, r.session, r.SessionAliveTime)
		}
	} else {
		var err error
		if r.session, err = r.GetSession(); err != nil {
			return err
		}
		if r.session == nil {
			r.session = map[string]interface{}{}
		}
		r.session[key] = val
		if r.sessionName == "" {
			if key, err := r.sessionFunc.getSessionKeyValue(); err == nil {
				return r.sessionFunc.setSession(key, r.session, r.SessionAliveTime)
			} else {
				return err
			}
		} else {
			return r.sessionFunc.setSession(key, r.session, r.SessionAliveTime)
		}
	}
}

// 设置Cookie
func (r *Response) SetCookie(c *http.Cookie) {
	http.SetCookie(r, c)
}

// 通过cookie名字移除Cookie
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

// 通过cookie移除Cookie
func (r *Response) RemoveCookie(ck *http.Cookie) {
	ck.Expires = time.Now()
	http.SetCookie(r, ck)
}

type buf struct {
	r *Response
}

func (buf *buf) Write(p []byte) (n int, err error) {
	return buf.r.ResponseWriter.Write(p)
}

// 写入前端的数据
func (r *Response) Write(b []byte) (int, error) {
	defer func() {
		r.Started = true
	}()
	if r.isGzip || (r.IsOpenGzip && r.NeedGzipLen < len(b) && !r.Started) {
		if !r.Started {
			r.isGzip = true
			r.ResponseWriter.Header().Set("Content-Encoding", "gzip")
			if r.ResponseWriter.Header().Get("Content-Type") == "" {
				r.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(b))
			}
			r.ResponseWriter.WriteHeader(r.status)
		}
		if r.Gzip == nil {
			buf := buf{r: r}
			r.Gzip = gzip.NewWriter(&buf)
		}
		return r.Gzip.Write(b)
	}
	if !r.Started {
		r.ResponseWriter.WriteHeader(r.status)
	}
	return r.ResponseWriter.Write(b)
}

// 写入前端的json数据
func (r *Response) WriteJson(i interface{}) error {
	if b, err := json.Marshal(i); err == nil {
		contentType := http.DetectContentType(b)
		f := regexp.MustCompile(`(;[\ ]?charset=.*)`)
		t := f.FindAllStringSubmatch(contentType, 1)
		contentType = "application/json"
		if len(t) > 0 && len(t[0]) > 0 {
			contentType = contentType + t[0][0]
		}
		r.ResponseWriter.Header().Set("Content-Type", contentType)
		_, err = r.Write(b)
		return err
	} else {
		return err
	}
}

func (r *Response) Redirect(redirectPath string) {
	http.Redirect(r, r.req, redirectPath, 302)
}
