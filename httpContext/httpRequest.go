package httpContext

import (
	"net/http"
)

type Request struct {
	*http.Request
}
