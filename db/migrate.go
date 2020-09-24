package db

import (
	"club/model"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.LogMode(false)

	db.DropTableIfExists(&model.Club{})
	db.DropTableIfExists(&model.ClubInform{})
	db.DropTableIfExists(&model.ClubMember{})
	db.DropTableIfExists(&model.ClubRecruitment{})
	db.DropTableIfExists(&model.RecruitMember{})

	if !db.HasTable(&model.Club{}) {
		db.CreateTable(&model.Club{})
	}
	if !db.HasTable(&model.ClubInform{}) {
		db.CreateTable(&model.ClubInform{})
	}
	if !db.HasTable(&model.ClubMember{}) {
		db.CreateTable(&model.ClubMember{})
	}
	if !db.HasTable(&model.ClubRecruitment{}) {
		db.CreateTable(&model.ClubRecruitment{})
	}
	if !db.HasTable(&model.RecruitMember{}) {
		db.CreateTable(&model.RecruitMember{})
	}

	db.AutoMigrate(&model.Club{}, &model.ClubInform{}, &model.ClubMember{}, &model.ClubRecruitment{}, &model.RecruitMember{})
	db.Model(&model.ClubInform{}).AddForeignKey("club_uuid", "clubs(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&model.ClubMember{}).AddForeignKey("club_uuid", "clubs(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&model.ClubRecruitment{}).AddForeignKey("club_uuid", "clubs(uuid)", "RESTRICT", "RESTRICT")
	db.Model(&model.RecruitMember{}).AddForeignKey("recruitment_uuid", "club_recruitments(uuid)", "RESTRICT", "RESTRICT")
}
