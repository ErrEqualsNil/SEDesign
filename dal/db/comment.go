package db

import (
	"SEDesign/model"
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
)

type MySQlConf struct {
	ip       string
	port     string
	protocol string
	user     string
	password string
	dbName   string
}

func getConnection() (*gorm.DB, error) {
	data, err := ioutil.ReadFile("model/db_conf.yml")
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
		conf.user, conf.password, conf.protocol, conf.ip, conf.port, conf.dbName)
	log.Printf("db conf: %v", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db get connection err: %v", err)
		return nil, err
	}
	return db, nil
}

func MGetComment(ids []uint64) ([]*model.Comment, error) {
	db, err := getConnection()
	if err != nil {
		log.Fatalf("call db getConnection err: %v", err)
		return nil, err
	}
	result := make([]*model.Comment, 0)
	err = db.Model(model.Comment{}).Where("id in (?)", ids).Find(&result).Error
	if err != nil {
		log.Fatalf("error to find from db, err: %v, ids: %v", err, ids)
		return nil, err
	}
	return result, nil
}

func UpdateComment(id uint64, params map[string]interface{}) error {
	db, err := getConnection()
	if err != nil {
		log.Fatalf("call db getConnection err: %v", err)
		return err
	}

	err = db.Model(model.Comment{}).Where("id = ?", id).Updates(params).Error
	if err != nil {
		log.Fatalf("error to update from db, err: %v, ids: %v", err, id)
		return err
	}

	return nil
}

func CreateComment(comment *model.Comment) error {
	db, err := getConnection()
	if err != nil {
		log.Fatalf("call db getConnection err: %v", err)
		return err
	}

	err = db.Create(comment).Error
	if err != nil {
		log.Fatalf("error to create from db, err: %v, comment: %v", err, comment)
		return err
	}
	return nil
}