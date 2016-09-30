package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
	"net/http"
)

type Router struct {
	rules     map[string]Handler
	dbSession *r.Session
}
type Handler func(*Client, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewRouter(dbSession *r.Session) *Router {
	return &Router{
		rules:     make(map[string]Handler),
		dbSession: dbSession,
	}
}

func (r *Router) Handle(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	client := NewClient(socket, e.FindHandler, e.dbSession)
	defer client.Close()
	go client.Write()
	client.Read()
}
