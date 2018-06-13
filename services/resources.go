package services

import (
	"github.com/free-way/riverwaveActors/models"
	"github.com/free-way/riverwaveActors/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)


func GetResources(ctx *gin.Context){
	var resources []*models.Resource
	utils.Db.Find(&resources)
	ctx.JSON(http.StatusOK,resources)

}

func AddResource(ctx *gin.Context) {
	var resource models.Resource
	ctx.BindJSON(&resource)

	if err = resource.Validate(); err != nil{
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"Errors": err.Error(),
		})
		ctx.Abort()
		return
	}
	if err = utils.Db.Create(&resource).Error; err != nil{
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK,resource)
}

func  EditResource(ctx *gin.Context) {
	resourceId,_ := strconv.Atoi(ctx.Param("resource"))
	var payload map[string]interface{}
	var resource models.Resource
	ctx.BindJSON(&payload)
	if utils.Db.Where("id = ?",resourceId).Find(&resource).RecordNotFound(){
		ctx.JSON(http.StatusNotFound, map[string]string{
			"Error": "Resource Not Found",
		})
		ctx.Abort()
		return
	}
	utils.Db.Model(&resource).Update(payload)
	ctx.JSON(http.StatusOK,resource)
}

func  DeleteResource(ctx *gin.Context){
	resourceId,_ := strconv.Atoi(ctx.Param("resource"))
	var resource models.Resource
	if utils.Db.Where("id = ?",resourceId).Find(&resource).RecordNotFound(){
		ctx.JSON(http.StatusNotFound, map[string]string{
			"Error": "Resource Not Found",
		})
		ctx.Abort()
		return
	}
	utils.Db.Delete(&resource)
	ctx.JSON(http.StatusOK,map[string]string{
		"Message":"Resource Deleted",
	})
}
