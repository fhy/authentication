package dingtalk

import (
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkcontact_1_0 "github.com/alibabacloud-go/dingtalk/contact_1_0"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"gorm.io/gorm"
)

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

func GetUserInfo(accessToken string) (_err error) {
	client, _err := ContactClient()
	if _err != nil {
		return _err
	}

	getUserHeaders := &dingtalkcontact_1_0.GetUserHeaders{}
	getUserHeaders.XAcsDingtalkAccessToken = tea.String(accessToken)
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, _err = client.GetUserWithOptions(tea.String("me"), getUserHeaders, &util.RuntimeOptions{})
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			err = _t
		} else {
			err.Message = tea.String(tryErr.Error())
		}
		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {
			// err 中含有 code 和 message 属性，可帮助开发定位问题
		}

	}
	return _err
}

func GetAccessToken(args []*string) (_err error) {
	client, _err := CreateClient()
	if _err != nil {
		return _err
	}

	getAccessTokenRequest := &dingtalkoauth2_1_0.GetAccessTokenRequest{
		AppKey:    tea.String("dingbnt5pgtoonmbjimz"),
		AppSecret: tea.String("XPKLcyyb8_e0-5KFz-dQfXwZFqOWUBv1CnuFG5DKivSkYthj7gJNvZWHIyg5nBMx"),
	}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, _err = client.GetAccessToken(getAccessTokenRequest)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			err = _t
		} else {
			err.Message = tea.String(tryErr.Error())
		}
		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {
			// err 中含有 code 和 message 属性，可帮助开发定位问题
		}

	}
	return _err
}
