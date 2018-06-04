package services

import (
	"context"
	"github.com/free-way/riverwaveCommon/definitions"
	"github.com/free-way/riverwaveActors/models"
	"github.com/free-way/riverwaveActors/utils"
	"github.com/jinzhu/copier"
	"errors"
)

type Resources struct {

}

func (Resources) GetResources(ctx context.Context,empty *definitions.Empty) (*definitions.ResourcesResponse,error){
	var resources []*models.Resource
	var err error
	var res definitions.ResourcesResponse
	err = utils.Db.Find(&resources).Error
	if err != nil{
		return &definitions.ResourcesResponse{
			Message:err.Error(),
			Status: -1,
		},err
	}
	copier.Copy(&res.Resources,&resources)
	return &res,nil

}

func (Resources) AddResource(ctx context.Context, req *definitions.CreateResourceParams)(*definitions.GeneralResponse,error)  {
	var resource models.Resource
	var err error
	resource.Name = req.Name
	resource.Path = req.Path
	if err = resource.Validate(); err != nil{
		return &definitions.GeneralResponse{
			Message:err.Error(),
			Status: -1,
		},err
	}
	if err = utils.Db.Create(&resource).Error; err != nil{
		return &definitions.GeneralResponse{
			Message:err.Error(),
			Status: -1,
		},err
	}
	return &definitions.GeneralResponse{
		Message:"Success",
		Status: 0,
	},nil
}

func (Resources) EditResource(ctx context.Context, req *definitions.EditResourceParams)(*definitions.GeneralResponse,error)  {
	var resource models.Resource
	var err error
	if utils.Db.Where("id = ?",req.ResourceId).Find(&resource).RecordNotFound(){
		err = errors.New("record not found")
		return &definitions.GeneralResponse{
			Message:err.Error(),
			Status: -1,
		},nil
	}
	resource.Path = req.Path
	resource.Name = req.Name
	if err = resource.Validate(); err != nil{
		return &definitions.GeneralResponse{
			Message: err.Error(),
			Status: -1,
		},err
	}
	utils.Db.Save(&resource)
	return &definitions.GeneralResponse{
		Message:"Success",
		Status: 0,
	},nil
}

func (Resources) DeleteResource(ctx context.Context, req *definitions.DeleteResourceParams)(*definitions.GeneralResponse,error){
	var resource models.Resource
	var err error
	if utils.Db.Where("id = ?",req.ResourceId).Find(&resource).RecordNotFound(){
		err = errors.New("record not found")
		return &definitions.GeneralResponse{
			Message:err.Error(),
			Status: -1,
		},nil
	}
	utils.Db.Delete(&resource)
	return &definitions.GeneralResponse{
		Message:"Success",
		Status: -1,
	},nil
}
