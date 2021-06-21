package method

import (
	"SEDesign/dal/db"
	"SEDesign/logic"
	"log"
	"time"
)

func CleanTask() error{
	var err error
	err = nil

	tasks, err := db.MGetAbnormalTask()
	if err != nil {
		log.Printf("db MGetAbnormalTask err: %v", err)
		return err
	}

	for _, task := range tasks {
		err := logic.DeleteTask(task)
		if err != nil {
			log.Printf("logic Delete Task err: %v, taskId: %v", err, task.Id)
			continue
		}
	}
	return err
}

func CleanTaskEachHour() {
	for {
		err := CleanTask()
		if err != nil {
			log.Printf("clean task err: %v", err)
		}
		time.Sleep(time.Hour * 1)
	}
}
