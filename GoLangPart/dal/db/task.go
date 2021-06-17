package db

import (
	"SEDesign/model"
	"gorm.io/gorm"
	"log"
)

func MGetTasks(ids []int64) ([]*model.Task, error) {
	conn, err := GetMySQLConn()
	if err != nil {
		log.Fatalf("call db GetMySQLConn err: %v", err)
		return nil, err
	}

	result := make([]*model.Task, 0)
	err = conn.Model(model.Task{}).Where("id in (?)", ids).Find(&result).Error
	if err != nil {
		log.Printf("error to find task from db, err: %v, ids: %v", err, ids)
		return nil, err
	}
	if len(result) == 0 {
		log.Printf("Tasks Not Found, ids: %v", ids)
		return nil, gorm.ErrRecordNotFound
	}
	return result, nil
}

func CreateTask(task *model.Task) error {
	conn, err := GetMySQLConn()
	if err != nil {
		log.Fatalf("call db GetMySQLConn err: %v", err)
		return err
	}

	err = conn.Model(model.Task{}).Create(task).Error
	if err != nil {
		log.Printf("error to create task at db, err: %v, task: %v", err, task)
		return err
	}
	return nil
}

func UpdateTaskStatus(id int64, newStatus model.TaskStatus) error {
	conn, err := GetMySQLConn()
	if err != nil {
		log.Fatalf("call db GetMySQLConn err: %v", err)
		return err
	}

	err = conn.Model(model.Task{}).Where("id=?", id).Update("status", newStatus).Error
	if err != nil {
		log.Printf("error to update task status at db, err: %v, taskId: %v", err, id)
		return err
	}
	return nil
}

func MGetUnSubmitTask() ([]*model.Task, error) {
	conn, err := GetMySQLConn()
	if err != nil {
		log.Fatalf("call db GetMySQLConn err: %v", err)
		return nil, err
	}

	result := make([]*model.Task, 0)
	err = conn.Model(model.Task{}).Where("status in (?)", []model.TaskStatus{model.TaskStatusCreating, model.TaskStatusUnknown}).Find(&result).Error
	if err != nil {
		log.Printf("error to find task (unsubmited), err:%v", err)
		return nil, err
	}
	if len(result) == 0 {
		log.Printf("unsubmited task not found")
		return nil, gorm.ErrRecordNotFound
	}
	return result, nil
}