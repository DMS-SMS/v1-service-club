package main

import (
	"club/adapter"
	"club/db"
	"club/db/access"
	"club/handler"
	authproto "club/proto/golang/auth"
	consulagent "club/tool/consul/agent"
	topic "club/utils/topic/golang"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client/selector"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport/grpc"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"math/rand"
	"net"
	"os"
)

func main() {
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
		ServiceName: "DMS.SMS.v1.service.club",
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerAddr,
		},
	}.NewTracer()
	if err != nil {
		log.Fatalf("error while creating new tracer for service, err: %v", err)
	}
	defer func() {
		_ = closer.Close()
	}()

	// create service
	version := os.Getenv("VERSION")
	if version == "" {
		log.Fatal("please set VERSION in environment variable")
	}
	port := getRandomPortNotInUsedWithRange(10100, 10200)
	service := micro.NewService(
		micro.Name(topic.ClubServiceName),
		micro.Version(version),
		micro.Transport(grpc.NewTransport()),
		micro.Address(fmt.Sprintf(":%d", port)),
	)

	// create rpc handler
	defaultAgent := consulagent.Default(
		consulagent.Strategy(selector.RoundRobin),
		consulagent.Client(consul),
	)
	authStudentSrv := authproto.NewAuthStudentService(topic.AuthServiceName, service.Client())
	rpcHandler := handler.Default(
		handler.AWSSession(nil),
		handler.AccessManager(defaultAccessManage),
		handler.Tracer(authSrvTracer),
		handler.ConsulAgent(defaultAgent),
		handler.AuthStudent(authStudentSrv),
	)
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
