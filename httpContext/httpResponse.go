package httpContext

import (
	"compress/gzip"
	"log"
	"net/http"
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
	ck.Expires = time.Now()
	http.SetCookie(r, ck)
}

type buf struct {
	r *Response
}

func (buf *buf) Write(p []byte) (n int, err error) {
	return buf.r.ResponseWriter.Write(p)
}

// 添加回复数据
func (r *Response) Write(b []byte) (int, error) {
	defer func() {
		r.Started = true
	}()
	log.Println(len(b))
	log.Println(r.NeedGzipLen)
	if r.isGzip || (r.IsOpenGzip && r.NeedGzipLen < len(b)) {
		if !r.Started {
			r.isGzip = true
			r.ResponseWriter.Header().Set("Content-Encoding", "gzip")
			r.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(b))
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
func (r *Response) Redirect(redirectPath string) {
	r.status = 302
	http.Redirect(r.ResponseWriter, r.req, redirectPath, 302)
}
