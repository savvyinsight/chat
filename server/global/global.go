package global

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	GVA_DB    *gorm.DB
	GVA_REDIS *redis.Client
	GVA_CTX   = context.Background()
)
