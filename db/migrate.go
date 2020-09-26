package db

import (
	"club/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	migrator := db.Migrator()

	_ = migrator.DropTable(&model.ClubMember{})
	_ = migrator.DropTable(&model.ClubInform{})
	_ = migrator.DropTable(&model.RecruitMember{})
	_ = migrator.DropTable(&model.ClubRecruitment{})
	_ = migrator.DropTable(&model.Club{})

	return db.AutoMigrate(&model.Club{}, &model.ClubInform{}, &model.ClubMember{}, &model.ClubRecruitment{}, &model.RecruitMember{})
}
