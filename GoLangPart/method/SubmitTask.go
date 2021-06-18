package method

import (
	"SEDesign/dal/db"
	"SEDesign/dal/mq"
	"SEDesign/model"
	"gorm.io/gorm"
	"log"
	"time"
)

func SubmitTaskEachTenMin() {
	for {
		time.Sleep(10 * time.Minute)
		err := SubmitTaskRun()
		if err != nil {
			log.Printf("SubmitTaskRun err: %v", err)
		}
	}
}

func SubmitTaskRun() error {
	tasks, err := db.MGetUnSubmitTask()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("No Task Need Resubmit")
			return nil
		}
		log.Printf("db MGetUnSubitTask err: %v", err)
		return err
	}

	taskIds := make([]int64, 0)
	for _, task := range tasks {
		taskIds = append(taskIds, task.Id)
		//提交到 message queue
		err = mq.SubmitTask(task)
		if err != nil {
			log.Printf("mq SubmitTask err: %v, TaskId: %v", err, task.Id)
			return err
		}
	}

	//更新状态
	err = db.UpdateTaskStatus(taskIds, model.TaskStatusQueueing)
	if err !=nil {
		log.Printf("db UpdateTaskStatus err: %v, TaskId: %v", err, taskIds)
		return err
	}
	return nil
}