package controller

import (
	"net/http"
	"saas-billing/app/commands"

	"github.com/gin-gonic/gin"
)

type TenantController struct {
	CreateTenantCommand commands.CreateTenantFromProductCommand
}

func NewTenantController(
	createTenantCommand commands.CreateTenantFromProductCommand,
) *TenantController {
	return &TenantController{
		CreateTenantCommand: createTenantCommand,
	}
}

func (c *TenantController) CreateTenant(ctx *gin.Context) {
	var params struct {
		PriceId    string `json:"price_id" binding:"required"`
		OrgId      string `json:"org_id" binding:"required"`
		TenantId   string `json:"tenant_id" binding:"required"`
		TenantName string `json:"tenant_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	req, err := commands.NewCreateTenantFromProductRequest(
		params.PriceId,
		params.OrgId,
		params.TenantId,
		params.TenantName,
	)
	if err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	if err := c.CreateTenantCommand.Execute(req); err != nil {
		ctx.JSON(500,
			gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		return
	}

	ctx.JSON(201,
		gin.H{
			"status":  http.StatusCreated,
			"message": "success",
		})
}
