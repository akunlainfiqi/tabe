package controller

import (
	"net/http"
	"saas-billing/app/commands"

	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	createOrganizationCommand commands.CreateOrganizationCommand
}

func NewOrganizationController(
	createOrganizationCommand commands.CreateOrganizationCommand,
) *OrganizationController {
	return &OrganizationController{
		createOrganizationCommand: createOrganizationCommand,
	}
}

func (c *OrganizationController) Create(ctx *gin.Context) {
	var params struct {
		OrgId          string `json:"org_id" binding:"required"`
		ContactName    string `json:"contact_name" binding:"required"`
		ContactEmail   string `json:"contact_email" binding:"required"`
		ContactPhone   string `json:"contact_phone" binding:"required"`
		ContactAddress string `json:"contact_address" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	req := commands.NewCreateOrganizationRequest(
		params.OrgId,
		params.ContactName,
		params.ContactEmail,
		params.ContactPhone,
		params.ContactAddress,
	)

	if err := c.createOrganizationCommand.Execute(req); err != nil {
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
