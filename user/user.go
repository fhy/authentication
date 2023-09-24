package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"webb-auth/common"
	"webb-auth/conf"

	"github.com/fhy/utils-golang/utils"
	"github.com/fhy/utils-golang/wggo"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	ID          int64 `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LoginAt     time.Time
	LogoutAt    time.Time
	Password    string
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Email       string         `json:"email"`
	MobilePhone string
	WeChatId    string
	DingTalkId  string
}

func (u *User) Init() {
	fmt.Println("init user")
}

func (u *User) FindWithWechat(wechatId string) error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		err = fmt.Errorf("error getting the user with wechat, error: %w", err)
		return err
	}
	res, err := password.Generate(32, 8, 8, false, false)
	if err != nil {
		err = fmt.Errorf("error getting the user with wechat, error: %w", err)
		return err
	}
	result := common.DB.Where(User{WeChatId: wechatId}).Attrs(User{ID: node.Generate().Int64(), Password: res}).FirstOrCreate(u)

	if result.Error != nil {
		err = fmt.Errorf("error getting the user with wechat, error: %w", result.Error)
		return err
	}
	return nil
}

func (u *User) FindWithDingtalk(dingtalkId string) error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		err = fmt.Errorf("error getting the user with dingtalk, error: %w", err)
		return err
	}
	res, err := password.Generate(32, 8, 8, false, false)
	if err != nil {
		err = fmt.Errorf("error getting the user with dingtalk, error: %w", err)
		return err
	}
	result := common.DB.Where(User{DingTalkId: dingtalkId}).Attrs(User{ID: node.Generate().Int64(), Password: res}).FirstOrCreate(u)

	if result.Error != nil {
		err = fmt.Errorf("error getting the user with dingtalk, error: %w", result.Error)
		return err
	}
	return nil
}

func (u *User) GenerateToken() (string, error) {
	nowTime := time.Now()
	expireTime := jwt.NewNumericDate(nowTime.Add(utils.TOKEN_EXPIRE))
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", fmt.Errorf("error generating user token, error: %w", err)
	}

	claims := utils.AccessClaims{
		UID:       u.ID,
		TokenType: utils.TOKEN_TYPE_ACCESS,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间
			ExpiresAt: expireTime,
			IssuedAt:  jwt.NewNumericDate(nowTime),
			ID:        node.Generate().String(),
			// 指定token发行人
			Issuer: "138e8",
		},
	}

	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	// token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(conf.Conf.Jwt.Secret)
	token, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(utils.LoadEdPrivateKeyFromDisk(conf.Conf.Jwt.PrivKeyPath))

	if err != nil {
		return "", fmt.Errorf("error creating the token for user: %d, error:%w", u.ID, err)
	}
	return token, nil
}

func GetFromId(id int64) (_user *User, _err error) {
	_user = &User{}
	result := common.DB.First(_user, id)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting user from id: %d, error: %w", id, result.Error)
	}
	return _user, nil
}

func (u *User) GenerateRefreshToken() (string, error) {
	nowTime := time.Now()
	expireTime := jwt.NewNumericDate(nowTime.Add(utils.TOKEN_EXPIRE))

	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", fmt.Errorf("error generating user refresh token, error: %w", err)
	}

	claims := utils.RefreshClaims{
		UID:       u.ID,
		TokenType: utils.TOKEN_TYPE_REFRESH,
		Salt:      u.LoginAt.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间
			ExpiresAt: expireTime,
			IssuedAt:  jwt.NewNumericDate(nowTime),
			ID:        node.Generate().String(),
			// 指定token发行人
			Issuer: "138e8",
		},
	}

	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	// token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(fmt.Sprintf("%s:%s:%d", conf.Conf.Jwt.Secret, u.Password, u.LogoutAt.Unix()))
	token, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(utils.LoadEdPrivateKeyFromDisk(conf.Conf.Jwt.PrivKeyPath))
	if err != nil {
		return "", fmt.Errorf("error creating the token for user: %d, error:%w", u.ID, err)
	}
	return token, nil
}

func (u *User) RefleshToken(t string) (string, error) {
	pubkey := utils.LoadEdPublicKeyFromDisk(conf.Conf.Jwt.PubKeyPath)
	token, err := utils.VerifyRefleshToken(t, &pubkey)
	if err != nil {
		return "", fmt.Errorf("error refreshing token for user:%d, error: %w", u.ID, err)
	}
	if err := u.FindWithID(token.Claims.(*utils.RefreshClaims).UID); err != nil {
		return "", fmt.Errorf("error refreshing token for user:%d, error: %w", u.ID, err)
	}
	if token.Claims.(*utils.RefreshClaims).Salt != u.LogoutAt.Unix() {
		return "", errors.New("error user logout")
	}
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", fmt.Errorf("error refreshing token for user:%d, error: %w", u.ID, err)
	}
	nowTime := time.Now()
	expireTime := jwt.NewNumericDate(nowTime.Add(utils.TOKEN_EXPIRE))

	claims := utils.AccessClaims{
		UID:       u.ID,
		TokenType: utils.TOKEN_TYPE_ACCESS,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间
			ExpiresAt: expireTime,
			IssuedAt:  jwt.NewNumericDate(nowTime),
			ID:        node.Generate().String(),
			// 指定token发行人
			Issuer: conf.Conf.Jwt.Issuer,
		},
	}

	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(utils.LoadEdPrivateKeyFromDisk(conf.Conf.Jwt.PrivKeyPath))
	if err != nil {
		return "", fmt.Errorf("error generate token when refreshing token for user:%d, error: %w", u.ID, err)
	}
	return accessToken, nil
}

func (u *User) LoginWithOfficeAccount(client *utils.ClientInfo) error {
	result := common.RC.Set(context.Background(), fmt.Sprintf("%s%s", utils.USER_SID_REDIS_PREFIX, client.SessionId), u.ID, utils.TOKEN_EXPIRE)
	if result.Err() != nil {
		return fmt.Errorf("error login setting session/uid to redis, error:%w", result.Err())
	}
	wggo.WgGo(func() {
		if err := u.updateLogin(); err != nil {
			logrus.Errorf("error %s logining error: %s", client.LogFormatShort(), err)
		}
		client.UserId = u.ID
		Login{}.log(client, common.LOGIN_WITH_WECHAT_OFFICIALACCOUNT)
	})
	return nil
}

func (u *User) LoginWithDingtalk(client *utils.ClientInfo) error {
	result := common.RC.Set(context.Background(), fmt.Sprintf("%s%s", utils.USER_SID_REDIS_PREFIX, client.SessionId), u.ID, utils.TOKEN_EXPIRE)
	if result.Err() != nil {
		return fmt.Errorf("error login setting session/uid to redis, error:%w", result.Err())
	}
	wggo.WgGo(func() {
		if err := u.updateLogin(); err != nil {
			logrus.Errorf("error %s logining error: %s", client.LogFormatShort(), err)
		}
		client.UserId = u.ID
		Login{}.log(client, common.LOGIN_WITH_DINGTALK)
	})
	return nil
}

func (u *User) GetFromSession(session string) error {
	if len(session) < 1 {
		return errors.New(utils.ERRMSG_INVALID_SESSION)
	}
	result, err := common.RC.Get(context.Background(), fmt.Sprintf("%s%s", utils.USER_SID_REDIS_PREFIX, session)).Result()
	if err != nil {
		logrus.Debugf("error getting user, getting id with session:%s from redis, error:%s", session, err)
		return fmt.Errorf("error getting user from session, error:%w", err)
	}
	if id, err := strconv.ParseInt(result, 10, 64); err != nil {
		logrus.Debugf("error getting user, getting id with session: %s from string:%s, error:%s", session, result, err)
		return fmt.Errorf("error getting user from session, error:%w", err)
	} else {
		if id < 1 {
			logrus.Debugf("error getting user, getting error id:%d with session: %s", id, session)
			return fmt.Errorf("error getting user from session, error:%w", errors.New(utils.ERRMSG_INVALID_Id))
		}
		if err := u.FindWithID(id); err != nil {
			logrus.Debugf("error getting user with session: %s from id:%d, error:%s", session, id, err)
			return fmt.Errorf("error getting user from session, error:%w", err)
		}
		return nil
	}
}

func (u *User) Logout(client *utils.ClientInfo) error {
	if len(client.SessionId) < 1 {
		return errors.New(utils.ERRMSG_INVALID_SESSION)
	}
	result := common.RC.Del(context.Background(), fmt.Sprintf("%s%s", utils.USER_SID_REDIS_PREFIX, client.SessionId))
	if result.Err() != nil {
		return fmt.Errorf("error login setting session/uid to redis, error:%w", result.Err())
	}
	wggo.WgGo(func() {
		if err := u.updateLogout(); err != nil {
			logrus.Errorf("error %s logouting error: %s", client.LogFormatShort(), err)
		}
		client.UserId = u.ID
		Login{}.log(client, common.LOGIN_WITH_LOGOUG)
	})
	return nil
}

func (u *User) FindWithID(id int64) error {
	result := common.DB.First(u, id)
	if result.Error != nil {
		err := fmt.Errorf("error getting the user with wechat, error: %w", result.Error)
		return err
	}
	return nil
}

// func (u *User) Create() error {
// 	result := common.DB.Create(u)
// 	if result.Error != nil {
// 		err := fmt.Errorf("error creating the user to db, error: %w", result.Error)
// 		return err
// 	}
// 	return nil
// }

func (u *User) updateLogin() error {
	result := common.DB.Model(u).Update("login_at", time.Now())
	if result.Error != nil {
		err := fmt.Errorf("error updating the user's login time, error: %w", result.Error)
		return err
	}
	return nil
}

func (u *User) updateLogout() error {
	result := common.DB.Model(u).Update("logout_at", time.Now())
	if result.Error != nil {
		err := fmt.Errorf("error updating the user's login time, error: %w", result.Error)
		return err
	}
	return nil
}

func (u *User) SaveMobilePhone(phone string) error {
	result := common.DB.Model(u).Update("mobile_phone", phone)
	if result.Error != nil {
		err := fmt.Errorf("error saving the user's mobile phone to db, error: %w", result.Error)
		return err
	}
	return nil
}
func (u *User) SaveEmail(email string) error {
	result := common.DB.Model(u).Update("email", email)
	if result.Error != nil {
		err := fmt.Errorf("error saving the user's email to db, error: %w", result.Error)
		return err
	}
	return nil
}

func (u *User) Save() error {
	result := common.DB.Save(u)
	if result.Error != nil {
		err := fmt.Errorf("error saving the user to db, error: %w", result.Error)
		return err
	}
	return nil
}
