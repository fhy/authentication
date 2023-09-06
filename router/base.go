package router

import (
	"github.com/fhy/utils-golang/utils"

	"github.com/gin-gonic/gin"
)

func GetGinEngine() *gin.Engine {
	r := gin.Default()
	r.Use(utils.OptionsHeader(), utils.ErrorHandler())
	initRouterV1(r)
	return r
}
