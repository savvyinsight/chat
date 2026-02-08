package api

import (
	"chat/global"
	"chat/middleware"
	"chat/model"
	"chat/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
	"gorm.io/gorm"
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

// GetCurrentUser returns the profile of the currently authenticated user
func GetCurrentUser(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	var uid uint
	if idf, ok := claims["id"].(float64); ok {
		uid = uint(idf)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	var user model.UserBasic
	if err := global.GVA_DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load user", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "data": user})
}

// Create User
// @Summary      Create User
// @Description  Create user
// @Tags         User
// @Param name query string false "User Name"
// @Param email query string false "User Email"
// @Param phone query string false "User Phone"
// @Param password query string false "Password"
// @Param repassword query string false "Confirm Password"
// @Success      200  {string}  json{"code","message"}
// @Router       /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := model.UserBasic{}
	user.Name = c.Query("name")
	user.Email = c.Query("email")
	user.Phone = c.Query("phone")
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

	if user.Email == "" && user.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or phone required",
			"error":   "MISSING_IDENTIFIER",
		})
		return
	}

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

// Register user (JSON)
// @Summary Register
// @Accept json
// @Produce json
// @Param request body map[string]string true "Register request {name, email, phone, password, repassword}"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /user/register [post]
func Register(c *gin.Context) {
	var req struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Password   string `json:"password"`
		Repassword string `json:"repassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	if req.Password != req.Repassword {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password and repassword do not match"})
		return
	}

	user := model.UserBasic{}
	user.Name = req.Name
	user.Email = req.Email
	user.Phone = req.Phone
	user.PassWord = req.Password

	if err := service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Create User Failed!", "error": err.Error()})
		return
	}

	// generate token for newly registered user
	token, terr := middleware.GenerateTokenForUser(&user)
	if terr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token", "error": terr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Register succeeded", "user_id": user.ID, "token": token})
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

// Update User
// @Summary      Update User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param id path int true "User ID"
// @Param request body model.UserBasic true "Update User Request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /user/{id} [put]
func UpdateUser(c *gin.Context) {
	// 1. 获取用户ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"error":   "INVALID_ID",
		})
		return
	}

	// 2. 先检查用户是否存在
	var existingUser model.UserBasic
	result := global.GVA_DB.First(&existingUser, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
				"error":   "USER_NOT_FOUND",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to query user",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// 3. 绑定更新数据
	var updateData model.UserBasic
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// 4. 设置要更新的字段（避免更新ID）
	updateData.ID = uint(id)

	// 5. Validate fields (email/phone)
	if err := updateData.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	// 6. 调用Service层更新
	err = service.UpdateUser(&updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Update User Failed!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update User Succeeded!",
		"user_id": id,
	})
}

// PartialUpdateUser 部分更新用户信息
// @Summary      Partial Update User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param id path int true "User ID"
// @Param request body map[string]interface{} true "Update Fields"
// @Success      200  {object}  map[string]interface{}
// @Router       /user/{id} [patch]
func PartialUpdateUser(c *gin.Context) {
	// 1. 获取用户ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
		})
		return
	}

	// 2. 解析更新字段
	var updateFields map[string]interface{}
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON data",
			"error":   err.Error(),
		})
		return
	}

	// 3. 移除不允许更新的字段（如ID）
	delete(updateFields, "id")
	delete(updateFields, "ID")
	delete(updateFields, "created_at")
	delete(updateFields, "CreatedAt")

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No fields to update",
		})
		return
	}

	// 4. 调用Service层更新
	err = service.UpdateUserPartial(uint(id), updateFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Update User Failed!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update User Succeeded!",
		"user_id": id,
	})
}

// Login user
// @Summary Login
// @Accept json
// @Produce json
// @Param request body map[string]string true "Login request {identifier, password}"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /user/login [post]
func Login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	user, err := service.AuthenticateUser(req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed", "error": err.Error()})
		return
	}

	// generate JWT token for the authenticated user
	token, terr := middleware.GenerateTokenForUser(user)
	if terr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token", "error": terr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login succeeded", "user_id": user.ID, "token": token})
}
