package models

import (
	"github.com/free-way/riverwaveActors/utils"
)
func RunMigrations()  {
	if utils.Db != nil{
		utils.Db.AutoMigrate(
			&User{},
			&Role{},
			&Permission{},
			&RolePermission{},
			&Resource{},
			&ResourcePermission{},
		)
	}

}
