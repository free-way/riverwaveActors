package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"fmt"
	"os"
	"net"
	"google.golang.org/grpc"
	"github.com/free-way/riverwaveCommon/definitions"
	"github.com/free-way/riverwaveActors/services"
	"log"
	"flag"
	"github.com/free-way/riverwaveActors/utils"
	"github.com/jinzhu/gorm"
	"github.com/free-way/riverwaveActors/models"
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

	//create grpc server
	listener, err := net.Listen("tcp", cfg.Section("Microservice").Key("port").String())
	if err != nil {
		fmt.Println("Can not listen on port")
		os.Exit(-1)
	}
	service := grpc.NewServer()
	definitions.RegisterActorsServiceServer(service, services.ActorsService{})
	definitions.RegisterAuthorizationServiceServer(service, services.AuthorizationService{})

	if err := service.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
