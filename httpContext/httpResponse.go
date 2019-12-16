package httpContext

import (
	"net/http"
	"time"
)

type Response struct {
	http.ResponseWriter
	req *http.Request
	// 是否调用过Write
	Started bool
	// 回复状态
	status int
	//写到前端的数据
	writeData []byte
	//Cookies []*http.Cookie
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
	http.SetCookie(r, ck)
}

// 添加回复数据
func (r *Response) Write(b []byte) (int, error) {
	r.Started = true
	r.writeData = append(r.writeData, b...)
	return len(b), nil
}

func (r *Response) GetWaitWriteData() []byte {
	return r.writeData
}
func (r *Response) ClearWaitWriteData() {
	r.writeData = nil
}
