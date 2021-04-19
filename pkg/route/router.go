package route

import (
	"goblog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)
var route  *mux.Router

// SetRoute 设置路由实例，以供 Name2URL 等函数使用
func SetRoute(r *mux.Router) {
    route = r
}

// RouteName2URL 通过路由名称来获取 URL
func RouteName2URL(routerName string, pars ...string) string {

    url, err := route.Get(routerName).URL(pars...)
    if err != nil {
        logger.LogError(err)
        return ""
    }
    return url.String()
}

// 获取请求参数
func GetRouterParam(parameterName string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[parameterName]
}