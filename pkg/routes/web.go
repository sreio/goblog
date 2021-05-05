package routes

import (
	"goblog/app/http/controllers"
	middwares "goblog/app/http/middlewares"
	// middwares "goblog/app/http/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(r *mux.Router) {
	
	//静态页面
	pc := new(controllers.PackageController)
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// 中间件：强制内容类型为 HTML
	// r.Use(middwares.ForceHTML)
	//文章模块
	ac := new(controllers.ArticlesController)
	r.HandleFunc("/articles/{id:[0-9]+}", middwares.Auth(ac.Show)).Methods("GET").Name("articles.show")
	r.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", middwares.Auth(ac.Edit)).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", middwares.Auth(ac.Update)).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/create", middwares.Auth(ac.Create)).Methods("GET").Name("articles.create")
    r.HandleFunc("/articles", middwares.Auth(ac.Store)).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", middwares.Auth(ac.Delete)).Methods("POST").Name("articles.delete")

	// 用户相关
	auc := new(controllers.AuthController)
	r.HandleFunc("/auth/register", middwares.Guest(auc.Register)).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", middwares.Guest(auc.DoRegister)).Methods("POST").Name("auth.doregister")
	r.HandleFunc("/auth/login", middwares.Guest(auc.Login)).Methods("GET").Name("auth.login")
    r.HandleFunc("/auth/dologin", middwares.Guest(auc.DoLogin)).Methods("POST").Name("auth.dologin")
    r.HandleFunc("/auth/logout", middwares.Auth(auc.Logout)).Methods("POST").Name("auth.logout")

	// 用户认证
    uc := new(controllers.UserController)
    r.HandleFunc("/users/{id:[0-9]+}", uc.Show).Methods("GET").Name("users.show")
	
	// --- 全局中间件 ---
    // 开始会话
    r.Use(middwares.StartSession)
}