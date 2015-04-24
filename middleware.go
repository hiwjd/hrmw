package hrmw

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MiddlewareHandle func(http.ResponseWriter, *http.Request, httprouter.Params, *Middleware)
type Middleware struct {
	cur      int
	handlers []MiddlewareHandle
}

func (m *Middleware) Next(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if m.cur < len(m.handlers) {
		h := m.handlers[m.cur]
		m.cur++
		h(w, r, ps, m)
	}
}

func (m *Middleware) Last(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	len := len(m.handlers)
	h := m.handlers[len-1]
	m.cur = len
	h(w, r, ps, m)
}

func Use(handlers ...MiddlewareHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		m := &Middleware{0, handlers}
		m.Next(w, r, ps)
	}
}
