package test

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
	"os"
	"path/filepath"
)

const (
	validAdminUUID = "admin-111111111111"
	validLeaderUUID = "student-111111111111"
	validClubUUID = "club-111111111111"

	validClubName = "DMS"
	validField = "SW 개발"
	validLocation = "2-2반 교실"
	validFloor = "3"
)

var (
	validImageByteArr []byte
	validSpanContextString string
	validXRequestID string
	validMemberUUIDs = []string{validLeaderUUID}
)

func init() {
	exampleTracerForAPIGateway, closer, err := jaegercfg.Configuration{ServiceName: "DMS.SMS.v1.api.gateway"}.NewTracer()
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new tracer for api, err: %v", err)) }
	defer func() { _ = closer.Close() }()
	exampleTracerForRPCService, closer, err := jaegercfg.Configuration{ServiceName: "DMS.SMS.v1.service.club"}.NewTracer()
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new tracer for service, err: %v", err)) }
	defer func() { _ = closer.Close() }()

	exampleSpanForAPIGateway := exampleTracerForAPIGateway.StartSpan("v1/clubs")
	exampleSpanForRPCService := exampleTracerForRPCService.StartSpan("CreateNewClub", opentracing.ChildOf(exampleSpanForAPIGateway.Context()))
	validSpanContextString = exampleSpanForRPCService.Context().(jaeger.SpanContext).String()

	absPath, err := filepath.Abs("./for_test/image/doraemon.png")
	if err != nil { log.Fatal(fmt.Sprintf("error while getting abstract file path, err: %v", err)) }
	file, err := os.Open(absPath)
	if err != nil { log.Fatal(fmt.Sprintf("error while opening new test image files, err: %v", err)) }
	fileInfo, err := file.Stat()
	if err != nil { log.Fatal(fmt.Sprintf("error while getting file information, err: %v", err)) }
	validImageByteArr = make([]byte, fileInfo.Size())
	fileReader := bufio.NewReader(file)
	_, err = fileReader.Read(validImageByteArr)
	if err != nil { log.Fatal(fmt.Sprintf("error while reading from image file, err: %v", err)) }

	validXRequestID = uuid.New().String()
}
