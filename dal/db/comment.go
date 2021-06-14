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
	Ip       string
	Port     string
	Protocol string
	User     string
	Password string
	Dbname   string
}

func GetConnection() (*gorm.DB, error) {
	data, err := ioutil.ReadFile("conf/db_conf.yml")
	if err != nil {
		log.Printf("read db_conf err: %v", err)
		return nil, err
	}
	conf := MySQlConf{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Printf("yaml Unmarshal err: %v", err)
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Protocol, conf.Ip, conf.Port, conf.Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("dsn: %v", dsn)
		log.Printf("db get connection err: %v", err)
		return nil, err
	}
	return db, nil
}

func MGetComment(ids []uint64) ([]*model.Comment, error) {
	db, err := GetConnection()
	if err != nil {
		log.Printf("call db GetConnection err: %v", err)
		return nil, err
	}
	result := make([]*model.Comment, 0)
	err = db.Model(model.Comment{}).Where("id in (?)", ids).Find(&result).Error
	if err != nil {
		log.Printf("error to find from db, err: %v, ids: %v", err, ids)
		return nil, err
	}
	if len(result) == 0 {
		log.Printf("Comment Not Found, ids: %v", ids)
		return nil, gorm.ErrRecordNotFound
	}
	return result, nil
}

func UpdateComment(id uint64, params map[string]interface{}) error {
	db, err := GetConnection()
	if err != nil {
		log.Printf("call db GetConnection err: %v", err)
		return err
	}

	err = db.Model(model.Comment{}).Where("id = ?", id).Updates(params).Error
	if err != nil {
		log.Printf("error to update from db, err: %v, ids: %v", err, id)
		return err
	}

	return nil
}

func CreateComment(comment *model.Comment) error {
	db, err := GetConnection()
	if err != nil {
		log.Printf("call db GetConnection err: %v", err)
		return err
	}

	err = db.Create(comment).Error
	if err != nil {
		log.Printf("error to create from db, err: %v, comment: %v", err, comment)
		return err
	}
	return nil
}