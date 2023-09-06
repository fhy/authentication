package wechat

import (
	"context"
	"fmt"
	"time"
	"webb-auth/common"
	"webb-auth/conf"

	"github.com/fhy/utils-golang/config"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	mpcfg "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	oacfg "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	officalAccount *officialaccount.OfficialAccount
	miniProgram    *miniprogram.MiniProgram
)

type WeChat struct {
	OpenId      string `json:"id" gorm:"id primaryKey"`
	UnionId     string
	Province    string
	City        string
	Country     string
	Unionid     string
	Sex         int32
	Nickname    string
	AvatarUrl   string `json:"headImageUrl"`
	PhoneNumber string
	SessionKey  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (w WeChat) Save() error {
	result := common.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&w)
	if result.Error != nil {
		err := fmt.Errorf("error upsert the wechat, error: %w", result.Error)
		return err
	}
	return nil
}

func GetOfficialRedirectURL(redirectURL string, scope string, state string) (string, error) {
	return officalAccount.GetOauth().GetRedirectURL(redirectURL, scope, state)
}

func GetOfficialAccountUserInfo(code string) (*oauth.UserInfo, error) {
	auth := officalAccount.GetOauth()
	if accessToken, err := auth.GetUserAccessToken(code); err != nil {
		return nil, fmt.Errorf("failed to get officialaccount user accesstoke, code: %s, error: %w", code, err)
	} else {
		if userinfo, err := auth.GetUserInfo(accessToken.AccessToken, accessToken.OpenID, ""); err != nil {
			return nil, fmt.Errorf("failed to get officialaccount user info, accessToken: %s,, openid: %s error: %w",
				accessToken.AccessToken, accessToken.OpenID, err)
		} else {
			return &userinfo, nil
		}
	}
}

func GetMiniProgromUserInfo(code string) (*WeChat, error) {
	auth := miniProgram.GetAuth()
	if userinfo, err := auth.Code2Session(code); err == nil {
		if userinfo.ErrCode == 0 {
			return &WeChat{OpenId: userinfo.OpenID, UnionId: userinfo.UnionID, SessionKey: userinfo.SessionKey}, err
		} else {
			return nil, fmt.Errorf("failed to get MiniProgrom user info, code: %s error: %s", code, userinfo.ErrMsg)
		}
	} else {
		return nil, fmt.Errorf("failed to get MiniProgrom user info, code: %s error: %w", code, err)
	}
}

func Init(redisCfg *config.RedisConfig, wcCfg *conf.WechatConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wc := wechat.NewWechat()
	wc.SetCache(cache.NewRedis(ctx, &cache.RedisOpts{
		Host:        fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:    redisCfg.Password,
		Database:    redisCfg.DB,
		MaxIdle:     redisCfg.MaxIdle,
		MaxActive:   redisCfg.MaxActive,
		IdleTimeout: redisCfg.IdleTimeout,
	}))
	officalAccount = wc.GetOfficialAccount(&oacfg.Config{
		AppID:     wcCfg.OfficialAccount.AppID,
		AppSecret: wcCfg.OfficialAccount.AppSecret,
		Token:     wcCfg.OfficialAccount.Token,
	})
	miniProgram = wc.GetMiniProgram(&mpcfg.Config{
		AppID:     wcCfg.MiniProgram.AppID,
		AppSecret: wcCfg.MiniProgram.AppSecret,
	})
}
