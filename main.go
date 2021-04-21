package main

import (
	middwares "goblog/app/http/middlewares"
	"goblog/bootstrap"
	"net/http"
)

func main() {
    bootstrap.SetUpDB()
    router := bootstrap.SetupRoute()

    http.ListenAndServe(":3000", middwares.RemoveTrailingSlash(router))
}