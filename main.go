package main

import (
	"club/adapter"
	"club/db"
	"club/db/access"
	"club/handler"
	authproto "club/proto/golang/auth"
	clubproto "club/proto/golang/club"
	"club/tool/closure"
	consulagent "club/tool/consul/agent"
	topic "club/utils/topic/golang"
	"fmt"
	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/checkers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	grpccli "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	// create service
	port := getRandomPortNotInUsedWithRange(10100, 10200)
	service := micro.NewService(
		micro.Name(topic.ClubServiceName),
		micro.Version("1.0.2"),
		micro.Transport(grpc.NewTransport()),
		micro.Address(fmt.Sprintf(":%d", port)),
	)
	srvID := fmt.Sprintf("%s-%s", service.Server().Options().Name, service.Server().Options().Id)

	// create consul connection
	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		log.Fatal("please set CONSUL_ADDRESS in environment variable")
	}
	consulCfg := api.DefaultConfig()
	consulCfg.Address = consulAddr
	consul, err := api.NewClient(consulCfg)
	if err != nil {
		log.Fatalf("consul connect fail, err: %v", err)
	}

	// create db access manager
	dbc, _, err := adapter.ConnectDBWithConsul(consul, "db/club/local")
	if err != nil {
		log.Fatalf("db connect fail, err: %v", err)
	}
	if err := db.Migrate(dbc); err != nil {
		log.Fatalf("db migration error, err: %v", err)
	}
	defaultAccessManage, err := db.NewAccessorManage(access.Default(dbc))
	if err != nil {
		log.Fatalf("db accessor create fail, err: %v", err)
	}

	// create jaeger connection
	jaegerAddr := os.Getenv("JAEGER_ADDRESS")
	if jaegerAddr == "" {
		log.Fatal("please set JAEGER_ADDRESS in environment variable")
	}
	authSrvTracer, closer, err := jaegercfg.Configuration{
		ServiceName: topic.ClubServiceName,
		Tags:        []opentracing.Tag{{"sid", srvID}},
		Reporter:    &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: jaegerAddr},
		Sampler:     &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1},
	}.NewTracer()
	if err != nil {
		log.Fatalf("error while creating new tracer for service, err: %v", err)
	}
	defer func() {
		_ = closer.Close()
	}()

	// create AWS session
	awsId := os.Getenv("SMS_AWS_ID")
	if awsId == "" {
		log.Fatal("please set SMS_AWS_ID in environment variable")
	}
	awsKey := os.Getenv("SMS_AWS_KEY")
	if awsKey == "" {
		log.Fatal("please set SMS_AWS_KEY in environment variable")
	}
	s3Region := os.Getenv("SMS_AWS_REGION")
	if s3Region == "" {
		log.Fatal("please set SMS_AWS_REGION in environment variable")
	}
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Region),
		Credentials: credentials.NewStaticCredentials(awsId, awsKey, ""),
	})
	if err != nil {
		log.Fatalf("error while creating new aws session, err: %v", err)
	}

	// create rpc handler
	defaultAgent := consulagent.Default(
		consulagent.Strategy(selector.RoundRobin),
		consulagent.Client(consul),
	)
	cliOpts := []client.Option{client.Transport(grpc.NewTransport())}
	authStudentSrv := authproto.NewAuthStudentService(topic.AuthServiceName, grpccli.NewClient(cliOpts...))
	rpcHandler := handler.Default(
		handler.AccessManager(defaultAccessManage),
		handler.Tracer(authSrvTracer),
		handler.ConsulAgent(defaultAgent),
		handler.AuthStudent(authStudentSrv),
		handler.AWSSession(awsSession),
	)

	service.Init(
		micro.AfterStart(closure.ConsulServiceRegistrar(service.Server(), consul)),
		micro.BeforeStop(closure.ConsulServiceDeregistrar(service.Server(), consul)),
	)

	_ = clubproto.RegisterClubAdminHandler(service.Server(), rpcHandler)
	_ = clubproto.RegisterClubStudentHandler(service.Server(), rpcHandler)
	_ = clubproto.RegisterClubLeaderHandler(service.Server(), rpcHandler)

	// DB Health checker 실행
	sqlDB, err := dbc.DB()
	if err != nil {
		log.Fatalf("unable to get sql DB from gorm DB, err: %v", err)
	}
	h := health.New()
	dbChecker, err := checkers.NewSQL(&checkers.SQLConfig{
		Pinger: sqlDB,
	})
	if err != nil {
		log.Fatalf("unable to create sql health checker, err: %v", err)
	}
	dbHealthCfg := &health.Config{
		Name:       "DB-Checker",
		Checker:    dbChecker,
		Interval:   time.Second * 5,
		OnComplete: closure.TTLCheckHandlerAboutDB(service.Server(), consul),
	}
	if err = h.AddChecks([]*health.Config{dbHealthCfg}); err != nil {
		log.Fatalf("unable to register health checks, err: %v", err)
	}
	if err = h.Start(); err != nil {
		log.Fatalf("unable to start health checks, err: %v", err)
	}

	// 서비스 실행
	if err := service.Run(); err != nil {
		log.Fatalf("error occurs while running service, err: %v", err)
	}
}

func getRandomPortNotInUsedWithRange(min, max int) (port int) {
	for {
		port = rand.Intn(max - min) + min
		conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}
		_ = conn.Close()
		break
	}
	return
}
