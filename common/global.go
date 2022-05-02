package common

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	RC *redis.Client
)
