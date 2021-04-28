package models

import (
	"goblog/pkg/types"
	"time"
)

type BaseModel struct{
	ID uint64

	CreatedAt time.Time
    UpdatedAt time.Time
}

func (a BaseModel) GetStringID() string {
	return types.Uint64ToString(a.ID)
}