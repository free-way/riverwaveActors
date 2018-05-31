package models

import "github.com/jinzhu/gorm"

type Permission struct {
	gorm.Model
	//ID uint `gorm:"primary_key"`
	Name string `gorm:"unique_index"`
}

type RolePermission struct {
	gorm.Model
	//ID uint `gorm:"primary_key"`
	Role   Role `gorm:"foreignkey:RoleID"`
	RoleID uint
	Permission Permission `gorm:"foreignkey:PermissionID"`
	PermissionID uint
}


func (Permission) TableName() string{
	return "rw_permissions"
}

func (RolePermission) TableName()string{
	return "rw_roles_permissions"
}
