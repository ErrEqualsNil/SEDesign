package method

import (
	"SEDesign/dal/db"
	"SEDesign/logic"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
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

	//检查task是否存在
	taskList, err := db.MGetTasks([]int64{handler.req.TaskId})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("delete task not found, taskid: %v", handler.req.TaskId)
		}
		log.Printf("call MGetTasks err: %v", err)
		return err
	}

	task := taskList[0]
	err = logic.DeleteTask(task)
	if err != nil {
		log.Printf("logic Delete Task err: %v, taskId: %v", err, task.Id)
		return err
	}

	handler.Ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"resp": "Delete Task success!",
	})
	return nil
}