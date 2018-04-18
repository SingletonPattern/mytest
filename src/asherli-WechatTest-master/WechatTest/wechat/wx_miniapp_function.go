package wechat

import (
	"fmt"

	"weixin/utils/xrequest"

	"github.com/labstack/gommon/log"
)

const (
	jsCode2SessionKey = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type MiniAppSessionKey struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
}

func (c *wechatConfig)MiniAppLogin(jsCode string) (*MiniAppSessionKey, []error) {
	result := &MiniAppSessionKey{}

	resp, _, errs := xrequest.New().Post(fmt.Sprintf(jsCode2SessionKey, c.config.MiniAppId(), c.config.MiniAppSecret(), jsCode)).EndStruct(result)
	if len(errs) > 0 {
		return nil, errs
	}
	err := unmarshalResponseToJson(resp, result)
	if err != nil {
		log.Error("小程序登录失败: ", err)
		return nil, []error{err}
	}
	//TODO Session存储当前不实现
	return result,nil
}
