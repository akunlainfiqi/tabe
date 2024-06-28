package controller

import (
	"net/http"
	"saas-billing/app/commands"
	"saas-billing/app/queries"

	"github.com/gin-gonic/gin"
)

type TenantController struct {
	CreateTenantCommand commands.CreateTenantOnboardingCommand

	TenantQuery queries.TenantQuery
}

func NewTenantController(
	createTenantCommand commands.CreateTenantOnboardingCommand,

	tenantQuery queries.TenantQuery,
) *TenantController {
	return &TenantController{
		CreateTenantCommand: createTenantCommand,

		TenantQuery: tenantQuery,
	}
}

func (c *TenantController) CreateTenant(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.JSON(401,
			gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
		return
	}

	var params struct {
		PriceId    string `json:"price_id" binding:"required"`
		OrgId      string `json:"org_id" binding:"required"`
		TenantId   string `json:"tenant_id" binding:"required"`
		TenantName string `json:"tenant_name" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	req, err := commands.NewCreateTenantOnboardingRequest(
		params.PriceId,
		params.OrgId,
		params.TenantId,
		params.TenantName,
		userId,
	)
	if err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	res, err := c.CreateTenantCommand.Execute(req)
	if err != nil {
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
			"data":    res,
		})
}

func (c *TenantController) GetTenant(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.JSON(401,
			gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "id is required",
			})
		return
	}

	tenant, err := c.TenantQuery.FindByID(id)
	if err != nil {
		ctx.JSON(500,
			gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		return
	}

	ctx.JSON(200,
		gin.H{
			"status":  http.StatusOK,
			"message": "success",
			"data":    tenant,
		})
}

func (c *TenantController) GetByOrgID(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.JSON(401,
			gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
		return
	}

	orgId := ctx.Param("org_id")
	if orgId == "" {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "org_id is required",
			})
		return
	}

	tenants, err := c.TenantQuery.FindByOrgID(orgId)
	if err != nil {
		ctx.JSON(500,
			gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		return
	}

	ctx.JSON(200,
		gin.H{
			"status":  http.StatusOK,
			"message": "success",
			"data":    tenants,
		})
}
