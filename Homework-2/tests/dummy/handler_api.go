package dummy

import "net/http"

type HandlerApi struct {
}

func NewHandlerApi() *HandlerApi {
	return &HandlerApi{}
}

func (*HandlerApi) Create(w http.ResponseWriter, r *http.Request) {}

func (*HandlerApi) Read(w http.ResponseWriter, r *http.Request) {}

func (*HandlerApi) Update(w http.ResponseWriter, r *http.Request) {}

func (*HandlerApi) Delete(w http.ResponseWriter, r *http.Request) {}
