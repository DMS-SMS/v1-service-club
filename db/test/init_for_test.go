package test

import (
	"club/adapter"
	"club/db"
	"club/db/access"
	"github.com/hashicorp/consul/api"
	"log"
	"sync"
)

func init() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	dbc, _, err := adapter.ConnectDBWithConsul(cli, "db/club/local_test")
	if err != nil {
		log.Fatal(err)
	}
	db.Migrate(dbc)

	manager, err = db.NewAccessorManage(access.Default(dbc))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		testGroup = sync.WaitGroup{}
		testGroup.Add(numberOfTestFunc)
		testGroup.Wait()
		_ = dbc.Close()
	}()
}
