package redis

import (
	"base/config"
	"base/db"
	"fmt"
	"os"
	"webb-auth/common"

	"github.com/sirupsen/logrus"
)

func Init(cfg *config.RedisConfig) {
	fmt.Println("init redis")
	if redisclient, err := db.RedisInit(cfg); err != nil {
		logrus.Errorf("failed to init redis, error: %s", err)
		os.Exit(1)
	} else {
		common.RC = redisclient
	}
}
