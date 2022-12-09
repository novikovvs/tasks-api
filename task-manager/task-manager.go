package main

import (
	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
	"net/http"
	"task-manager/http/controllers"
	_ "task-manager/http/controllers"

	"log"
	DBService "task-manager/db"
	model "task-manager/models"
)

func main() {
	var err error

	log.Println("Start app")

	DBService.CreateConn()

	err = createSchema()
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	s := r.Group("/tasks")
	{
		s.GET("/", controllers.TasksIndex)
		s.PUT("/create", controllers.TasksCreate)
	}

	r.Run()
}

func createSchema() error {
	err := DBService.CreateConn().AutoMigrate(&model.Task{})
	if err != nil {
		return err
	}
	return nil
}
