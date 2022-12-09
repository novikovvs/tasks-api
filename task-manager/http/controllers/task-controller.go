package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"log"
	_ "net/http"
	DBService "task-manager/db"
	HttpUtility "task-manager/http"
	model "task-manager/models"
)

func TasksIndex(c *gin.Context) {
	var tasks []model.Task
	db := DBService.CreateConn()

	err := db.
		Model(&model.Task{}).
		Preload(clause.Associations, model.PreloadTasks).
		Where("parent_id is null").
		Find(&tasks).
		Error

	if err != nil {
		log.Println(err)
	}

	HttpUtility.Success(c, tasks, "", "")
}

func TasksCreate(c *gin.Context) {
	var task model.Task

	if err := c.ShouldBind(&task); err == nil {
		err = DBService.CreateConn().Create(&task).Error
		if err != nil {
			HttpUtility.Success(c, nil, "", err.Error())
		} else {
			HttpUtility.Success(c, task, "Created", "")
		}
	} else {
		HttpUtility.Success(c, nil, "", err.Error())
	}

}
