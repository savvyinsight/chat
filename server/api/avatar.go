package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"chat/global"
)

// UploadAvatar handles multipart avatar upload for current user.
// @Summary Upload avatar for current user
// @Accept multipart/form-data
// @Param avatar formData file true "avatar file"
// @Success 200 {object} map[string]interface{}
// @Router /user/avatar [post]
func UploadAvatar(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	var uid uint
	if idf, ok := claims["id"].(float64); ok {
		uid = uint(idf)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "avatar file required", "error": err.Error()})
		return
	}

	// ensure upload dir exists
	uploadDir := filepath.Join("web", "avatars")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create upload dir", "error": err.Error()})
		return
	}

	ext := filepath.Ext(file.Filename)
	fname := fmt.Sprintf("%d_%d%s", uid, time.Now().Unix(), ext)
	dst := filepath.Join(uploadDir, fname)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save avatar", "error": err.Error()})
		return
	}

	// store public path (served under /static)
	avatarURL := "/static/avatars/" + fname

	// update user record
	if err := global.GVA_DB.Model("user_basic").Where("id = ?", uid).Update("avatar_url", avatarURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user avatar", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "avatar_url": avatarURL})
}
