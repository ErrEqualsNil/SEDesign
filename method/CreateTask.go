package method

import (
	"SEDesign/dal/db"
	"SEDesign/dal/mq"
	"SEDesign/model"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type CreateTaskReqParam struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type CreateTaskHandler struct {
	Ctx *gin.Context
	req CreateTaskReqParam
}

func (handler CreateTaskHandler) checkValid (req CreateTaskReqParam) bool {
	if len(req.Name) == 0 {
		log.Printf("Name Not Found")
		return false
	}
	return true
}

func (handler CreateTaskHandler) Run () error {
	if handler.Ctx.ShouldBind(&handler.req) != nil {
		log.Printf("Invalid req: %v", handler.Ctx.Params)
		return errors.New("invalid req")
	}
	//参数检查
	if !handler.checkValid(handler.req) {
		log.Printf("Invalid params")
		return errors.New("invalid params")
	}

	//写入mysql
	task := &model.Task{
		ItemName: handler.req.Name,
		Status: model.TaskStatusCreating,
	}
	err := db.CreateTask(task)
	if err != nil {
		log.Printf("db create comment err: %v", err)
		return err
	}

	//写入mq
	err = mq.SubmitTask(task)
	if err != nil {
		log.Printf("mq submit task err: %v", err)
		return err
	}

	//更新排队状态
	err = db.UpdateTaskStatus(task.Id, model.TaskStatusQueueing)
	if err != nil {
		log.Printf("db update task err: %v", err)
		return err
	}
	handler.Ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"resp": "create task success!",
	})
	return nil
}
