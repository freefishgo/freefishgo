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
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"regexp"
	"time"
)

type IResponse interface {
	// 升级为WebSocket服务 upgrades为空时采用默认的参数 为多个时只采用第一个作为WebSocket参数
	WebSocket(upgrades ...*websocket.Upgrader) (conn *websocket.Conn, err error)
	Hijack() (net.Conn, *bufio.ReadWriter, error)
	setISession(i ISession)
	getSessionKeyValue() (string, error)
	RemoveSession()
	GetSession(key string) interface{}
	getSession() error
	UpdateSession() error
	SetSession(key string, val interface{})
	// 设置Cookie
	SetCookie(c *http.Cookie)
	// 设置Cookie
	SetCookieUseKeyValue(key string, val string)
	// 通过cookie名字移除Cookie
	RemoveCookieByName(name string)
	WriteHeader(statusCode int)
	ReadStatusCode() int
	// 通过cookie移除Cookie
	RemoveCookie(ck *http.Cookie)
	// 写入前端的数据
	Write(b []byte) (int, error)
	// 获取写入前端的缓存
	GetWriteCache() []byte
	// 清除写入前端的缓存
	ClearWriteCache()
	// 写入前端的json数据
	WriteJson(i interface{}) error
	Redirect(redirectPath string)
	getYourself() *Response
	GetStarted() bool
	GetIsWriteInCache() bool
	SetIsWriteInCache(bool)
	Header() http.Header
}

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
	MsgData     map[interface{}]interface{}

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

// 升级为WebSocket服务 upgrades为空时采用默认的参数 为多个时只采用第一个作为WebSocket参数
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

func (r *Response) getSessionKeyValue() (string, error) {
	r.isUpdateSessionKey = true
	return r.sessionFunc.GetSessionKeyValue()
}

func (r *Response) RemoveSession() {
	r.session = nil
	r.isGetSession = false
	r.sessionIsUpdate = false
	r.sessionFunc.RemoveBySessionID(r.SessionId)
	r.SessionId = ""
}

func (r *Response) GetSession(key string) interface{} {
	if r.SessionId == "" {
		return nil
	}
	if !r.isGetSession {
		var err error
		if err = r.getSession(); err == nil {
			r.isGetSession = true
			if r.session == nil {
				return nil
			}
		}
	}
	v, _ := r.session[key]
	return v
}

func (r *Response) getSession() error {
	if r.SessionId == "" {
		return errors.New("没有设置Session")
	}
	var err error
	if !r.isGetSession {
		if r.session, err = r.sessionFunc.GetSession(r.SessionId); err == nil {
			if r.session == nil {
				r.SessionId = ""
			}
			r.isGetSession = true
		}
	}
	return err
}

func (r *Response) UpdateSession() error {
	if r.SessionId == "" {
		return nil
	}
	if r.sessionIsUpdate {
		return r.sessionFunc.SetSession(r.SessionId, r.session)
	}
	return nil
}

func (r *Response) SetSession(key string, val interface{}) {
	r.sessionIsUpdate = true
	r.getSession()
	if r.SessionId == "" {
		r.isUpdateSessionKey = true
		if SessionName, err := r.sessionFunc.GetSessionKeyValue(); err == nil {
			r.SessionId = SessionName
		}
	}
	if r.session == nil {
		r.session = map[interface{}]interface{}{}
	}
	r.session[key] = val
}

// 设置Cookie
func (r *Response) SetCookie(c *http.Cookie) {
	http.SetCookie(r, c)
}

// 设置Cookie
func (r *Response) SetCookieUseKeyValue(key string, val string) {
	http.SetCookie(r.ResponseWriter, &http.Cookie{Name: key, Value: val})
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

// 写入前端的数据
func (r *Response) Write(b []byte) (int, error) {
	if r.isWriteInCache && len(r.writeCache) < r.maxResponseCacheLen {
		r.writeCache = append(r.writeCache, b...)
		return len(b), nil
	}
	defer func() {
		r.Started = true
	}()
	if r.isGzip || (r.IsOpenGzip && r.NeedGzipLen < len(b)+len(r.writeCache) && !r.Started) {
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
			r.ResponseWriter.Header().Set("Content-Encoding", "gzip")
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

// 获取写入前端的缓存
func (r *Response) GetWriteCache() []byte {
	return r.writeCache
}

// 清除写入前端的缓存
func (r *Response) ClearWriteCache() {
	r.writeCache = nil
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
