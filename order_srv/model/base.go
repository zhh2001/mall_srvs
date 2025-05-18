package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GormList []string

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

func (g *GormList) Value() (driver.Value, error) {
	return json.Marshal(*g)
}

type BaseModel struct {
	ID        int32          `json:"id" gorm:"primary_key;type:int"`
	IsDeleted bool           `json:"-" gorm:"column:is_deleted"`
	CreatedAt time.Time      `json:"-" gorm:"column:create_time"`
	UpdatedAt time.Time      `json:"-" gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
