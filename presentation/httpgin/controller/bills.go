package controller

import (
	"net/http"
	"saas-billing/app/commands"

	"github.com/gin-gonic/gin"
)

type BillController struct {
	ExpireBillsCommand commands.ExpireBillsCommand
	PayBillsCommand    commands.PayBillsCommand
	CreateBillCommand  commands.CreateBillsCommand
}

func NewBillController(
	expireBillsCommand commands.ExpireBillsCommand,
	payBillsCommand commands.PayBillsCommand,
	createBillCommand commands.CreateBillsCommand,
) *BillController {
	return &BillController{
		ExpireBillsCommand: expireBillsCommand,
		PayBillsCommand:    payBillsCommand,
		CreateBillCommand:  createBillCommand,
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
		TenantID            string `json:"tenant_id" binding:"required"`
		UseRemainingBalance bool   `json:"use_remaining_balance"`
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
		params.UseRemainingBalance,
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
