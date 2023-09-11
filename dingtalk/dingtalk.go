package dingtalk

import (
	"fmt"
	"time"

	"webb-auth/common"

	"github.com/fhy/utils-golang/config"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkcontact_1_0 "github.com/alibabacloud-go/dingtalk/contact_1_0"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var getAccessTokenRequest dingtalkoauth2_1_0.GetUserTokenRequest

type DingTalk struct {
	OpenId      string `json:"id" gorm:"id primaryKey"`
	UnionId     string
	Nickname    string
	AvatarUrl   string `json:"headImageUrl"`
	PhoneNumber string
	StateCode   string
	SessionKey  string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (dk DingTalk) Create() error {
	result := common.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&dk)
	if result.Error != nil {
		err := fmt.Errorf("error upsert the dingtalk, error: %w", result.Error)
		return err
	}
	return nil
}

// This file is auto-generated, don't edit it. Thanks.

/**
 * 使用 Token 初始化账号Client
 * @return Client
 * @throws Exception
 */
func CreateClient() (_result *dingtalkoauth2_1_0.Client, _err error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	_result = &dingtalkoauth2_1_0.Client{}
	_result, _err = dingtalkoauth2_1_0.NewClient(config)
	return _result, _err
}

func ContactClient() (_result *dingtalkcontact_1_0.Client, _err error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	_result = &dingtalkcontact_1_0.Client{}
	_result, _err = dingtalkcontact_1_0.NewClient(config)
	return _result, _err
}

func GetUserInfo(accessToken string) (_dingtalk *DingTalk, _err error) {
	client, _err := ContactClient()
	if _err != nil {
		return nil, _err
	}

	getUserHeaders := &dingtalkcontact_1_0.GetUserHeaders{}
	getUserHeaders.XAcsDingtalkAccessToken = tea.String(accessToken)
	_result, _err := client.GetUserWithOptions(tea.String("me"), getUserHeaders, &util.RuntimeOptions{})

	if _err != nil {
		return nil, _err
	}

	return &DingTalk{
		OpenId:      *_result.Body.OpenId,
		UnionId:     *_result.Body.UnionId,
		Nickname:    *_result.Body.Nick,
		AvatarUrl:   *_result.Body.AvatarUrl,
		PhoneNumber: *_result.Body.Mobile,
		StateCode:   *_result.Body.StateCode,
		Email:       *_result.Body.Email,
	}, nil

}

func GetAccessToken(code string) (accessToken string, _err error) {
	client, _err := CreateClient()
	if _err != nil {
		return "", _err
	}

	getAccessTokenRequest.SetCode(code)
	result, _err := client.GetUserToken(&getAccessTokenRequest)
	if _err != nil {
		return "", _err
	}

	return *result.Body.AccessToken, nil
}

func Init(dkCfg *config.DingTalkConfig) {
	getAccessTokenRequest = dingtalkoauth2_1_0.GetUserTokenRequest{
		// ClientId:     tea.String("dingbnt5pgtoonmbjimz"),
		ClientId: tea.String(dkCfg.AppKey),
		// ClientSecret: tea.String("XPKLcyyb8_e0-5KFz-dQfXwZFqOWUBv1CnuFG5DKivSkYthj7gJNvZWHIyg5nBMx"),
		ClientSecret: tea.String(dkCfg.AppSecret),
	}
	getAccessTokenRequest.SetGrantType("authorization_code")
}
