package method

import (
	"SEDesign/dal/db"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type MGetCommentByTaskIdReqParam struct {
	TaskId uint64 `form:"taskId"`
	Offset int  `form:"offset"`
	Limit int `form:"limit"`
}

type MGetCommentByTaskIdHandler struct {
	Ctx *gin.Context
	req MGetCommentByTaskIdReqParam
}

func (handler MGetCommentByTaskIdHandler) checkValid (req MGetCommentByTaskIdReqParam) bool {
	if req.TaskId == 0 {
		log.Printf("TaskId Not Found")
		return false
	}
	return true
}

func (handler MGetCommentByTaskIdHandler) Run () error {
	err := handler.Ctx.ShouldBind(&handler.req)
	if err != nil {
		log.Printf("Invalid err: %v", err)
		return errors.New("invalid req")
	}
	//参数检查
	if !handler.checkValid(handler.req) {
		log.Printf("Invalid params, req: %v", handler.req)
		return errors.New("invalid params")
	}
	comments, err := db.MGetCommentByTaskId(handler.req.TaskId, handler.req.Offset, handler.req.Limit)
	if err != nil {
		log.Printf("db call MGetCommentByTaskId err: %v", err)
		return err
	}
	handler.Ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"resp": "get comment by task id success!",
		"Comments": comments,
		"limit": handler.req.Limit,
		"offset": handler.req.Offset,
	})
	return nil
}
