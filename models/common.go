package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CommonModel struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Enum interface {
	IsValid() bool
}

func ValidateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Enum)
	return value.IsValid()
}
