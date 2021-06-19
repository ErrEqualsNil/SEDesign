package method

import (
	"SEDesign/dal/db"
	"SEDesign/dal/es"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SearchTaskReqParam struct {
	TaskId int64 `form:"taskId"`
	ItemId int64 `form:"itemId"`
	Name string `form:"name"`
	Offset int  `form:"offset"`
	Limit int `form:"limit"`
}

type SearchTaskHandler struct {
	Ctx *gin.Context
	req SearchTaskReqParam
}

func (handler SearchTaskHandler) checkValid () bool {
	if handler.req.Offset < 0 || handler.req.Limit < 0{
		log.Printf("Invalid Param: %v", handler.req)
		return false
	}
	return true
}

func (handler SearchTaskHandler) Run () error {
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

	param := &es.SearchTaskReqParams{
		Id:   handler.req.TaskId,
		Name: handler.req.Name,
		ItemId: handler.req.ItemId,
		Limit: handler.req.Limit,
		Offset: handler.req.Offset,
	}
	totalCount, taskIds, err := es.SearchTaskByName(param)
	if err != nil {
		log.Printf("es SearchTaskByName err: %v", err)
		return err
	}
	tasks, err := db.MGetTasks(taskIds)
	if err != nil {
		log.Printf("db MGetTasks err: %v", err)
		return err
	}
	handler.Ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"resp": "search task success!",
		"total_count": totalCount,
		"tasks": tasks,
		"limit": handler.req.Limit,
		"offset": handler.req.Offset,
	})
	return nil
}
