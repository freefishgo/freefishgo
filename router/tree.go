package router

import "net/http"

type tree struct {
	Tree   *tree
	Path   string
	Handle http.Handler
}
