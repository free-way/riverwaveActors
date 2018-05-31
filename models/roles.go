package models

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	//ID   uint `gorm:"primary_key"`
	Name string
	Slug string `gorm:"unique_index"`
}

func (Role) TableName() string  {
	return "rw_roles"
}
