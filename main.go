package main

import (
	"base/utils"
	"base/wggo"
	"webb-auth/conf"
	"webb-auth/redis"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	logrus.Info("init auth's config")
	utils.InitLogger(&conf.Conf.Log)
	redis.Init(&conf.Conf.Redis)
}

func main() {
	wggo.Wg.Wait()
}
