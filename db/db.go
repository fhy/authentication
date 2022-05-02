package db

import (
	"base/config"
	"base/db"
	"fmt"
	"os"
	"webb-auth/common"
	"webb-auth/user"
	"webb-auth/wechat"

	logger "github.com/sirupsen/logrus"
)

func Init(config *config.DbConfig) {
	fmt.Println("init db")
	if db, err := db.DbInit(config); err != nil {
		logger.Errorf("failed to init db, error: %s", err)
		os.Exit(1)
	} else {
		db.AutoMigrate(
			&user.User{},
			&wechat.WeChat{})
		common.DB = db
	}
}
