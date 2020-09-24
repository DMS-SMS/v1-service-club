package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"strings"
)

type DBConfig struct {
	Dialect string `json:"dialect" validate:"required"`
	Host    string `json:"host" validate:"required"`
	Port 	int	   `json:"port" validate:"required"`
	User    string `json:"user" validate:"required"`
	DB		string `json:"db" validate:"required"`
}

func ConnectDBWithConsul(cli *api.Client, key string) (db *gorm.DB, conf DBConfig, err error) {
	kv, _, err := cli.KV().Get(key, nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("unable to get %s KV from consul, err: %v", key, err.Error()))
		return
	}

	if err = json.Unmarshal(kv.Value, &conf); err != nil {
		err = errors.New(fmt.Sprintf("error occurs while unmarshal KV value into struct, err: %v", err.Error()))
		return
	}

	if err = validator.New().Struct(&conf); err != nil {
		err = errors.New(fmt.Sprintf("invalid %s KV value, err: %v", key, err.Error()))
		return
	}

	conf.Dialect = strings.ToLower(conf.Dialect)
	switch conf.Dialect {
	case "mysql":
		db, err = connectToMysql(conf)
	default:
		err = errors.New(fmt.Sprintf("%s is not supported db in this service.", conf.Dialect))
	}
	return
}

func connectToMysql(conf DBConfig) (db *gorm.DB, err error) {
	pwd := os.Getenv("DB_PASSWORD")
	if pwd == "" {
		err = errors.New("please set DB_PASSWORD environment variable")
		return
	}
	args := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.User, pwd, conf.Host, conf.DB)
	db, err = gorm.Open(conf.Dialect, args)
	return
}