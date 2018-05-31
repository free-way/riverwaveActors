package models

import (
	"github.com/jinzhu/gorm"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)


type User struct {
	gorm.Model
	//ID        uint `gorm:"primary_key"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique_index"`
	Password  string
	Role   Role `gorm:"foreignkey:RoleID"`
	RoleID uint
}

func (User) TableName() string {
	return "rw_users"
}
func (user User) Validation() error {
	return validation.ValidateStruct(&user,
		validation.Field(&user.FirstName, validation.Required, validation.Length(3, 50)),
		validation.Field(&user.LastName, validation.Required,validation.Length(3,50)),
		validation.Field(&user.Email, validation.Required,is.Email),
		validation.Field(&user.Password,validation.Required,),
	)
}
