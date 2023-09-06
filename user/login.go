package user

import (
	"webb-auth/common"

	"github.com/fhy/utils-golang/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Login struct {
	gorm.Model
	UserId    int64
	LoginWith string
	IP        string
	Session   string
	OS        string
	Platform  string
	UserAgent string
}

func (login Login) log(client *utils.ClientInfo, loginWith string) {
	result := common.DB.Create(&Login{
		UserId:    client.UserId,
		LoginWith: loginWith,
		IP:        client.Ip,
		Session:   client.SessionId,
		OS:        client.Os,
		Platform:  client.Platform,
		UserAgent: client.UserAgent,
	})
	if result.Error != nil {
		logrus.Errorf("error loging user(%s)'s login,error: %s", client.LogFormatShort(), result.Error.Error())
	}
}
