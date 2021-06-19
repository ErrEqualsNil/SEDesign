package method

import (
	"SEDesign/dal/cache"
	"SEDesign/dal/db"
	"SEDesign/dal/es"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
)

type DeleteTaskReqParam struct {
	TaskId int64 `form:"taskId" json:"taskId"`
}

type DeleteTaskHandler struct {
	Ctx *gin.Context
	req DeleteTaskReqParam
}

func (handler DeleteTaskHandler) checkValid() bool {
	if handler.req.TaskId == 0 {
		log.Printf("req Illegal, req: %v", handler.req)
		return false
	}
	return true
}

func (handler DeleteTaskHandler) Run() error {
	err := handler.Ctx.ShouldBind(&handler.req)
	if err != nil {
		log.Printf("Invalid err: %v", err)
		return errors.New("invalid req")
	}
	//参数检查
	if !handler.checkValid() {
		log.Printf("Invalid params")
		return errors.New("invalid params")
	}

	err = cache.DeleteTaskById(handler.req.TaskId)
	if err != nil {
		log.Printf("call cache DeleteTaskById err: %v", err)
		return err
	}

	err = es.DeleteTaskById(handler.req.TaskId)
	if err != nil {
		log.Printf("call es DeleteTaskById err: %v", err)
		return err
	}

	err = db.DeleteTask(handler.req.TaskId)
	if err != nil {
		log.Printf("call db DeleteTask err: %v", err)
		return err
	}

	return nil
}