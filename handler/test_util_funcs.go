package handler

import (
	consulagent "club/consul/agent"
	"club/db"
	"club/db/access"
	authproto "club/proto/golang/auth"
	"fmt"
	"github.com/stretchr/testify/mock"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
)

func newDefaultMockHandler(mock *mock.Mock) *_default {
	exampleTracerForRPCService, closer, err := jaegercfg.Configuration{ServiceName: "DMS.SMS.v1.service.auth"}.NewTracer()
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new tracer for service, err: %v", err)) }
	defer func() { _ = closer.Close() }()

	mockAccessManage, err := db.NewAccessorManage(access.Mock(mock))
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new access manage with mock, err: %v", err)) }

	mockConsulAgent := consulagent.Mock(mock)
	mockAuthStudent := authproto.MockAuthStudentService(mock)

	return Default(
		AccessManager(mockAccessManage),
		Tracer(exampleTracerForRPCService),
		ConsulAgent(mockConsulAgent),
		AuthStudent(mockAuthStudent),
	)
}
