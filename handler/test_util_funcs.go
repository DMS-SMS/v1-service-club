package handler

import (
	"club/db"
	"club/db/access"
	authproto "club/proto/golang/auth"
	consulagent "club/tool/consul/agent"
	"fmt"
	"github.com/stretchr/testify/mock"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
)

func setAndGetTestEnv() (newMock *mock.Mock, h *_default) {
	newMock = new(mock.Mock)

	exampleTracerForRPCService, closer, err := jaegercfg.Configuration{ServiceName: "DMS.SMS.v1.service.auth"}.NewTracer()
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new tracer for service, err: %v", err)) }
	defer func() { _ = closer.Close() }()

	mockAccessManage, err := db.NewAccessorManage(access.Mock(newMock))
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new access manage with mock, err: %v", err)) }

	mockConsulAgent := consulagent.Mock(newMock)
	mockAuthStudent := authproto.MockAuthStudentService(newMock)

	h = Default(
		AccessManager(mockAccessManage),
		Tracer(exampleTracerForRPCService),
		ConsulAgent(mockConsulAgent),
		AuthStudent(mockAuthStudent),
	)
	return
}
