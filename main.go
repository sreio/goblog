package main

import (
	"goblog/bootstrap"
	"goblog/pkg/database"
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
    database.Initialize()
    bootstrap.SetUpDB()
    router = bootstrap.SetupRoute()
    router.Use(foreaHtmlMiddlewaer)

    // 通过命名路由获取 URL 示例
    // homeURL, _ := router.Get("home").URL()
    // fmt.Println("homeURL: ", homeURL)
    // articleURL, _ := router.Get("articles.show").URL("id", "23")
    // fmt.Println("articleURL: ", articleURL)

    http.ListenAndServe(":3000", removeTrailingSlash(router))
}