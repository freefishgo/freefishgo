package httpContext

import "net/http"

type Response struct {
	http.ResponseWriter
	Cookies *http.Cookie
}

func (r *Response) SetCookie(c *http.Cookie) {

}
