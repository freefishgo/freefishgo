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

import "net/http"

type HttpContext struct {
	Response IResponse
	Request  *Request
}

func (h *HttpContext) setContext(rw http.ResponseWriter, r *http.Request) {
	h.Response = &Response{ResponseWriter: rw, req: r, Started: false}
	h.Response.WriteHeader(200)
	h.Request = &Request{Request: r}
}
