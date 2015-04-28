package hrmw

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MiddlewareHandle func(http.ResponseWriter, *http.Request, httprouter.Params, *Middleware)
type Middleware struct {
	cur      int
	handlers []MiddlewareHandle
	data     map[string]interface{}
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

func (m *Middleware) Set(key string, v interface{}) {
	m.data[key] = v
}

func (m *Middleware) Get(key string) interface{} {
	v := m.data[key]
	delete(m.data, key)
	return v
}

func Use(handlers ...MiddlewareHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		m := &Middleware{}
		m.cur = 0
		m.handlers = handlers
		m.data = make(map[string]interface{})
		m.Next(w, r, ps)
	}
}

type Pattern struct {
	first MiddlewareHandle
	last  MiddlewareHandle
}

func NewPattern() *Pattern {
	return &Pattern{nil, nil}
}

func (p *Pattern) First(handle MiddlewareHandle) *Pattern {
	p.first = handle
	return p
}

func (p *Pattern) Last(handle MiddlewareHandle) *Pattern {
	p.last = handle
	return p
}

func (p *Pattern) Use(handlers ...MiddlewareHandle) httprouter.Handle {
	hs := make([]MiddlewareHandle, 0)
	if p.first != nil {
		hs = append(hs, p.first)
	}
	hs = append(hs, handlers...)
	if p.last != nil {
		hs = append(hs, p.last)
	}

	return Use(hs...)
}
