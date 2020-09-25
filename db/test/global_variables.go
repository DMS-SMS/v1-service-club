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

const numberOfTestFunc = 5

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

	clubMemberClubUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Club_Test_DB"),
		TableName:      model.ClubMemberInstance.TableName(),
		ConstraintName: model.ClubMemberInstance.ClubUUIDConstraintName(),
		AttrName:       model.ClubMemberInstance.ClubUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.ClubInstance.TableName(),
		AttrName:  model.ClubInstance.UUID.KeyName(),
	})

	clubRecruitmentClubUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Club_Test_DB"),
		TableName:      model.ClubRecruitmentInstance.TableName(),
		ConstraintName: model.ClubRecruitmentInstance.ClubUUIDConstraintName(),
		AttrName:       model.ClubRecruitmentInstance.ClubUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.ClubInstance.TableName(),
		AttrName:  model.ClubInstance.UUID.KeyName(),
	})

	recruitMemberRecruitmentUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Club_Test_DB"),
		TableName:      model.RecruitMemberInstance.TableName(),
		ConstraintName: model.RecruitMemberInstance.RecruitmentUUIDConstraintName(),
		AttrName:       model.RecruitMemberInstance.RecruitmentUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.ClubRecruitmentInstance.TableName(),
		AttrName:  model.ClubRecruitmentInstance.UUID.KeyName(),
	})
)