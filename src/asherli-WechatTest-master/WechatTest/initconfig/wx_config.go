/*
 *    Copyright 2016-2018 Li ZongZe
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package initconfig

import (
	"errors"
	"fmt"
	"sync"
	"time"

	. "weixin/utils"
)

var WxConfig WxOpenConfigInterface

//初始化微信配置参数
func init() {
	var err error
	WxConfig, err = NewWxOpenConfig("wxbdcb6afc21014792", "wxce28e7deb273a7ae", "c6b5dfad77e8165b156deb5b056f702c", "ZrH0uLHcUeF1hdj5Y6hL8CvnSwfF92ZwPC1bHNAjAeJ", "initiateWeb", "wxc44e124ba5d0d053", "ea86c29d5c74c39773ffd69e992d1000")
	if err != nil {
		fmt.Println(err)
		return
	}
}

const (
	//微信公众号接口令牌
	ACCESS_TOKEN_KEY              = "WX_OPEN_ACCESS_TOKEN"
	ACCESS_REFRESH_TOKEN_KEY      = "WX_OPEN_ACCESS_REFRESH_TOKEN"
	ACCESS_TOKEN_EXPIRES_TIME_KEY = "WX_OPEN_ACCESS_TOKEN_EXPIRES_TIME"

	//开放平台接口令牌
	COMPONENT_ACCESS_TOKEN_KEY              = "WX_COMPONENT_ACCESS_TOKEN"
	COMPONENT_ACCESS_TOKEN_EXPIRES_TIME_KEY = "WX_COMPONENT_ACCESS_TOKEN_EXPIRES_TIME"

	COMPONENT_VERIFY_YICKIT = "WX_COMPONENT_VERIFY_TICKET"

	//网页JSSDK
	JS_API_TICKET_KEY              = "WX_OPEN_JS_API_TICKET"
	JS_API_TICKET_EXPIRES_TIME_KEY = "WX_OPEN_JS_API_TICKET_EXPIRES_TIME"
)

type WxOpenConfig struct {
	appId string

	accessToken             string
	accessTokenExpiresTime  int64
	accessTokenRefreshToken string
	accessTokenLock         sync.RWMutex

	jsApiTicket            string
	jsApiTicketExpiresTime int64
	jsApiTicketLock        sync.RWMutex

	componentAppId          string
	componentAppSecret      string
	componentEncodingAesKey string
	componentToken          string

	componentVerifyTicket string

	componentAccessToken            string
	componentAccessTokenExpiresTime int64
	componentAccessTokenLock        sync.RWMutex

	miniAppId     string
	miniAppSecret string
}

func (w *WxOpenConfig) MiniAppSecret() string {
	return w.miniAppSecret
}

func (w *WxOpenConfig) MiniAppId() string {
	return w.miniAppId
}

type WxOpenConfigInterface interface {
	AppId() string
	ComponentAppId() string
	ComponentAppSecret() string
	ComponentEncodingAesKey() string
	ComponentToken() string

	//开放平台接口令牌兑换凭证
	ComponentVerifyTicket() string
	UpdateComponentVerifyTicket(verifyTicket string) error

	//开放平台接口令牌
	ComponentAccessTokenLock() sync.RWMutex
	ComponentAccessToken() string
	UpdateComponentAccessToken(componentAccessToken string, componentAccessTokenExpiresTime int64) error
	ExpireComponentAccessToken() error
	IsComponentAccessTokenExpired() bool

	//公众号接口令牌
	AccessToken() string
	AccessTokenLock() sync.RWMutex
	ExpireAccessToken() error
	IsAccessTokenExpired() bool
	UpdataAccessToken(accessToken, refreshToken string, accessTokenExpiresTime int64) error
	AccessTokenRefreshToken() string

	//JSSDK
	JsApiTicket() string
	JsApiTicketLock() sync.RWMutex
	UpdateJsApiTicket(jsapiTicket string, expiresInSeconds int64) error
	ExpireJsApiTicket() error
	IsJsApiTicketExpired() bool

	MiniAppSecret() string
	MiniAppId() string
}

func NewWxOpenConfig(appId, componentAppId, componentAppSecret, componentEncodingAesKey, componentToken, miniAppId, miniAppSecret string) (WxOpenConfigInterface, error) {
	if IsEmptyStr(appId) {
		return nil, errors.New("公众号AppId不能为空")
	}
	if IsEmptyStr(componentAppId) {
		return nil, errors.New("开放平台AppId不能为空")
	}
	if IsEmptyStr(componentAppSecret) {
		return nil, errors.New("开放平台AppSecret不能为空")
	}
	if IsEmptyStr(componentEncodingAesKey) {
		return nil, errors.New("开放平台EncodingAesKey不能为空")
	}
	if IsEmptyStr(componentToken) {
		return nil, errors.New("开放平台Token不能为空")
	}
	if IsEmptyStr(miniAppId) {
		return nil, errors.New("小程序AppId不能为空")
	}
	if IsEmptyStr(miniAppSecret) {
		return nil, errors.New("小程序AppSecret不能为空")
	}

	return &WxOpenConfig{
		appId:                   appId,
		componentAppId:          componentAppId,
		componentAppSecret:      componentAppSecret,
		componentEncodingAesKey: componentEncodingAesKey,
		componentToken:          componentToken,
		miniAppId:               miniAppId,
		miniAppSecret:           miniAppSecret,
	}, nil
}

func (w *WxOpenConfig) ComponentToken() string {
	return w.componentToken
}

func (w *WxOpenConfig) ComponentEncodingAesKey() string {
	return w.componentEncodingAesKey
}

func (w *WxOpenConfig) ComponentAppSecret() string {
	return w.componentAppSecret
}

func (w *WxOpenConfig) ComponentAppId() string {
	return w.componentAppId
}

func (w *WxOpenConfig) ComponentVerifyTicket() string {
	componentVerifyTicket, err := Redis.Get(COMPONENT_VERIFY_YICKIT).Result()
	if err != nil {
		return ""
	}
	return componentVerifyTicket
}

func (w *WxOpenConfig) ComponentAccessTokenLock() sync.RWMutex {
	return w.componentAccessTokenLock
}

func (w *WxOpenConfig) ComponentAccessToken() string {
	componentAccessToken, err := Redis.Get(COMPONENT_ACCESS_TOKEN_KEY).Result()
	if err != nil {
		return ""
	}
	return componentAccessToken
}

func (w *WxOpenConfig) JsApiTicketLock() sync.RWMutex {
	return w.jsApiTicketLock
}

func (w *WxOpenConfig) JsApiTicket() string {
	jsApiTicket, err := Redis.Get(JS_API_TICKET_KEY).Result()
	if err != nil {
		return ""
	}
	return jsApiTicket
}

func (w *WxOpenConfig) AccessTokenLock() sync.RWMutex {
	return w.accessTokenLock
}

func (w *WxOpenConfig) AccessTokenRefreshToken() string {
	accessTokenRefreshToken, err := Redis.Get(ACCESS_REFRESH_TOKEN_KEY).Result()
	if err != nil {
		return ""
	}
	return accessTokenRefreshToken
}

func (w *WxOpenConfig) AccessToken() string {
	accessToken, err := Redis.Get(ACCESS_TOKEN_KEY).Result()
	if err != nil {
		return ""
	}
	return accessToken
}

func (w *WxOpenConfig) AppId() string {
	return w.appId
}

func (woc *WxOpenConfig) UpdateComponentVerifyTicket(verifyTicket string) error {
	return Redis.Set(COMPONENT_VERIFY_YICKIT, verifyTicket, 0).Err()
}

func (woc *WxOpenConfig) UpdateComponentAccessToken(componentAccessToken string, componentAccessTokenExpiresTime int64) error {
	err := Redis.Set(COMPONENT_ACCESS_TOKEN_KEY, componentAccessToken, 0).Err()
	if err != nil {
		return err
	}
	return Redis.Set(COMPONENT_ACCESS_TOKEN_EXPIRES_TIME_KEY, time.Now().UnixNano()/1e6+(componentAccessTokenExpiresTime-200)*1000, 0).Err()
}

func (woc *WxOpenConfig) ExpireComponentAccessToken() error {
	return Redis.Set(COMPONENT_ACCESS_TOKEN_EXPIRES_TIME_KEY, 0, 0).Err()
}

func (woc *WxOpenConfig) IsComponentAccessTokenExpired() bool {
	componentAccessTokenExpiresTime, err := Redis.Get(COMPONENT_ACCESS_TOKEN_EXPIRES_TIME_KEY).Int64()
	if err != nil {
		return true
	}
	if componentAccessTokenExpiresTime < time.Now().UnixNano()/1e6 {
		return true
	}
	return false
}

func (woc *WxOpenConfig) UpdataAccessToken(accessToken, refreshToken string, accessTokenExpiresTime int64) error {
	err := Redis.Set(ACCESS_TOKEN_KEY, accessToken, 0).Err()
	if err != nil {
		return err
	}
	err = Redis.Set(ACCESS_REFRESH_TOKEN_KEY, refreshToken, 0).Err()
	if err != nil {
		return err
	}
	return Redis.Set(ACCESS_TOKEN_EXPIRES_TIME_KEY, time.Now().UnixNano()/1e6+(accessTokenExpiresTime-200)*1000, 0).Err()
}

func (woc *WxOpenConfig) ExpireAccessToken() error {
	return Redis.Set(ACCESS_TOKEN_EXPIRES_TIME_KEY, 0, 0).Err()
}

func (woc *WxOpenConfig) IsAccessTokenExpired() bool {
	componentAccessTokenExpiresTime, err := Redis.Get(ACCESS_TOKEN_EXPIRES_TIME_KEY).Int64()
	if err != nil {
		return true
	}
	if componentAccessTokenExpiresTime < time.Now().UnixNano()/1e6 {
		return true
	}
	return false
}

func (woc *WxOpenConfig) UpdateJsApiTicket(jsapiTicket string, expiresInSeconds int64) error {
	err := Redis.Set(JS_API_TICKET_KEY, jsapiTicket, 0).Err()
	if err != nil {
		return err
	}
	return Redis.Set(JS_API_TICKET_EXPIRES_TIME_KEY, time.Now().UnixNano()/1e6+(expiresInSeconds-200)*1000, 0).Err()
}

func (woc *WxOpenConfig) ExpireJsApiTicket() error {
	return Redis.Set(JS_API_TICKET_EXPIRES_TIME_KEY, 0, 0).Err()
}

func (woc *WxOpenConfig) IsJsApiTicketExpired() bool {
	expiresInSeconds, err := Redis.Get(JS_API_TICKET_EXPIRES_TIME_KEY).Int64()
	if err != nil {
		return true
	}
	if expiresInSeconds < time.Now().UnixNano()/1e6 {
		return true
	}
	return false
}
