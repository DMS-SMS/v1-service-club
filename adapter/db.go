package adapter

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

type DBConfig struct {
	Dialect string `json:"dialect" validate:"required"`
	Host    string `json:"host" validate:"required"`
	Port 	int	   `json:"port" validate:"required"`
	User    string `json:"user" validate:"required"`
	DB		string `json:"db" validate:"required"`
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