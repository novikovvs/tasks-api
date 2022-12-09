package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Task struct {
	gorm.Model
	Title       string    `json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"-"`
	SubTasks    []Task    `gorm:"foreignKey:ParentID;references:ID"  json:"sub_tasks"`
	Done        bool      `json:"done,omitempty" form:"done" binding:"-"`
	ParentID    NullInt64 `json:"parent_id,string,omitempty" form:"parent_id,string" binding:"-"`
}

type NullInt64 struct {
	sql.NullInt64
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(nil)
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

func (u Task) String() string {
	return fmt.Sprintf("Task<%d %s %s>", u.ID, u.Title, u.Description)
}

func PreloadTasks(d *gorm.DB) *gorm.DB {
	return d.Preload("SubTasks."+clause.Associations, PreloadTasks)
}
