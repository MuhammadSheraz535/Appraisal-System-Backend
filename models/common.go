package models

import (
	"time"

	"gorm.io/gorm"
)

type CommonModel struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
