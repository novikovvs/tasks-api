package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Task struct {
	gorm.Model
	Title       string `json:"title" json:"title,omitempty"`
	Description string `json:"description" json:"description,omitempty"`
	SubTasks    []Task `json:"sub_tasks" gorm:"foreignKey:ParentID;references:ID" json:"sub_tasks"`
	Done        bool   `json:"done" json:"done,omitempty"`
	ParentID    uint   `json:"parent_id" json:"parent_id,omitempty"`
}

func (u Task) String() string {
	return fmt.Sprintf("Task<%d %s %s>", u.ID, u.Title, u.Description)
}

func PreloadTasks(d *gorm.DB) *gorm.DB {
	return d.Preload("SubTasks."+clause.Associations, PreloadTasks)
}
