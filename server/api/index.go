package api

import (
	"github.com/gin-gonic/gin"
)

// GetIndex godoc
// @Summary      Index Endpoint
// @Description  Responds with a welcome message
// @Tags         Index
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the index page!",
	})
}
