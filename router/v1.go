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
	v1 := r.Group("/api/v1")
	everyBody := v1.Use(requestid.New(), utils.TrailClient(conf.Conf.Cookie.Domain, conf.Conf.Cookie.AgeMax))
	{
		everyBody.POST("login", controllers.Login)
		everyBody.GET("wechat/offiaccount/redirecturl", controllers.GetOfficialRedirectURL)
		everyBody.GET("wechat/offiaccount/auth", controllers.OfficialAccountAuth)    // 微信公众号认证
		everyBody.GET("wechat/miniprogram/login/:code", controllers.MiniProgramAuth) // 微信小程序登录凭证校验
		everyBody.GET("user/info", controllers.GetUserInfo)
	}
	user := v1.Use(requestid.New(), utils.TrailClient(conf.Conf.Cookie.Domain, conf.Conf.Cookie.AgeMax),
		utils.GetUser(&verifyKey))
	{
		user.GET("logout", controllers.Logout)
		user.GET("token", controllers.GetToken)
		user.POST("token", controllers.RefleshToken)
		user.GET("wechat/qrcode", controllers.GetMiniProgromQrcode) // 微信小程序生成二维码
		user.POST("user/registry", controllers.Registry)
	}

}
