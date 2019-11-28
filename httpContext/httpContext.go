package httpContext

import "net/http"

type HttpContext struct {
	Response http.ResponseWriter
	Request  *http.Request
}
