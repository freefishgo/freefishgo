package httpContext

import "net/http"

type Request struct {
	*http.Request
}

func (r *Request) get() {
	r.Cookies()
}
