package method

import (
	"SEDesign/dal/db"
	"SEDesign/logic"
	"SEDesign/model"
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

	newTasks := make([]*model.Task, 0)
	for _, task := range tasks {
		//非商品只有100条评论的情况下， 保留task内容，便于重新生成
		if !(task.CommentCount <= 100 && task.Status == model.TaskStatusComplete) {
			newTask := &model.Task{
				ItemName: task.ItemName,
				ItemId: task.ItemId,
				Status: model.TaskStatusCreating,
				CommentCount: 0,
			}
			newTasks = append(newTasks, newTask)
		}

		err = logic.DeleteTask(task)
		if err != nil {
			log.Printf("logic Delete Task err: %v, taskId: %v", err, task.Id)
			continue
		}
	}

	err = logic.MCreateTask(newTasks)
	if err != nil {
		log.Printf("logic MCreateTask err: %v", err)
		return err
	}

	return err
}

func CleanTaskEachHour() {
	for {
	    time.Sleep(time.Hour * 1)
		err := CleanTask()
		if err != nil {
			log.Printf("clean task err: %v", err)
		}
	}
}
