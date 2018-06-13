package models

import (
	"github.com/jinzhu/gorm"
	"github.com/go-ozzo/ozzo-validation"
)

type Resource struct {
	gorm.Model
	Name string `json:"name"`
	Path string `json:"path"`
}

type ResourcePermission struct {
	gorm.Model
	//ID uint `gorm:"primary_key"`
	Resource     Resource   `gorm:"foreignKey:ResourceID"`
	ResourceID   uint
	Permission   Permission `gorm:"foreignKey:PermissionID"`
	PermissionID uint
}

func (Resource) TableName() string {
	return "rw_resources"
}

func (ResourcePermission) TableName() string {
	return "rw_resources_permissions"
}

func (r Resource) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Path, validation.Required),
	)
}
