package main

import (
	middwares "goblog/app/http/middlewares"
	"goblog/bootstrap"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var router *mux.Router

func foreaHtmlMiddlewaer (next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        next.ServeHTTP(w, r)
    })
}

func removeTrailingSlash (next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
        if r.URL.Path != "/" {
            r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    bootstrap.SetUpDB()
    router = bootstrap.SetupRoute()

    http.ListenAndServe(":3000", middwares.RemoveTrailingSlash(router))
}