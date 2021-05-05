package bootstrap

import (
	"goblog/pkg/route"
	"goblog/pkg/routes"
	"github.com/gorilla/mux"
)

// SetupRoute 路由初始化
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
    route.SetRoute(router)
	routes.RegisterWebRoutes(router)
	return router
}
