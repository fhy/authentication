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
	"github.com/silenceper/wechat/v2/officialaccount"
	oacfg "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	officalAccount *officialaccount.OfficialAccount
)

type WeChat struct {
	OpenId      string `json:"id" gorm:"id primaryKey"`
	UnionId     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Nickname    string
	AvatarUrl   string `json:"headImageUrl"`
	PhoneNumber string
}

func (w WeChat) Create() error {
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
		logrus.Errorf("failed to get officialaccount user accesstoke, code: %s, error: %s", code, err)
		return nil, err
	} else {
		if userinfo, err := auth.GetUserInfo(accessToken.AccessToken, accessToken.OpenID, ""); err != nil {
			logrus.Errorf("failed to get officialaccount user info, accessToken: %s,, openid: %s error: %s",
				accessToken.AccessToken, accessToken.OpenID, err)
			return nil, err
		} else {
			return &userinfo, nil
		}
	}
}

func Init(redisCfg *config.RedisConfig, wcCfg *conf.WechatConfig) {
	wc := wechat.NewWechat()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	officalAccount = wc.GetOfficialAccount(&oacfg.Config{
		AppID:     wcCfg.OfficialAccount.AppID,
		AppSecret: wcCfg.OfficialAccount.AppSecret,
		Token:     wcCfg.OfficialAccount.Token,
	})
	wc.SetCache(cache.NewRedis(ctx, &cache.RedisOpts{
		Host:        fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:    redisCfg.Password,
		Database:    redisCfg.DB,
		MaxActive:   redisCfg.MaxActive,
		IdleTimeout: redisCfg.IdleTimeout,
	}))
}
