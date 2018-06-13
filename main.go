package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"fmt"
	"os"
	"flag"
	"github.com/free-way/riverwaveActors/utils"
	"github.com/jinzhu/gorm"
	"github.com/free-way/riverwaveActors/models"
	"github.com/gin-gonic/gin"
	"github.com/free-way/riverwaveActors/services"
)


var (
	err error
	cfg *ini.File
)

func init(){
	//load environment variables
	cfgFlag := flag.String("config", "", "Env file")
	flag.Parse()
	cfg, err = ini.Load(*cfgFlag)
	if err != nil {
		fmt.Println("could not load configuration file due to: ", err.Error())
		//os.Exit(-1)
	}
	connectionString := cfg.Section("Database").Key("conn").String()
	enableLogger,err := cfg.Section("Logger").Key("enable_logs").Bool()
	if err != nil{
		fmt.Printf(err.Error())
	}
	//connect to the database
	utils.Db, err = gorm.Open("mysql", connectionString)
	utils.Db.LogMode(enableLogger)
	if err != nil {
		fmt.Println("Error While trying to connect to the database: ", err.Error())
		os.Exit(1)
	}
	models.RunMigrations()
}

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users",services.GetAllUsers)
		v1.POST("/users",services.CreateUser)
		v1.PUT("/users/:user",services.EditUser)
		v1.DELETE("/users/:user",services.DeleteUser)
		v1.POST("/authenticate",services.Authenticate)
		v1.POST("/validate-token",services.ValidateToken)

		v1.GET("/resources",services.GetResources)
		v1.POST("/resources",services.AddResource)
		v1.PUT("/resources/:resource",services.EditResource)
		v1.DELETE("/resources/:resource",services.DeleteResource)
	}

	r.Run(cfg.Section("Microservice").Key("port").String())


}
