package validate

import (
	"github.com/go-playground/validator/v10"
	"log"
	"strconv"
	"strings"
)

var DBValidator *validator.Validate

func init() {
	DBValidator = validator.New()

	if err := DBValidator.RegisterValidation("uuid", isValidateUUID); err != nil { log.Fatal(err) } // 문자열 전용
	if err := DBValidator.RegisterValidation("time", isTime);         err != nil { log.Fatal(err) } // 문자열 전용
	if err := DBValidator.RegisterValidation("range", isWithinRange); err != nil { log.Fatal(err) } // 정수 전용
}

func isValidateUUID(fl validator.FieldLevel) bool {
	switch fl.Param() {
	case "admin":
		return adminUUIDRegex.MatchString(fl.Field().String())
	case "student":
		return studentUUIDRegex.MatchString(fl.Field().String())
	case "teacher":
		return teacherUUIDRegex.MatchString(fl.Field().String())
	case "parent":
		return parentUUIDRegex.MatchString(fl.Field().String())
	case "club":
		return clubUUIDRegex.MatchString(fl.Field().String())
	case "recruitment":
		return recruitmentUUIDRegex.MatchString(fl.Field().String())
	}
	return false
}

func isTime(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}

	if len(field) != 10 {
		return false
	}

	return timeRegex.MatchString(fl.Field().String())
}

func isWithinRange(fl validator.FieldLevel) bool {
	_range := strings.Split(fl.Param(), "~")
	if len(_range) != 2 {
		log.Fatal("please set param of range like (int)~(int)")
	}

	start, err := strconv.Atoi(_range[0])
	if err != nil {
		log.Fatalf("please set param of range like (int)~(int), err: %v", err)
	}
	end, err := strconv.Atoi(_range[1])
	if err != nil {
		log.Fatalf("please set param of range like (int)~(int), err: %v", err)
	}

	field := int(fl.Field().Int())
	return field >= start && field <= end
}