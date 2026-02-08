package api

import (
    "net/http"
    "strconv"

    jwt "github.com/appleboy/gin-jwt/v2"
    "github.com/gin-gonic/gin"

    "chat/service"
)

// GetMessages godoc
// @Summary Get chat messages between two users
// @Param with query int true "Other user id"
// @Param limit query int false "Limit"
// @Success 200 {object} map[string]interface{}
// @Router /messages [get]
func GetMessages(c *gin.Context) {
    // derive current user id from JWT claims
    claims := jwt.ExtractClaims(c)
    var uid int
    if idf, ok := claims["id"].(float64); ok {
        uid = int(idf)
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
        return
    }
    withStr := c.Query("with")
    if withStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "with is required"})
        return
    }
    wid, err := strconv.Atoi(withStr)
    if err != nil || wid <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid with"})
        return
    }
    limit := 100
    if lstr := c.Query("limit"); lstr != "" {
        if lv, err := strconv.Atoi(lstr); err == nil && lv > 0 {
            limit = lv
        }
    }
    msgs, err := service.GetMessagesBetween(uint(uid), uint(wid), limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load messages", "error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "ok", "data": msgs})
}
