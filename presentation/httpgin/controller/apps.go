package controller

import (
	"net/http"
	"saas-billing/app/queries"

	"github.com/gin-gonic/gin"
)

type AppController struct {
	appQuery queries.AppQueries
}

func NewAppController(appQuery queries.AppQueries) *AppController {
	return &AppController{appQuery}
}

func (c *AppController) GetAll(ctx *gin.Context) {
	apps, err := c.appQuery.GetAll()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": http.StatusOK,
		"data":   apps,
	})
}
