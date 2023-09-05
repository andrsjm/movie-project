package util

import (
	"movie-project/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnError(c *gin.Context, err error) {
	var response model.Response
	response.Status = 400
	response.Message = err.Error()
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadRequest, response)
}

func UnauthorizedError(c *gin.Context) {
	var response model.Response
	response.Status = 400
	response.Message = "Unauthorized Access"
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusUnauthorized, response)
}
