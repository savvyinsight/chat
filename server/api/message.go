package api

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "chat/service"
)

// GetMessages godoc
// @Summary Get chat messages between two users
// @Param user_id query int true "Current user id"
// @Param with query int true "Other user id"
// @Param limit query int false "Limit"
// @Success 200 {object} map[string]interface{}
// @Router /messages [get]
func GetMessages(c *gin.Context) {
    userStr := c.Query("user_id")
    withStr := c.Query("with")
    if userStr == "" || withStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "user_id and with are required"})
        return
    }
    uid, err := strconv.Atoi(userStr)
    if err != nil || uid <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user_id"})
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
