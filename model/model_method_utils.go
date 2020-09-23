package model

import (
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
	reflect.ValueOf(gormModelExceptTable).Elem().FieldByName("DeletedAt").Set(reflect.ValueOf((*time.Time)(nil)))
	return
}
