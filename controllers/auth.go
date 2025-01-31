package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"webb-auth/common"
	"webb-auth/models"
	"webb-auth/wechat"

	"github.com/fhy/utils-golang/utils"
	"github.com/fhy/utils-golang/wggo"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Login(c *gin.Context) {
	utils.ResponseSuccessJson(c, "ok")
}

func Logout(c *gin.Context) {
	utils.ResponseSuccessJson(c, "ok")
}

func GetToken(c *gin.Context) {
	if client, err := utils.GetClientInfo(c); err != nil {
		logrus.Errorf("error getting token for %s, error: %s", client.SessionId, err)
		utils.ResponseFailedJson(c, utils.ERRCODE_REQUEST_PARAM_ERROR, utils.ERRMSG_REQUEST_PARAM_ERROR, nil, http.StatusBadRequest)
	} else {
		if token, err := models.GetTokenWithSession(client); err != nil {
			logrus.Errorf("error getting token for %s, error: %s", client.LogFormatLong(), err)
			if errors.Is(errors.New(utils.ERRMSG_INVALID_SESSION), err) {
				utils.ResponseFailedJson(c, utils.ERRCODE_INVALID_SESSION, utils.ERRMSG_INVALID_SESSION, nil, http.StatusBadRequest)
			} else {
				utils.ResponseFailedJson(c, utils.ERRCODE_SERVER_ERROR, utils.ERRMSG_SERVER_ERROR, nil, http.StatusBadGateway)
			}
		} else {
			utils.ResponseSuccessJson(c, token)
		}
	}
}

func RefleshToken(c *gin.Context) {
	token := c.Param("refresh_token")
	if client, err := utils.GetClientInfo(c); err != nil {
		logrus.Errorf("client %s error getting, error: %s", client.SessionId, err)
		utils.ResponseFailedJson(c, utils.ERRCODE_REQUEST_PARAM_ERROR, utils.ERRMSG_REQUEST_PARAM_ERROR, nil, http.StatusBadRequest)
	} else {
		if token, err := models.RefleshToken(token, client); err != nil {
			logrus.Errorf("error refleshing token for %s, error: %s", client.LogFormatLong(), err)
			utils.ResponseFailedJson(c, utils.ERRCODE_INVALID_TOKEN, utils.ERRMSG_INVALID_TOKEN, nil, http.StatusBadGateway)
		} else {
			utils.ResponseSuccessJson(c, token)
			result := common.RC.Expire(context.Background(), fmt.Sprintf("%s%s", utils.USER_SID_REDIS_PREFIX, client.SessionId), utils.TOKEN_EXPIRE)
			if result.Err() != nil {
				logrus.Errorf("error expire setting session/uid to redis, error:%s", result.Err())
			}
		}
	}
}

func GetOfficialRedirectURL(c *gin.Context) {
	if sessionId, err := c.Cookie(utils.SESSION_COOKIE_NAME); err == nil {
		if len(sessionId) > 1 {
			wechat.GetOfficialRedirectURL(common.REDIRECT_URI, "", sessionId)
		}
	}
	utils.ResponseSuccessJson(c, "ok")
}

func MiniProgramAuth(c *gin.Context) {
	utils.ResponseSuccessJson(c, "ok")
}

func OfficialAccountAuth(c *gin.Context) {
	code := c.Query(common.CODE_KEY)
	state := c.Query(common.STATE_KEY)
	utils.ResponseSuccessJson(c, "ok")
	if len(code) > 0 || len(state) > 0 {
		wggo.WgGo(func() {
			if client, err := utils.GetClientInfo(c); err != nil {
				logrus.Errorf("error authenticating with officialaccount, error: %s", err)
			} else {
				client.SessionId = state
				if err := models.OfficialAccountAuth(code, client); err != nil {
					logrus.Errorf("error authenticating with officialaccount, error: %s", err)
				}
			}
		})
	} else {
		logrus.Errorf("official account auth, error code: %s or state: %s", code, state)
	}
}

func DingTalkAuth(c *gin.Context) {
	authCode := c.Query(common.AUTHCODE_KEY)
	state := c.Query(common.STATE_KEY)
	//utils.ResponseSuccessJson(c, "ok")
	utils.ResponseText(c, utils.SUCCESS_MSG_CUSTOM)
	if len(authCode) > 0 || len(state) > 0 {
		wggo.WgGo(func() {
			if client, err := utils.GetClientInfo(c); err != nil {
				logrus.Errorf("error authenticating with dingtalk, error: %s", err)
			} else {
				client.SessionId = state
				if err := models.DingTalkAuth(authCode, client); err != nil {
					logrus.Errorf("error authenticating with dingtalk, error: %s", err)
				}
			}
		})
	} else {
		logrus.Errorf("dingtalk auth, error code: %s or state: %s", authCode, state)
	}
}

func GetMiniProgromQrcode(c *gin.Context) {
	utils.ResponseSuccessJson(c, "ok")
}

func GetUserInfo(c *gin.Context) {
	if client, err := utils.GetClientInfo(c); err != nil {
		logrus.Errorf("error getting token for %s, error: %s", client.SessionId, err)
		utils.ResponseFailedJson(c, utils.ERRCODE_REQUEST_PARAM_ERROR, utils.ERRMSG_REQUEST_PARAM_ERROR, nil, http.StatusBadRequest)
	} else {
		if userinfo, err := models.GetUserInfo(client); err != nil {
			logrus.Errorf("error getting user info for %s, error: %s", client.LogFormatLong(), err)
			utils.ResponseFailedJson(c, utils.ERRCODE_INVALID_Id, utils.ERRMSG_INVALID_Id, nil, http.StatusBadGateway)
		} else {
			utils.ResponseSuccessJson(c, userinfo)
		}
	}
}

func Registry(c *gin.Context) {
	utils.ResponseSuccessJson(c, "ok")
}
