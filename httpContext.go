package freeFishGo

import "net/http"

type HttpContext struct {
	Response Response
	Request  *Request
}

func (h *HttpContext) SetContext(rw http.ResponseWriter, r *http.Request) {
	h.Response = Response{ResponseWriter: rw, req: r, Started: false}
	h.Response.WriteHeader(200)
	h.Request = &Request{Request: r}
}
