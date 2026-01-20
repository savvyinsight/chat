package service

import (
	"chat/global"
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
// @Router       /userList [get]
func GetUserList(c *gin.Context) {
	db := global.GVA_DB.Model(model.UserBasic{})

	var users []model.UserBasic
	db.Find(&users)
	// data := model.GetUserList()
	data := users
	c.JSON(200, gin.H{
		"message": "User list retrieved successfully",
		"data":    data,
	})
}
