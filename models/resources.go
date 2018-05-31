package models

import "github.com/jinzhu/gorm"

type Resource struct {
	gorm.Model
	//ID uint `gorm:"primary_key"`
	Name string
}

type ResourcePermission struct {
	gorm.Model
	//ID uint `gorm:"primary_key"`
	Resource Resource `gorm:"foreignKey:ResourceID"`
	ResourceID uint
	Permission Permission `gorm:"foreignKey:PermissionID"`
	PermissionID uint
}


func (Resource)TableName() string  {
	return "rw_resources"
}

func (ResourcePermission) TableName()string{
	return "rw_resources_permissions"
}
