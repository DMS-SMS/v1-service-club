package main

import (
	"club/adapter"
	"club/db"
	"club/db/access"
	"github.com/hashicorp/consul/api"
	log "github.com/micro/go-micro/v2/logger"
	jaegercfg "github.com/uber/jaeger-client-go/config"
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
}

