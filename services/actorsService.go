package services

import (
	"github.com/free-way/riverwaveCommon/definitions"
	"context"
	"github.com/free-way/riverwaveActors/models"
	"github.com/free-way/riverwaveActors/utils"
	"github.com/jinzhu/copier"
	"fmt"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/fatih/structs"
	"github.com/dgrijalva/jwt-go"
)

var (
	err error
)
type ActorsService struct {

}

func (ActorsService) GetAllUsers(ctx context.Context,empty *definitions.Empty) (*definitions.GetAllUsersResponse,error)  {
	users := []models.User{}
	var resUser []*definitions.User
	if err = utils.Db.Find(&users).Error; err != nil{
		fmt.Printf(err.Error())
		goto ErrorOccurred
	}
	if err = copier.Copy(&resUser,&users); err != nil{
		goto ErrorOccurred
	}

	return &definitions.GetAllUsersResponse{
		Status: 0,
		Message: "Success",
		Users:resUser,
	},nil

	ErrorOccurred:
		return  &definitions.GetAllUsersResponse{
			Status: -1,
			Message: "Error Occured",
		},err
}

func (ActorsService) CreateUser(ctx context.Context, request *definitions.CreateUserRequest) (*definitions.GeneralResponse,error)  {
	var user models.User
	var role models.Role
	var pswdByte []byte
	if err = copier.Copy(&user,request); err != nil{
		goto ErrorOccurred
	}
	if err = user.Validation(); err != nil{
		goto ErrorOccurred
	}

	//encrypt password
	pswdByte,err = bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.MinCost)
	user.Password = string(pswdByte)
	if request.RoleID > 0{
		 if err = utils.Db.Where("ID = ?",request.RoleID).Find(&role).Error; err != nil{
			goto ErrorOccurred
		}
	}
	if err = utils.Db.Create(&user).Error; err != nil{
		goto ErrorOccurred
	}
	return &definitions.GeneralResponse{
		Status:0,
		Message:"Success",
	},nil


	ErrorOccurred:
		return  &definitions.GeneralResponse{
			Status:-1,
			Message: "Error Occurred",
		},err
}

func (ActorsService) EditUser(ctx context.Context, request *definitions.EditUserRequest) (*definitions.GeneralResponse,error){
	var user models.User
	var payload map[string]interface{}
	var role models.Role
	if err = copier.Copy(&user,request); err != nil{
		goto ErrorOccurred
	}
	if err = user.Validation(); err != nil{
		goto ErrorOccurred
	}
	if request.RoleID != 0{
		if err = utils.Db.Where("ID = ?",request.RoleID).Find(&role).Error; err != nil{
			goto ErrorOccurred
		}
	}
	if request.ID < 1{
		err = errors.New("ID can not be null")
		goto ErrorOccurred
	}else{
		payload = structs.Map(&user)
		if err = utils.Db.Find(&user).Error; err != nil{
			goto ErrorOccurred
		}else{

			if err = utils.Db.Model(&user).Omit("password").Update(payload).Error; err != nil{
				goto ErrorOccurred
			}
		}
	}
	return &definitions.GeneralResponse{
		Status:0,
		Message:"Success",
	},nil

ErrorOccurred:
	return  &definitions.GeneralResponse{
		Status:-1,
		Message: "Error Occurred",
	},err
}

func (ActorsService) DeleteUser(ctx context.Context, request *definitions.DeleteUserRequest) (*definitions.GeneralResponse,error){
	var user models.User
	if err = utils.Db.Where("ID = ?",request.UserID).Find(&user).Error; err != nil {
		goto ErrorOccurred
	}
	if err = utils.Db.Delete(&user).Where("ID = ?",request.UserID).Error; err != nil{
		goto ErrorOccurred
	}
	return &definitions.GeneralResponse{
		Status:0,
		Message:"Success",
	},nil
ErrorOccurred:
	return  &definitions.GeneralResponse{
		Status:-1,
		Message: "Error Occurred",
	},err
}

func (ActorsService) Authenticate(ctx context.Context, request *definitions.AuthenticationParams) (*definitions.GeneralResponse,error){
	user := models.User{}
	token := &jwt.Token{}
	var tokenString string
	if utils.Db.Where("Email = ?",request.Email).Find(&user).RecordNotFound(){
		errors.New("email not found")
		goto ErrorOccurred
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(request.Password)); err != nil{
		errors.New("password does not match")
		goto ErrorOccurred
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"id":user.ID,
		"first_name":user.FirstName,
		"last_name": user.LastName,
		"email":user.Email,
		"role_id": user.RoleID,

	})
	tokenString,err = token.SignedString([]byte("riverwave"))
	if err != nil{
		goto ErrorOccurred
	}
	return &definitions.GeneralResponse{
		Status:0,
		Message:tokenString,
	},nil

ErrorOccurred:
	return  &definitions.GeneralResponse{
		Status:-1,
		Message: "Error Occurred",
	},err
}
