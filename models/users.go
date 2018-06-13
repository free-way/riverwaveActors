package models

import (
	"github.com/jinzhu/gorm"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)


type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique_index";json:"email"`
	Password  string `json:"password"`
	Role   Role `gorm:"foreignkey:RoleID"`
	RoleID uint `json:"role_id"`
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
		validation.Field(&user.RoleID,validation.Required),
	)
}
