package main

import (
	middwares "goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/config"
	c "goblog/pkg/config"
	"net/http"
)

func init() {
	config.Initialize()
}

func main() {
    bootstrap.SetUpDB()
    router := bootstrap.SetupRoute()

    http.ListenAndServe(":" + c.GetString("app.port"), middwares.RemoveTrailingSlash(router))
}