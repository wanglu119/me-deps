package mux

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type Router = mux.Router

func NewRouter() *Router {
	return mux.NewRouter()
}

func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// addr: ip:port
func Serve(router *Router, addr string) {
	var listener net.Listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Listening on: ", addr)
	err = http.Serve(listener, router)
	if err != nil {
		panic(err)
	}
}
