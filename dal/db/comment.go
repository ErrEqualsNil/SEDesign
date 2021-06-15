package db

import (
	"SEDesign/model"
	"gorm.io/gorm"
	"log"
)

func MGetComment(ids []uint64) ([]*model.Comment, error) {
	db, err := GetConn()
	if err != nil {
		log.Printf("call db GetConn err: %v", err)
		return nil, err
	}
	result := make([]*model.Comment, 0)
	err = db.Model(model.Comment{}).Where("id in (?)", ids).Find(&result).Error
	if err != nil {
		log.Printf("error to find comment at db, err: %v, ids: %v", err, ids)
		return nil, err
	}
	if len(result) == 0 {
		log.Printf("Comment Not Found, ids: %v", ids)
		return nil, gorm.ErrRecordNotFound
	}
	return result, nil
}

func MGetCommentByTaskId(taskId uint64, offset int, limit int) ([]*model.Comment, error) {
	db, err := GetConn()
	if err != nil {
		log.Printf("call db GetConn err: %v", err)
		return nil, err
	}
	result := make([]*model.Comment, 0)
	err = db.Model(model.Comment{}).Where("task_id=?", taskId).Limit(limit).Offset(offset).Find(&result).Error
	if err != nil {
		log.Printf("error to find comment by task id, err: %v", err)
		return nil, err
	}
	if len(result) == 0 {
		log.Printf("Comment Not Found, taskId: %v", taskId)
		return nil, gorm.ErrRecordNotFound
	}
	return result, nil
}