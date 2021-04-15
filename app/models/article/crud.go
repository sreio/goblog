package article

import (
	"goblog/pkg/model"
	"goblog/pkg/types"
)




func Get(idStr string) (Article, error){
	var article Article
	id := types.StringToInt(idStr)
	err := model.DB.First(&article, id).Error
	if err != nil {
		return article, err
	}
	return article, nil
}