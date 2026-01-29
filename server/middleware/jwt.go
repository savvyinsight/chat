package middleware

import (
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"chat/model"
	"chat/service"

	"gorm.io/gorm"
)

var identityKey = "id"

// JWTMiddleware returns a configured Gin-JWT middleware instance.
func JWTMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "chat zone",
		Key:         []byte(getSecret()),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.UserBasic); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			if idf, ok := claims[identityKey].(float64); ok {
				return &model.UserBasic{Model: gorm.Model{ID: uint(idf)}}
			}
			return &model.UserBasic{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals struct {
				Identifier string `json:"identifier" binding:"required"`
				Password   string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			user, err := service.AuthenticateUser(loginVals.Identifier, loginVals.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*model.UserBasic); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"message": message})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error: " + err.Error())
	}
	return authMiddleware
}

func getSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "secret"
}
