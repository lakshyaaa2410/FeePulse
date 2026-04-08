package controller

import (
	"fee-reminder/service"

	"github.com/gin-gonic/gin"
)

type ControllerInterface interface {
	AddMembersFromCSV(c *gin.Context)
	GetAllMembers(ctx *gin.Context)
}

type Controller struct {
	service service.ServiceInterface
}

func InitializeController(service service.ServiceInterface) *Controller {
	return &Controller{service: service}
}
