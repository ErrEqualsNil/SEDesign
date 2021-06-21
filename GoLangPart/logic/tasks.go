package logic

import (
	"SEDesign/dal/cache"
	"SEDesign/dal/db"
	"SEDesign/dal/es"
	"SEDesign/model"
	"log"
)

func DeleteTask(task *model.Task) error {
	//删除cache缓存
	err := cache.DeleteTaskById(task.ItemId)
	if err != nil {
		log.Printf("call cache DeleteTaskById err: %v", err)
		return err
	}

	err = es.DeleteTaskById(task.Id)
	if err != nil {
		log.Printf("call es DeleteTaskById err: %v", err)
		return err
	}

	err = db.DeleteTask(task.Id)
	if err != nil {
		log.Printf("call db DeleteTask err: %v", err)
		return err
	}

	return nil
}