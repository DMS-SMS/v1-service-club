package model

import (
	"gorm.io/gorm"
	"reflect"
	"time"
)

// 매개변수로 전달 받은 변수의 DeepCopy 본사본을 생성하여 반환하는 함수
// 제약조건 -> 매개변수로 넘길 변수는 포인터 변수여야 함!! (X -> panic 발생)
func deepCopyModel(model interface{}) interface{} {
	duplicateModel := reflect.New(reflect.ValueOf(model).Elem().Type())
	duplicateModel.Elem().Set(reflect.ValueOf(model).Elem())
	return duplicateModel.Interface()
}

// deepCopyModel 함수를 이용하여 본사본 생성 후 gorm.Model 필드 값 초기화 후 해당 변수 반환 함수
// 제약조건 -> 매개변수로 넘길 변수는 구조체인 동시에 gorm.Model 객체의 필드들을 가지고 있어야 함!! (X -> panic 발생)
func exceptGormModel(model interface{}) (gormModelExceptTable interface{}) {
	gormModelExceptTable = deepCopyModel(model)

	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("ID").Set(reflect.ValueOf(uint(0)))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(time.Time{}))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Time{}))
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("DeletedAt").Set(reflect.ValueOf(gorm.DeletedAt{}))
	return
}

// DeepCopy 메서드 -> 리시버 변수에 대한 DeepCopy 본사본 생성 및 반환 메서드
func (c *Club)             DeepCopy() *Club            { return deepCopyModel(c).(*Club) }
func (ci *ClubInform)      DeepCopy() *ClubInform      { return deepCopyModel(ci).(*ClubInform) }
func (cm *ClubMember)      DeepCopy() *ClubMember      { return deepCopyModel(cm).(*ClubMember) }
func (cr *ClubRecruitment) DeepCopy() *ClubRecruitment { return deepCopyModel(cr).(*ClubRecruitment) }
func (rm *RecruitMember)   DeepCopy() *RecruitMember   { return deepCopyModel(rm).(*RecruitMember) }

// ExceptGormModel 메서드 -> 리시버 변수로부터 gorm.Model(임베딩 객체)에 포함되어있는 필드 값 초기화 후 반환 메서드
func (c *Club)             ExceptGormModel() *Club            { return exceptGormModel(c).(*Club) }
func (ci *ClubInform)      ExceptGormModel() *ClubInform      { return exceptGormModel(ci).(*ClubInform) }
func (cm *ClubMember)      ExceptGormModel() *ClubMember      { return exceptGormModel(cm).(*ClubMember) }
func (cr *ClubRecruitment) ExceptGormModel() *ClubRecruitment { return exceptGormModel(cr).(*ClubRecruitment) }
func (rm *RecruitMember)   ExceptGormModel() *RecruitMember   { return exceptGormModel(rm).(*RecruitMember) }

// XXXConstraintName 메서드 -> XXX PK의 Constraint Name 값 반환 메서드
func (ci *ClubInform)      ClubUUIDConstraintName()        string { return "fk_club_informs_club" }
func (cm *ClubMember)      ClubUUIDConstraintName()        string { return "fk_club_members_club" }
func (cr *ClubRecruitment) ClubUUIDConstraintName()        string { return "fk_club_recruitments_club" }
func (rm *RecruitMember)   RecruitmentUUIDConstraintName() string { return "fk_recruit_members_club" }

// TableName 메서드 -> 리시버 변수에 해당되는 테이블의 이름 반환 메서드
func (c *Club)             TableName() string { return "clubs" }
func (ci *ClubInform)      TableName() string { return "club_informs" }
func (cm *ClubMember)      TableName() string { return "club_members" }
func (cr *ClubRecruitment) TableName() string { return "club_recruitments" }
func (rm *RecruitMember)   TableName() string { return "recruit_members" }
