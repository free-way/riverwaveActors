package services

import (
	"github.com/free-way/riverwaveActors/models"
	"github.com/free-way/riverwaveActors/utils"
	"fmt"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	err error
)

func GetAllUsers(ctx *gin.Context) {
	var users []models.User
	utils.Db.Preload("Role").Find(&users)
	ctx.JSON(http.StatusOK, users)
}

func CreateUser(ctx *gin.Context) {
	var user models.User
	var role models.Role
	var pswdByte []byte
	ctx.BindJSON(&user)
	if err = user.Validation(); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"Errors": err.Error(),
		})
		ctx.Abort()
		return
	}
	//encrypt password
	pswdByte, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	user.Password = string(pswdByte)
	if utils.Db.Where("id = ?", user.RoleID).Find(&role).RecordNotFound(); err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"Error": "Role Not Found",
		})
		ctx.Abort()
		return
	}
	if err = utils.Db.Preload("Role").Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func EditUser(ctx *gin.Context) {
	var err = errors.New("something went wrong")
	var user models.User
	var payload map[string]interface{}
	userId, _ := strconv.Atoi(ctx.Param("user"))
	if err = ctx.BindJSON(&payload);err != nil{
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}

	if utils.Db.Preload("Role").Where("id = ?",userId).Find(&user).RecordNotFound() {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"Error": "User Not Found",
		})
		ctx.Abort()
		return
	}

	if err = utils.Db.Model(&user).Omit("password","role_id").Update(payload).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK,user)
}

func DeleteUser(ctx *gin.Context) {
	userId,_ := strconv.Atoi(ctx.Param("user"))
	var user models.User
	if utils.Db.Where("ID = ?", userId).Find(&user).RecordNotFound() {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"Error": "User" + err.Error(),
		})
		ctx.Abort()
		return
	}
	if err = utils.Db.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK,map[string]string{
		"Message" : "User successfully Deleted",
	})
}

func Authenticate(ctx *gin.Context) {
	var payload map[string]string
	if err = ctx.BindJSON(&payload); err != nil{
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	user := models.User{}
	token := &jwt.Token{}
	var tokenString string
	if utils.Db.Where("Email = ?", payload["email"]).Find(&user).RecordNotFound() {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"Error": "User Not Find",
		})
		ctx.Abort()
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload["password"])); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"role_id":    user.RoleID,
	})
	tokenString, err = token.SignedString([]byte("riverwave"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK,map[string]string{
		"Token":tokenString,
	})

}

func ValidateToken(ctx *gin.Context) {
	var payload map[string]string
	ctx.BindJSON(&payload)
	if token,ok := payload["token"]; ok{
		if token == "" {
			err = errors.New("token not provided")
			ctx.JSON(http.StatusBadRequest, map[string]string{
				"Error": err.Error(),
			})
			ctx.Abort()
			return
		}
	}else{
		err = errors.New("token not provided")
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"Error": "Token not set in request",
		})
		ctx.Abort()
		return
	}


	token, err := jwt.Parse(payload["token"], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("riverwave"), nil
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := models.User{}
		if utils.Db.Where("Email = ? and id = ?", claims["email"], claims["id"]).Find(&user).RecordNotFound() {

			ctx.JSON(http.StatusUnauthorized, map[string]string{
				"Error": "Unauthorized",
			})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, map[string]string{
			"Error": "Authorized",
		})

	}else{
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": "Could not fetch claims from payload",
		})
		ctx.Abort()
		return
	}

}
