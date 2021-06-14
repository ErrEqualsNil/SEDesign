package main

import (
	"SEDesign/method"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main(){
	r := gin.Default()
	r.POST("/create", func(c *gin.Context){
		err := method.CreateTask(c)
		if err != nil {
			log.Printf("Call create task err: %v\n", err)
			c.JSON(http.StatusOK, gin.H{
				"resp": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"resp": "create task success!",
			})
		}
	})
	err := r.Run(":8000")
	if err != nil {
		return
	}
}