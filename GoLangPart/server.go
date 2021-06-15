package main

import (
	"SEDesign/method"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main(){
	//定时检查未提交的任务
	go method.AbnormalTaskCheckEachHour()

	r := gin.Default()
	r.POST("/create_task", func(context *gin.Context){
		handler := method.CreateTaskHandler{
			Ctx: context,
		}
		err := handler.Run()
		if err != nil {
			log.Printf("Call CreateTask err: %v\n", err)
			context.JSON(http.StatusOK, gin.H{
				"status_code": http.StatusInternalServerError,
				"resp": err.Error(),
			})
		}
	})

	r.POST("/search_comment_by_task_id", func(context *gin.Context) {
		handler := method.MGetCommentByTaskIdHandler{
			Ctx: context,
		}
		err := handler.Run()
		if err != nil {
			log.Printf("Call MGetCommentByTaskId err: %v\n", err)
			context.JSON(http.StatusOK, gin.H{
				"status_code": http.StatusInternalServerError,
				"resp": err.Error(),
			})
		}
	})

	err := r.Run(":8000")
	if err != nil {
		return
	}
}