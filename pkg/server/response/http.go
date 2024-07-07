package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"Data":    data,
	})
}

func BadRequest(c *gin.Context, err error, data any) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Success": false,
		"Message": err.Error(),
		"Data":    data,
	})
}

func NotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, gin.H{
		"Success": false,
		"Message": err.Error(),
	})
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"Success": false,
		"Message": err.Error(),
	})
}
