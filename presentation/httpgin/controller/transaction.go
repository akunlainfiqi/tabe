package controller

import (
	"net/http"
	"saas-billing/app/queries"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionQuery queries.TransactionQuery
}

func NewTransactionController(
	transactionQuery queries.TransactionQuery,
) *TransactionController {
	return &TransactionController{
		TransactionQuery: transactionQuery,
	}
}

func (c *TransactionController) GetByOrgID(ctx *gin.Context) {
	orgID := ctx.Param("org_id")

	transactions, err := c.TransactionQuery.FindByOrgID(orgID)
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
			"data":    transactions,
		})
}

func (c *TransactionController) GetByBillsID(ctx *gin.Context) {
	billID := ctx.Param("bill_id")

	transaction, err := c.TransactionQuery.FindByBillsID(billID)
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
			"data":    transaction,
		})
}
