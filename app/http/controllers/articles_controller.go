package controllers

import (
	"database/sql"
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
)

type ArticlesController struct{}

func (*ArticlesController) Show(w http.ResponseWriter,r *http.Request) {
	id := route.GetRouterParam("id", r)
    article, err := getArtilceByID(id)

    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            // 3.2 数据库错误
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        // tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
        tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
            "RouteName2URL": route.RouteName2URL,
            "Int64ToString": types.Int64ToString,
        }).ParseFiles("resources/views/articles/show.gohtml")
        logger.LogError(err)
        tmpl.Execute(w, article)
    }
}