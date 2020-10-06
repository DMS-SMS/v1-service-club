package test

import (
	"club/tool/random"
	"gorm.io/gorm"
	"time"
)

func createGormModelOnCurrentTime() gorm.Model {
	currentTime := time.Now()
	return gorm.Model{
		ID:        uint(random.Int64WithLength(3)),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		DeletedAt: nil,
	}
}