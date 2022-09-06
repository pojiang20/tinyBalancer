package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func maxAllowedMiddleware(n uint) mux.MiddlewareFunc {
	//定义通道容量
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			//消耗1
			acquire()
			//释放1
			defer release()
			next.ServeHTTP(w, request)
		})
	}
}
