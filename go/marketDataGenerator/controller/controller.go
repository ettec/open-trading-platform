package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ShowAccount godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Router /accounts/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {
	id := ctx.Param("id")

	b := Boom{id, "bob", 3}


	ctx.JSON(http.StatusOK, b)
}

type Boom struct {
	id string
	name string
	age int
}

type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}
