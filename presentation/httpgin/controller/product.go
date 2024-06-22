package controller

import (
	"net/http"
	"saas-billing/app/commands"
	"saas-billing/app/queries"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductQuery queries.ProductQuery

	createProductCommand commands.CreateProductCommand
}

func NewProductController(
	productQueries queries.ProductQuery,
	createProductCommand commands.CreateProductCommand,
) *ProductController {
	return &ProductController{
		ProductQuery:         productQueries,
		createProductCommand: createProductCommand,
	}
}

func (pc *ProductController) GetAll(ctx *gin.Context) {
	products, err := pc.ProductQuery.FindAll()
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
			"data":    products,
		})
}

func (pc *ProductController) GetByAppID(ctx *gin.Context) {
	appID := ctx.Param("app_id")

	products, err := pc.ProductQuery.FindByAppID(appID)
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
			"data":    products,
		})
}

func (c *ProductController) GetByID(ctx *gin.Context) {
	var params struct {
		ID string `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	product, err := c.ProductQuery.FindByID(params.ID)
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
			"data":    product,
		})
}

func (c *ProductController) Create(ctx *gin.Context) {
	var params struct {
		Name  string                              `json:"name" binding:"required"`
		Tiers []commands.CreateProductTierRequest `json:"tiers" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
		return
	}

	req := commands.NewCreateProductRequest(
		params.Name,
		params.Tiers,
	)

	if err := c.createProductCommand.Execute(req); err != nil {
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
