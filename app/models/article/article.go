package article

import (
	"goblog/app/models"
	"goblog/app/models/user"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/route"
)


type Article struct {
	models.BaseModel

	Title string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body  string `gorm:"type:longtext;not null;" valid:"body"`
	
	UserID uint64 `gorm:"not null;index"`
    User   user.User
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

// CreatedAtDate 创建日期
func (a Article) CreatedAtDate() string {
    return a.CreatedAt.Format("2006-01-02")
}