package test

import (
	"club/db"
	"club/model"
	"club/tool/mysqlerr"
	"strings"
	"sync"
)

var (
	manager db.AccessorManage
	testGroup sync.WaitGroup
)

const numberOfTestFunc = 2

var (
	clubInformClubUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Club_Test_DB"),
		TableName:      model.ClubInformInstance.TableName(),
		ConstraintName: model.ClubInformInstance.ClubUUIDConstraintName(),
		AttrName:       model.ClubInformInstance.ClubUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.ClubInstance.TableName(),
		AttrName:  model.ClubInstance.UUID.KeyName(),
	})
)