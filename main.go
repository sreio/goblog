package main

import (
    "fmt"
    "net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type","text/html; charset=utf-8")
    switch r.URL.Path {
        case "/":
            fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
        case "/home":
            fmt.Fprint(w, "<h1>Hello, 这里是 home</h1>")
        case "/me":
            fmt.Fprint(w, "<h1>Hello, I‘m Sreio。</h1>")
        default:
            fmt.Fprint(w, "默认跳转路由")
    }
}

func main() {
    http.HandleFunc("/", handlerFunc)
    http.ListenAndServe(":3000", nil)
}