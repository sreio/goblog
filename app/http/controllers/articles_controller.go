package controllers

import (
	"database/sql"
	"fmt"
	"goblog/app/models/article"
	"goblog/app/requests"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"

	"gorm.io/gorm"
)

type ArticlesController struct{}

//详情
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouterParam("id", r)
	article, err := article.Get(id)

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
        view.Render(w, view.D{"Article":article}, "articles.show", "articles._form_field")
	}
}

//list列表
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	//获取结果集
	articles, err := article.GetAll()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器错误")
	} else {
		view.Render(w, view.D{"Articles":articles}, "articles.index")
	}
}

//编辑页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouterParam("id", r)
	article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		// 4. 读取成功，显示编辑文章表单
        view.Render(w, view.D{
            "Article": article,
            "Errors":  nil,
        }, "articles.edit", "articles._form_field")
	}
}

//修改
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouterParam("id", r)
	_article, err := article.Get(id)

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
        _article = article.Article{
            Title: r.PostFormValue("title"),
            Body: r.PostFormValue("body"),
        }
		errors := requests.ValidateArticleForm(_article)
		if len(errors) == 0 {

			rowsAffected, err := _article.Update()
			
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
            	fmt.Fprint(w, "500 服务器内部错误")
				return
			}

			if rowsAffected > 0 {
				showUrl := route.RouteName2URL("articles.show", "id", id)
				http.Redirect(w, r, showUrl, http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改！")
			}
		} else {
            // 4.3 表单验证不通过，显示理由
            view.Render(w, view.D{
                "Article": _article,
                "Errors":  errors,
            }, "articles.edit", "articles._form_field")
		}
		
	}
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
    view.Render(w, view.D{}, "articles.edit", "articles._form_field")
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

    _article := article.Article{
        Title: r.PostFormValue("title"),
        Body: r.PostFormValue("body"),
    }

    errors := requests.ValidateArticleForm(_article)

    // 检查是否有错误
    if len(errors) == 0 {
        _article.Create()
        if _article.ID > 0 {
            fmt.Fprint(w, "插入成功，ID 为"+ _article.GetStringID())
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "创建文章失败，请联系管理员")
        }
    } else {
        view.Render(w, view.D{
            "Article": _article,
            "Errors": errors,
        }, "articles.edit", "articles._form_field")
    }
}

// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {

    // 1. 获取 URL 参数
    id := route.GetRouterParam("id", r)

    // 2. 读取对应的文章数据
    _article, err := article.Get(id)

    // 3. 如果出现错误
    if err != nil {
        if err == gorm.ErrRecordNotFound {
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
        rowsAffected, err := _article.Delete()

        // 4.1 发生错误
        if err != nil {
            // 应该是 SQL 报错了
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        } else {
            // 4.2 未发生错误
            if rowsAffected > 0 {
                // 重定向到文章列表页
                indexURL := route.RouteName2URL("articles.index")
                http.Redirect(w, r, indexURL, http.StatusFound)
            } else {
                // Edge case
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprint(w, "404 文章未找到")
            }
        }
    }
}