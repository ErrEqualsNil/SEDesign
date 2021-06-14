package method

import (
	"SEDesign/dal/db"
	"SEDesign/dal/mq"
	"SEDesign/model"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
)

type ReqParam struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func CheckValid(req ReqParam) bool {
	if len(req.Name) == 0 {
		log.Printf("Name Not Found")
		return false
	}
	return true
}

func CreateTask(c *gin.Context) error {
	req := ReqParam{}
	if c.ShouldBind(&req) != nil {
		log.Printf("Invalid req: %v", c.Params)
		return errors.New("invalid req")
	}
	//参数检查
	if !CheckValid(req) {
		log.Printf("Invalid params")
		return errors.New("invalid params")
	}

	task := &model.Comment{
		ItemName: req.Name,
		Status: model.TaskStatusQueueing,
	}
	err := db.CreateComment(task)
	if err != nil {
		log.Printf("db create comment err: %v", err)
		return err
	}

	param := &mq.MqTaskParam{
		Id: task.Id,
		Name: task.ItemName,
	}
	err = mq.SubmitTask(param)
	if err != nil {
		log.Printf("mq submit task err: %v", err)
		return err
	}
	return nil
}
