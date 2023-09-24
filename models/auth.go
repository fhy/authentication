package models

import (
	"errors"
	"fmt"
	"webb-auth/dingtalk"
	"webb-auth/user"
	"webb-auth/wechat"

	"github.com/fhy/utils-golang/utils"
)

func OfficialAccountAuth(code string, client *utils.ClientInfo) error {
	if userinfo, err := wechat.GetOfficialAccountUserInfo(code); err != nil {
		return fmt.Errorf("failed to get user info via officiala, code: %s, error: %s", code, err)
	} else {
		if userinfo.ErrCode == 0 {
			u := user.User{}
			if err := u.FindWithWechat(userinfo.Unionid); err != nil {
				return fmt.Errorf("error logining with officeaccount, error: %w", err)
			}
			if err := u.LoginWithOfficeAccount(client); err != nil {
				return fmt.Errorf("error logining with officeaccount, error: %w", err)
			}
			w := wechat.WeChat{
				OpenId:    userinfo.OpenID,
				UnionId:   userinfo.Unionid,
				Nickname:  userinfo.Nickname,
				AvatarUrl: userinfo.HeadImgURL,
			}
			if err := w.Create(); err != nil {
				return fmt.Errorf("error logining with create wechat, error: %w", err)
			}
			return nil
		} else {
			return errors.New(userinfo.ErrMsg)
		}
	}
}

func DingTalkAuth(code string, client *utils.ClientInfo) error {
	if accessToken, err := dingtalk.GetAccessToken(code); err != nil {
		return fmt.Errorf("failed to get access token via dingtalk, code: %s, error: %s", code, err)
	} else {
		if userinfo, err := dingtalk.GetUserInfoViaToken(accessToken); err != nil {
			return fmt.Errorf("failed to get user info via dingtalk, code: %s, error: %s", code, err)
		} else {
			u := user.User{}
			if err := u.FindWithDingtalk(userinfo.OpenId); err != nil {
				return fmt.Errorf("error logining with dingtalk, error: %w", err)
			}
			if err := u.LoginWithDingtalk(client); err != nil {
				return fmt.Errorf("error logining with dingtalk, error: %w", err)
			}
			if err := userinfo.Create(); err != nil {
				return fmt.Errorf("error logining with create dingtalk, error: %w", err)
			}
			return nil
		}
	}
}

func GetTokenWithSession(client *utils.ClientInfo) (interface{}, error) {
	u := user.User{}
	if err := u.GetFromSession(client.SessionId); err != nil {
		return nil, fmt.Errorf("error geting token with session, error: %w", err)
	}
	accessToken, err := u.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("error geting token with session, error: %w", err)
	}
	refreshToken, err := u.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("error geting token with session, error: %w", err)
	}
	return &struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func RefleshToken(token string, client *utils.ClientInfo) (interface{}, error) {
	u := user.User{}
	accessToken, err := u.RefleshToken(token)
	if err != nil {
		return nil, fmt.Errorf("error refreshing token, error: %w", err)
	}
	return &struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{AccessToken: accessToken, RefreshToken: token}, nil
}

func GetUserInfo(_client *utils.ClientInfo) (interface{}, error) {
	if _client.UserId > 0 {
		_user, err := user.GetFromId(_client.UserId)
		if err != nil {
			return nil, fmt.Errorf("error get userinfo, error: %s", err)
		} else {
			if _user.DingTalkId != "" {
				_dingtalk, err := dingtalk.GetFromId(_user.DingTalkId)
				if err == nil {
					return &struct {
						Nick      string
						Mobile    string
						AvatarUrl string
					}{_dingtalk.Nickname, _dingtalk.PhoneNumber, _dingtalk.AvatarUrl}, nil
				}
			}
			return nil, fmt.Errorf("error get userinfo, error: %s", err)
		}
	} else {
		return nil, fmt.Errorf("error get userinfo, error id: %d", _client.UserId)
	}
}
