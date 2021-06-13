package main

import (
	"SEDesign/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main(){
	dsn := "pengpeng:hdeilm1718@tcp(121.36.31.113:3306)/sedesign?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("mysql err: ", err)
	}
	db.AutoMigrate(&model.Comment{})
}