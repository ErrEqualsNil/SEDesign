package db

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
)

type MySQlConf struct {
	Ip       string
	Port     string
	Protocol string
	User     string
	Password string
	Dbname   string
}

func GetMySQLConn() (*gorm.DB, error) {
	data, err := ioutil.ReadFile("conf/db_conf.yml")
	if err != nil {
		log.Fatalf("read db_conf err: %v", err)
		return nil, err
	}
	conf := MySQlConf{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("yaml Unmarshal err: %v", err)
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Protocol, conf.Ip, conf.Port, conf.Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db get connection err: %v", err)
		return nil, err
	}
	return db, nil
}