package router

import (
	"base/utils"
	"webb-auth/conf"
	"webb-auth/controllers"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func initRouterV1(r *gin.Engine) {
	verifyKey := utils.LoadEdPublicKeyFromDisk(conf.Conf.Jwt.PubKeyPath)
	v1 := r.Group("/v1")
	i := v1.Use(requestid.New(), utils.TrailClient(conf.Conf.Cookie.Domain, conf.Conf.Cookie.AgeMax),
		utils.GetUser(&verifyKey))
	{
		i.POST("login", controllers.Login)
		i.GET("logout", controllers.Logout)
		i.GET("token", controllers.GetToken)
		i.GET("wechat/offiaccount/redirecturl", controllers.GetOfficialRedirectURL)
		i.GET("wechat/offiaccount/auth", controllers.OfficialAccountAuth)   // 微信公众号认证
		i.GET("wechat/miniprogram/auth/:code", controllers.MiniProgramAuth) // 微信小程序登录凭证校验
		i.GET("wechat/qrcode", controllers.GetMiniProgromQrcode)            // 微信小程序生成二维码
		i.POST("user/registry", controllers.Registry)
		i.GET("user/info", controllers.GetUserInfo)
	}
}
