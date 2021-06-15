package method

import (
	"SEDesign/dal/db"
	"SEDesign/dal/mq"
	"SEDesign/model"
	"gorm.io/gorm"
	"log"
	"time"
)

func AbnormalTaskCheckEachHour() {
	for {
		time.Sleep(1 * time.Hour)
		err := AbnormalTaskCheckEachHourRun()
		if err != nil {
			log.Printf("AbnormalTaskCheckEachHourRun err: %v", err)
		}
	}
}

func AbnormalTaskCheckEachHourRun() error {
	tasks, err := db.MGetUnSubmitTask()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("No Task Need Resubmit")
			return nil
		}
		log.Printf("db MGetUnSubitTask err: %v", err)
		return err
	}

	for _, task := range tasks {
		err = mq.SubmitTask(task)
		if err != nil {
			log.Printf("mq SubmitTask err: %v, TaskId: %v", err, task.Id)
			return err
		}

		err = db.UpdateTaskStatus(task.Id, model.TaskStatusQueueing)
		if err !=nil {
			log.Printf("db UpdateTaskStatus err: %v, TaskId: %v", err, task.Id)
			return err
		}
	}
	return nil
}