package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var db *sql.DB
func initDB() {
    var err error
    config := mysql.Config{
        User : "root",
        Passwd: "",
        Addr: "127.0.0.1:3306",
        Net: "tcp",
        DBName: "goblog",
        AllowNativePasswords: true,
    }
    //连接数据库
    db, err = sql.Open("mysql", config.FormatDSN())
    checkError(err)

    db.SetMaxOpenConns(25) //最大连接数
    db.SetMaxIdleConns(25) //最大空闲数
    db.SetConnMaxLifetime(5 * time.Minute)//每个连接的过期时间

    err = db.Ping()
    checkError(err)
}

func checkError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func init() {
    initDB()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
        "<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    fmt.Fprint(w, "文章 ID："+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "访问文章列表")
}

type ArticlesFormData struct {
    Title, Body string
    URL *url.URL
    Errors map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
    title := r.PostFormValue("title")
    body := r.FormValue("body")

    errors := make(map[string]string)

    if title == "" {
        errors["title"] = "标题不能为空"
    } else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 20{
        errors["title"] = "标题长度需介于 3-20"
    }

    if body == "" {
        errors["body"] = "内容不能为空"
    } else if utf8.RuneCountInString(body) < 10{
        errors["body"] = "内容长度需大于或等于 10 个字节"
    }

    if len(errors) == 0 {
        fmt.Fprint(w, "验证通过!<br>")
        fmt.Fprintf(w, "title 的值为: %v <br>", title)
        fmt.Fprintf(w, "title 的长度为: %v <br>", utf8.RuneCountInString(title))
        fmt.Fprintf(w, "body 的值为: %v <br>", body)
        fmt.Fprintf(w, "body 的长度为: %v <br>", utf8.RuneCountInString(body))
    } else {
        storeURL, _ := router.Get("articles.store").URL()
        data := ArticlesFormData{
            Title:  title,
            Body:   body,
            URL:    storeURL,
            Errors: errors,
        }
        // tmpl, err := template.New("create-form").Parse(html)
        tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
        if err != nil {
            panic(err)
        }

        tmpl.Execute(w, data)
    }

}

func articlesAddHandler(w http.ResponseWriter, r *http.Request) {
    storeURL, _ := router.Get("articles.store").URL()
    tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
    if err != nil {
        panic(err)
    }
    data := ArticlesFormData{
        Title: "",
        Body: "",
        URL: storeURL,
        Errors: nil,
    }
    tmpl.Execute(w, data)
   
}


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
    router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
    router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

    router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
    router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
    router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
    router.HandleFunc("/articles/add", articlesAddHandler).Methods("GET").Name("articles.add")

    // 自定义 404 页面
    router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

    router.Use(foreaHtmlMiddlewaer)

    // 通过命名路由获取 URL 示例
    // homeURL, _ := router.Get("home").URL()
    // fmt.Println("homeURL: ", homeURL)
    // articleURL, _ := router.Get("articles.show").URL("id", "23")
    // fmt.Println("articleURL: ", articleURL)

    http.ListenAndServe(":3000", removeTrailingSlash(router))
}