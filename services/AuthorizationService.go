package services

import (
	"context"
	"github.com/free-way/riverwaveCommon/definitions"
)

type AuthorizationService struct {

}

func (AuthorizationService) CanUser(ctx context.Context,r *definitions.AuthorizationParams)(*definitions.AuthorizationResponse,error){
	return &definitions.AuthorizationResponse{},nil
}
