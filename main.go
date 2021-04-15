package main

import (
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"database/sql"

	"github.com/gorilla/mux"
)

var router *mux.Router
var db *sql.DB


type Article struct{
    Title, Body string
    ID int64
}

type ArticlesFormData struct {
    Title, Body string
    URL *url.URL
    Errors map[string]string
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprint(w, "访问文章列表")
    rows, err := db.Query("SELECT * from articles")
    logger.LogError(err)
    defer rows.Close()

    var articles []Article

    // 遍历数据
    for rows.Next() {
        var article Article
        // 扫描每一行的结果并赋值到一个 article 对象中
        err := rows.Scan(&article.ID, &article.Title, &article.Body)
        logger.LogError(err)
        articles = append(articles, article)
    }
    // 检查遍历的时候是否发生错误
    err = rows.Err()
    logger.LogError(err)

    tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
    logger.LogError(err)
    tmpl.Execute(w, articles)
}


func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
    title := r.PostFormValue("title")
    body := r.FormValue("body")
    errors := validateArticleFormData(title, body)

    if len(errors) == 0 {
        lastInsterID, err := saveArticleToDB(title, body)
       if lastInsterID > 0 {
            fmt.Fprint(w, "新增成功，ID：" + strconv.FormatInt(lastInsterID, 10))
       } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器错误")
       }
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

func saveArticleToDB(title string, body string) (int64, error) {

    // 变量初始化
    var (
        id   int64
        err  error
        rs   sql.Result
        stmt *sql.Stmt
    )

    // 1. 获取一个 prepare 声明语句
    stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
    // 例行的错误检测
    if err != nil {
        return 0, err
    }

    // 2. 在此函数运行结束后关闭此语句，防止占用 SQL 连接
    defer stmt.Close()

    // 3. 执行请求，传参进入绑定的内容
    rs, err = stmt.Exec(title, body)
    if err != nil {
        return 0, err
    }

    // 4. 插入成功的话，会返回自增 ID
    if id, err = rs.LastInsertId(); id > 0 {
        return id, nil
    }

    return 0, err
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

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
    id := getRouterParam("id", r)
    article, err := getArtilceByID(id)

    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器错误")
        }
    } else {
        updateUrl, _ := router.Get("articles.update").URL("id", id)
        data := ArticlesFormData{
            Title: article.Title,
            Body: article.Body,
            URL: updateUrl,
            Errors: nil,
        }
        tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
        logger.LogError(err)
        tmpl.Execute(w, data)
    }
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
    // 1. 获取 URL 参数
    id := getRouterParam("id", r)

    // 2. 读取对应的文章数据
    _, err := getArtilceByID(id)

    // 3. 如果出现错误
    if err != nil {
        if err == sql.ErrNoRows {
            // 3.1 数据未找到
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            // 3.2 数据库错误
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        // 4. 未出现错误

        // 4.1 表单验证
        title := r.PostFormValue("title")
        body := r.PostFormValue("body")
        errors := validateArticleFormData(title, body)

        if len(errors) == 0 {

            // 4.2 表单验证通过，更新数据
            query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
            rs, err := db.Exec(query, title, body, id)

            if err != nil {
                logger.LogError(err)
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "500 服务器内部错误")
            }

            // √ 更新成功，跳转到文章详情页
            if n, _ := rs.RowsAffected(); n > 0 {
                showURL, _ := router.Get("articles.show").URL("id", id)
                http.Redirect(w, r, showURL.String(), http.StatusFound)
            } else {
                fmt.Fprint(w, "您没有做任何更改！")
            }
        } else {

            // 4.3 表单验证不通过，显示理由

            updateURL, _ := router.Get("articles.update").URL("id", id)
            data := ArticlesFormData{
                Title:  title,
                Body:   body,
                URL:    updateURL,
                Errors: errors,
            }
            tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
            logger.LogError(err)

            tmpl.Execute(w, data)
        }
    }
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
    // 1. 获取 URL 参数
    id := getRouterParam("id", r)

    // 2. 读取对应的文章数据
    article, err := getArtilceByID(id)

    // 3. 如果出现错误
    if err != nil {
        if err == sql.ErrNoRows {
            // 3.1 数据未找到
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            // 3.2 数据库错误
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        // 4. 未出现错误，执行删除操作
        rowsAffected, err := article.Delete()

        // 4.1 发生错误
        if err != nil {
            // 应该是 SQL 报错了
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        } else {
            // 4.2 未发生错误
            if rowsAffected > 0 {
                // 重定向到文章列表页
                indexURL, _ := router.Get("articles.index").URL()
                http.Redirect(w, r, indexURL.String(), http.StatusFound)
            } else {
                // Edge case
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprint(w, "404 文章未找到")
            }
        }
    }
}

func (t Article) Delete() (rowsAffected int64, err error) {
    rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(t.ID, 10))

    if err != nil {
        return 0, err
    }

    // √ 删除成功，跳转到文章详情页
    if n, _ := rs.RowsAffected(); n > 0 {
        return n, nil
    }

    return 0, nil
}


func getArtilceByID(id string) (Article, error){
    article := Article{}
    query := "SELECT * FROM articles WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
    return article, err
}

func validateArticleFormData(title string, body string) map[string]string {
    errors := make(map[string]string)
    
    // 验证标题
    if title == "" {
        errors["title"] = "标题不能为空"
    } else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
        errors["title"] = "标题长度需介于 3-40"
    }

    // 验证内容
    if body == "" {
        errors["body"] = "内容不能为空"
    } else if utf8.RuneCountInString(body) < 10 {
        errors["body"] = "内容长度需大于或等于 10 个字节"
    }

    return errors
}

func (t Article) Link() string {
    showUrl, err := router.Get("articles.show").URL("id", strconv.FormatInt(t.ID, 10))
    if err != nil {
        logger.LogError(err)
        return ""
    }
    return showUrl.String()
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

// 获取请求参数
func getRouterParam(parameterName string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[parameterName]
}

func main() {
    database.Initialize()
    db = database.DB
    bootstrap.SetUpDB()
    router = bootstrap.SetupRoute()

    

    router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
    router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
    router.HandleFunc("/articles/add", articlesAddHandler).Methods("GET").Name("articles.add")
    router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
    router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
    router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

    router.Use(foreaHtmlMiddlewaer)

    // 通过命名路由获取 URL 示例
    // homeURL, _ := router.Get("home").URL()
    // fmt.Println("homeURL: ", homeURL)
    // articleURL, _ := router.Get("articles.show").URL("id", "23")
    // fmt.Println("articleURL: ", articleURL)

    http.ListenAndServe(":3000", removeTrailingSlash(router))
}