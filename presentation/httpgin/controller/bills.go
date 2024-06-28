package controller

import (
	"net/http"
	"saas-billing/app/commands"
	"saas-billing/app/queries"
	"saas-billing/errors"

	"github.com/gin-gonic/gin"
)

type BillController struct {
	ExpireBillsCommand  commands.ExpireBillsCommand
	PayBillsCommand     commands.PayBillsCommand
	CreateBillCommand   commands.CreateBillsCommand
	CheckPaymentCommand commands.CheckPaymentCommand

	IamUserOrganizationQuery queries.IamUserOrganizationQuery
	BillQuery                queries.BillQuery
}

func NewBillController(
	expireBillsCommand commands.ExpireBillsCommand,
	payBillsCommand commands.PayBillsCommand,
	createBillCommand commands.CreateBillsCommand,
	checkPaymentCommand commands.CheckPaymentCommand,

	iamUserOrganizationQuery queries.IamUserOrganizationQuery,
	billQuery queries.BillQuery,
) *BillController {
	return &BillController{
		ExpireBillsCommand:  expireBillsCommand,
		PayBillsCommand:     payBillsCommand,
		CreateBillCommand:   createBillCommand,
		CheckPaymentCommand: checkPaymentCommand,

		IamUserOrganizationQuery: iamUserOrganizationQuery,
		BillQuery:                billQuery,
	}
}

func (c *BillController) InternalExpire(ctx *gin.Context) {
	var params struct {
		BillId string `json:"bill_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	req := commands.NewExpireBillsRequest(
		params.BillId,
	)

	if err := c.ExpireBillsCommand.Execute(req); err != nil {
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
		})
}

func (c *BillController) InternalPay(ctx *gin.Context) {
	var params struct {
		BillId               string `json:"bill_id" binding:"required"`
		TransactionType      string `json:"transaction_type" binding:"required"`
		TransactionTimestamp int64  `json:"transaction_timestamp" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	req := commands.NewPayBillsRequest(
		params.BillId,
		params.TransactionType,
		params.TransactionTimestamp,
	)

	if err := c.PayBillsCommand.Execute(req); err != nil {
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
		})
}

func (c *BillController) Create(ctx *gin.Context) {
	var params struct {
		TenantID string `json:"tenant_id" binding:"required"`
		BillType string `json:"bill_type" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	req := commands.NewCreateBillsRequest(
		params.TenantID,
		params.BillType,
	)

	if err := c.CreateBillCommand.Execute(req); err != nil {
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

func (c *BillController) GetBillDetail(ctx *gin.Context) {
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
		BillID string `form:"bill_id" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	bill, err := c.BillQuery.GetByID(params.BillID)
	if err != nil {
		if err == errors.ErrBillsNotFound {
			ctx.JSON(404,
				gin.H{
					"status":  http.StatusNotFound,
					"message": err.Error(),
				})
			return
		}
		ctx.JSON(500,
			gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		return
	}

	isOwner := c.IamUserOrganizationQuery.IsOwner(bill.OrganizationId, userId)
	if !isOwner {
		ctx.JSON(403,
			gin.H{
				"status":  http.StatusForbidden,
				"message": "Forbidden",
			})
		return
	}

	ctx.JSON(200,
		gin.H{
			"status": http.StatusOK,
			"data":   bill,
		})
}

func (c *BillController) GetOrganizationBills(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.JSON(401,
			gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
		return
	}

	organizationID := ctx.Param("org_id")
	if organizationID == "" {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Organization ID is required",
			})
		return
	}

	isOwner := c.IamUserOrganizationQuery.IsOwner(organizationID, userId)
	if !isOwner {
		ctx.JSON(403,
			gin.H{
				"status":  http.StatusForbidden,
				"message": "Forbidden",
			})
		return
	}

	bills, err := c.BillQuery.GetByOrganizationID(organizationID)
	if err != nil {
		if err == errors.ErrBillsNotFound {
			ctx.JSON(404,
				gin.H{
					"status":  http.StatusNotFound,
					"message": err.Error(),
				})
			return
		}
		ctx.JSON(500,
			gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		return
	}

	ctx.JSON(200,
		gin.H{
			"status": http.StatusOK,
			"data":   bills,
		})
}

func (c *BillController) InternalCheckPayment(ctx *gin.Context) {
	if err := c.CheckPaymentCommand.Execute(); err != nil {
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
		})
}

func (c *BillController) PaymentCallback(ctx *gin.Context) {
	var params struct {
		OrderID string `json:"order_id" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	if err := c.CheckPaymentCommand.Execute(); err != nil {
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
		})
}
