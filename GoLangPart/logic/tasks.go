package logic

import (
	"SEDesign/dal/cache"
	"SEDesign/dal/db"
	"SEDesign/dal/es"
	"SEDesign/dal/mq"
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

	err = mq.DeleteTask(task.Id)
	if err != nil {
		log.Printf("call mq DeleteTask err: %v", err)
		return err
	}

	return nil
}

func MCreateTask(tasks []*model.Task) error {
	//批量写入mysql
	err := db.MCreateTask(tasks)
	if err != nil {
		log.Printf("db MCreateTask err: %v", err)
		return err
	}

	//写入redis
	for _, task := range tasks{
		err = cache.AddTaskToCache(task.ItemId)
		if err != nil {
			log.Printf("cache AddTaskToCache err: %v", err)
		}
	}

	//批量写入es
	err = es.AddTask(tasks)
	if err != nil {
		log.Printf("es add task err: %v", err)
		return err
	}

	return nil
}