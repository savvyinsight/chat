package api

import (
	"chat/model"
	"chat/service"
	"net/http"
	"strconv"

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
	data := service.GetUserList()
	c.JSON(200, gin.H{
		"message": "User list retrieved successfully",
		"data":    data,
	})
}

// Create User
// @Summary      Create User
// @Description  Create user
// @Tags         User
// @Param name query string false "User Name"
// @Param password query string false "Password"
// @Param repassword query string false "Confirm Password"
// @Success      200  {string}  json{"code","message"}
// @Router       /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := model.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	if password != repassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password not consistent!",
			"error":   "PASSWORD_MISMATCH",
		})
		return
	}

	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password cannot be empty!",
			"error":   "EMPTY_PASSWORD",
		})
		return
	}

	user.PassWord = password

	err := service.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Create User Failed!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Create User Succeeded!",
		"user_id": user.ID,
	})
}

// Delete User
// @Summary      Delete User
// @Tags         User
// @Param id path int true "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /user/{id} [delete]
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"error":   "INVALID_ID",
		})
		return
	}

	user := model.UserBasic{}
	user.ID = uint(id)

	err = service.DeleteUser(&user)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
				"error":   "USER_NOT_FOUND",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Delete User Failed!",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete User Succeeded!",
		"user_id": id,
	})
}
