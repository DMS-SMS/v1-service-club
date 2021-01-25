package test

import (
	"club/db"
	"club/db/access"
	"github.com/hashicorp/consul/api"
	"log"
)

func init() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	dbc, _, err := db.ConnectDBWithConsul(cli, "db/club/local_test")
	if err != nil {
		log.Fatal(err)
	}
	db.Migrate(dbc)

	manager, err = db.NewAccessorManage(access.Default(dbc))
	if err != nil {
		log.Fatal(err)
	}
}
