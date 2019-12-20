package httpContext

import (
	"net/http"
	"time"
)

type Request struct {
	*http.Request
	sessionFunc  ISession
	session      map[string]interface{}
	isGetSession bool
	sessionName  string
}

// Session接口
type ISession interface {
	getSession(KeyValue string) (map[string]interface{}, error)
	getSessionKeyValue() (string, error)
	setSession(KeyValue string, m map[string]interface{}, duration time.Duration) error
	removeSession(KeyValue string)
}

func (r *Request) RemoveSession() {
	r.sessionFunc.removeSession(r.sessionName)
}
