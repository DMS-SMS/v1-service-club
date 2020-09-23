package model

import (
	"github.com/jinzhu/gorm"
)

type Club struct {
	gorm.Model
	UUID       uuid       `gorm:"PRIMARY_KEY;Type:char(17);UNIQUE;INDEX" validate:"uuid=club,len=17"`
	LeaderUUID leaderUUID `gorm:"Type:char(20);NOT NULL;UNIQUE;" validate:"uuid=student,len=20"`
}

type ClubInform struct {
	gorm.Model
	ClubUUID     clubUUID     `gorm:"Type:char(17);NOT NULL;UNIQUE" validate:"uuid=club,len=17"`
	Name         name         `gorm:"Type:varchar(30);NOT NULL;UNIQUE" validate:"min=1,max=30"`
	ClubConcept  clubConcept  `gorm:"Type:varchar(40)" validate:"max=40"`
	Introduction introduction `gorm:"Type:varchar(100)" validate:"max=150"`
	Field        field        `gorm:"Type:varchar(20);NOT NULL" validate:"min=1,max=20"`
	Location     location     `gorm:"Type:varchar(20);NOT NULL;UNIQUE" validate:"min=1,max=20"`
	Floor        floor        `gorm:"Type:tinyint(1);NOT NULL" validate:"range=1~5"`
	Link         link         `gorm:"Type:varchar(100)" validate:"max=1--"`
	LogoURI      logoURI      `gorm:"Type:varchar(100);NOT NULL" validate:"min=1,max=100"`
}
