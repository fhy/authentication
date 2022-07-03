package models

import (
	"base/utils"
	"errors"
	"fmt"
	"webb-auth/user"
	"webb-auth/wechat"
)

func MiniProgramAuth(code string, session string, client *utils.ClientInfo) (interface{}, error) {
	if w, err := wechat.GetMiniProgromUserInfo(code); err == nil {
		u, err := user.FindOrCreateUserWithWechat(w.Unionid)
		if err != nil {
			return nil, fmt.Errorf("error logining with miniprogram, error: %w", err)
		}
		if err := u.LoginWithMiniProgrom(session, client); err != nil {
			return nil, fmt.Errorf("error logining with minigrogram, error: %w", err)
		}
		if err := w.Save(); err != nil {
			return nil, fmt.Errorf("error logining with minigropram, error: %w", err)
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
			Id           int64  `json:"id"`
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}{Id: u.ID, AccessToken: accessToken, RefreshToken: refreshToken}, nil
	} else {
		return nil, fmt.Errorf("error logining with miniprogram, error: %w", err)
	}
}

func OfficialAccountAuth(code string, client *utils.ClientInfo) error {
	if userinfo, err := wechat.GetOfficialAccountUserInfo(code); err != nil {
		return fmt.Errorf("failed to get user info via officiala, code: %s, error: %s", code, err)
	} else {
		if userinfo.ErrCode == 0 {
			u, err := user.FindOrCreateUserWithWechat(userinfo.Unionid)
			if err != nil {
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
				Sex:       userinfo.Sex,
				Province:  userinfo.Province,
				City:      userinfo.City,
				Country:   userinfo.Country,
			}
			if err := w.Save(); err != nil {
				return fmt.Errorf("error logining with create wechat, error: %w", err)
			}
			return nil
		} else {
			return errors.New(userinfo.ErrMsg)
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
