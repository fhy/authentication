package router

import (
	"webb-auth/conf"
	"webb-auth/controllers"

	"github.com/fhy/utils-golang/utils"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func initRouterV1(r *gin.Engine) {
	verifyKey := utils.LoadEdPublicKeyFromDisk(conf.Conf.Jwt.PubKeyPath)
	v1 := r.Group("/v1")
	redirect := v1.Use(requestid.New(), utils.TrailClient(conf.Conf.Cookie.Domain, conf.Conf.Cookie.AgeMax))
	{
		redirect.POST("login", controllers.Login)
		redirect.GET("token", controllers.GetToken)
		redirect.GET("dingtalk/auth", controllers.DingTalkAuth) // 钉钉认证
		redirect.GET("wechat/offiaccount/redirecturl", controllers.GetOfficialRedirectURL)
		redirect.GET("wechat/offiaccount/auth", controllers.OfficialAccountAuth)   // 微信公众号认证
		redirect.GET("wechat/miniprogram/auth/:code", controllers.MiniProgramAuth) // 微信小程序登录凭证校验
		redirect.GET("wechat/qrcode", controllers.GetMiniProgromQrcode)            // 微信小程序生成二维码
		redirect.POST("user/registry", controllers.Registry)
	}
	i := v1.Use(requestid.New(), utils.TrailClient(conf.Conf.Cookie.Domain, conf.Conf.Cookie.AgeMax),
		utils.UserAuthen(verifyKey))
	{
		i.GET("logout", controllers.Logout)
		i.POST("token", controllers.RefleshToken)
		i.GET("user/info", controllers.GetUserInfo)
	}
}
