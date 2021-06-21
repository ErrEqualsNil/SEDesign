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

func MCreateTask(tasks []*model.Task) error {
	conn, err := GetMySQLConn()
	if err != nil {
		log.Fatalf("call db GetMySQLConn err: %v", err)
		return err
	}

	err = conn.Model(model.Task{}).Create(tasks).Error
	if err != nil {
		log.Printf("error to create task at db, err: %v, tasks: %v", err, tasks)
		return err
	}
	return nil
}

func DeleteTask(taskId int64) error {
	conn, err := GetMySQLConn()
	if err != nil {
		log.Fatalf("call db GetMySQLConn err: %v", err)
		return err
	}

	err = conn.Model(model.Task{}).Where("id=?", taskId).Delete(&model.Task{}).Error
	if err != nil {
		log.Printf("error to delete task at db, err: %v, taskId: %v", err, taskId)
		return err
	}

	err = conn.Model(model.Comment{}).Where("task_id=?", taskId).Delete(&model.Comment{}).Error
	if err != nil {
		log.Printf("error to delete comments of task, err: %v, taskId: %v", err, taskId)
		return err
	}
	return nil
}

func UpdateTaskStatus(taskIds []int64, newStatus model.TaskStatus) error {
	conn, err := GetMySQLConn()
	if err != nil {
		log.Fatalf("call db GetMySQLConn err: %v", err)
		return err
	}

	err = conn.Model(model.Task{}).Where("id in (?)", taskIds).Updates(map[string]interface{}{"status": newStatus}).Error
	if err != nil {
		log.Printf("error to update task status at db, err: %v, taskIds: %v", err, taskIds)
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

func MGetAbnormalTask() ([]*model.Task, error) {
	conn, err := GetMySQLConn()
	if err != nil {
		log.Fatalf("call db GetMySQLConn err: %v", err)
		return nil, err
	}

	result := make([]*model.Task, 0)
	err = conn.Model(model.Task{}).
		Where("status in (?)", []model.TaskStatus{model.TaskStatusProcessing, model.TaskStatusComplete}).
		Where(
			conn.Where("comment_count<100").Or("good_rate=0"),
			).
		Find(&result).Error
	if err != nil {
		log.Printf("MGetAbnormalTask err: %v", err)
		return nil, err
	}
	return result, nil
}