package test

import (
	clubproto "club/proto/golang/club"
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
	"log"
)

type GetClubsSortByUpdateTimeCase struct {
	UUID string
	Field, Name string
	Start, Count uint32
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedClubUUID  string
}

func (test *GetClubsSortByUpdateTimeCase) ChangeEmptyValueToValidValue() {
	if test.UUID == EmptyString              { test.UUID = validAdminUUID }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *GetClubsSortByUpdateTimeCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString                     { test.UUID = "" }
	if test.XRequestID == EmptyReplaceValueForString               { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString        { test.SpanContextString = "" }
}


func (test *GetClubsSortByUpdateTimeCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetClubsSortByUpdateTimeCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "GetClubInformsSortByUpdateTime":
		mock.On(string(method), test.Start, test.Count, test.Field, test.Name).Return(returns...)
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *GetClubsSortByUpdateTimeCase) SetRequestContextOf(req *clubproto.GetClubsSortByUpdateTimeRequest) {
	req.UUID = test.UUID
	req.Start = test.Start
	req.Count = test.Count
	req.Field = test.Field
	req.Name = test.Name

}

func (test *GetClubsSortByUpdateTimeCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	return
}
