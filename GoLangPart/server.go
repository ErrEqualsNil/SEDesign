package main

import (
	"SEDesign/method"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main(){
	//定时检查未提交的任务
	go method.SubmitTaskEachHour()
	//定时清理异常任务
	go method.CleanTaskEachHour()

	r := gin.Default()
	r.POST("/create_task", func(context *gin.Context){
		handler := method.CreateTaskHandler{
			Ctx: context,
		}
		err := handler.Run()
		if err != nil {
			log.Printf("Call CreateTask err: %v", err)
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
			log.Printf("Call MGetCommentByTaskId err: %v", err)
			context.JSON(http.StatusOK, gin.H{
				"status_code": http.StatusInternalServerError,
				"resp": err.Error(),
			})
		}
	})

	r.POST("/search_task", func(context *gin.Context) {
		handler := method.SearchTaskHandler{
			Ctx: context,
		}
		err := handler.Run()
		if err != nil {
			log.Printf("Call SearchTask err: %v", err)
			context.JSON(http.StatusOK, gin.H{
				"status_code": http.StatusInternalServerError,
				"resp": err.Error(),
			})
		}
	})

	r.POST("/submit_task", func(context *gin.Context) {
		err := method.SubmitTaskRun()
		if err != nil {
			log.Printf("Call SubmitTaskRun err: %v", err)
			context.JSON(http.StatusOK, gin.H{
				"status_code": http.StatusInternalServerError,
				"resp":        err.Error(),
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"status_code": http.StatusOK,
				"resp": "submit task success!",
			})
		}
	})

	r.POST("/delete_task", func(context *gin.Context) {
		handler := method.DeleteTaskHandler{
			Ctx: context,
		}
		err := handler.Run()
		if err != nil {
			log.Printf("Call DeleteTask err: %v", err)
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