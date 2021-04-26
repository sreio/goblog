package article

import (
	"goblog/app/models"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/route"
)


type Article struct {
	models.BaseModel

	Title string
	Body  string
}

// Link 方法用来生成文章链接
func (a Article) Link() string {
    return route.RouteName2URL("articles.show", "id", a.GetStringID())
}

func (article *Article) Update() (rowsAffected int64, err error){
	result := model.DB.Save(&article)

	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// Delete 删除文章
func (article *Article) Delete() (rowsAffected int64, err error) {
    result := model.DB.Delete(&article)
    if err = result.Error; err != nil {
        logger.LogError(err)
        return 0, err
    }

    return result.RowsAffected, nil
}