package route

import (
	"goblog/pkg/routes"
	"net/http"

	"github.com/gorilla/mux"
)



var Router  *mux.Router

// Initialize 初始化路由
func Initialize() {
    Router = mux.NewRouter()
	routes.RegisterWebRoutes(Router)
}

// RouteName2URL 通过路由名称来获取 URL
func RouteName2URL(routerName string, pars ...string) string {
    url, err := Router.Get(routerName).URL(pars...)
    if err != nil {
        // checkError(err)
        return ""
    }
    return url.String()
}

// 获取请求参数
func GetRouterParam(parameterName string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[parameterName]
}