package service

import (
	"chat/model"

	"github.com/gin-gonic/gin"
)

// GetUserList godoc
// @Summary      Get User List
// @Description  Retrieves a list of users
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*model.UserBasic, 10)
	data = model.GetUserList()
	c.JSON(200, gin.H{
		"message": "User list retrieved successfully",
		"data":    data,
	})
}
