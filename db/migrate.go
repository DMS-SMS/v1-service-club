package db

import (
	"club/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) (err error) {
	migrator := db.Migrator()

	//_ = migrator.DropTable(&model.ClubMember{})
	//_ = migrator.DropTable(&model.ClubInform{})
	//_ = migrator.DropTable(&model.RecruitMember{})
	//_ = migrator.DropTable(&model.ClubRecruitment{})
	//_ = migrator.DropTable(&model.Club{})

	if !migrator.HasTable(&model.Club{}) {
		if err = migrator.CreateTable(&model.Club{}); err != nil { return }
	}
	if !migrator.HasTable(&model.ClubInform{}) {
		if err = migrator.CreateTable(&model.ClubInform{}); err != nil { return }
	}
	if !migrator.HasTable(&model.ClubMember{}) {
		if err = migrator.CreateTable(&model.ClubMember{}); err != nil { return }
	}
	if !migrator.HasTable(&model.ClubRecruitment{}) {
		if err = migrator.CreateTable(&model.ClubRecruitment{}); err != nil { return }
	}
	if !migrator.HasTable(&model.RecruitMember{}) {
		if err = migrator.CreateTable(&model.RecruitMember{}); err != nil { return }
	}

	return db.AutoMigrate(&model.Club{}, &model.ClubInform{}, &model.ClubMember{}, &model.ClubRecruitment{}, &model.RecruitMember{})
}
