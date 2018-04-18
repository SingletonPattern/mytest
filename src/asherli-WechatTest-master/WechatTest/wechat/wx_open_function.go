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

package wechat

import (
	"fmt"
	"net/url"

	"weixin/utils/xrequest"

	"github.com/labstack/gommon/log"
)

const (
	weixinComponentHost    = "https://api.weixin.qq.com/cgi-bin/component"
	apiComponentToken      = weixinComponentHost + "/api_component_token"
	apiCreatePreAuthCode   = weixinComponentHost + "/api_create_preauthcode?component_access_token=%s"
	apiQueryAuth           = weixinComponentHost + "/api_query_auth?component_access_token=%s"
	apiAuthorizerToken     = weixinComponentHost + "/api_authorizer_token?component_access_token=%s"
	apiGetAuthorizerInfo   = weixinComponentHost + "/api_get_authorizer_info?component_access_token=%s"
	apiGetAuthorizerOption = weixinComponentHost + "/api_get_authorizer_option?component_access_token=%s"
	apiSetAuthorizerOption = weixinComponentHost + "/api_set_authorizer_option?component_access_token=%s"

	oauthUrl = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s"
)

type ComponentAccessTokenRequest struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppSecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

type ComponentAccessToken struct {
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int64  `json:"expires_in"`
}

//获取第三方平台component_access_token
func (c *wechatConfig) GetComponentAccessToken() string {
	componentAccessTokenLock := c.config.ComponentAccessTokenLock()
	componentAccessTokenLock.Lock()
	defer componentAccessTokenLock.Unlock()
	if c.config.IsComponentAccessTokenExpired() {
		postData := ComponentAccessTokenRequest{
			ComponentAppid:        c.config.ComponentAppId(),
			ComponentAppSecret:    c.config.ComponentAppSecret(),
			ComponentVerifyTicket: c.config.ComponentVerifyTicket(),
		}
		result := &ComponentAccessToken{}
		resp, _, errs := xrequest.New().Post(apiComponentToken).SendStruct(postData).EndStruct(result)
		if errs != nil {
			log.Error("获取开放平台接口令牌失败: ", errs)
			return ""
		}
		err := unmarshalResponseToJson(resp, result)
		if err != nil {
			log.Error("Token数据解析失败: ", err)
			return ""
		}
		err = c.config.UpdateComponentAccessToken(result.ComponentAccessToken, result.ExpiresIn);
		if err != nil {
			log.Error("开放平台接口令牌保存失败: ", err)
			return ""
		}
	}
	return c.config.ComponentAccessToken()
}

//刷新第三方平台component_access_token
func (c *wechatConfig) RefreshComponentAccessToken(forceRefresh bool) string {
	refreshLock := c.config.ComponentAccessTokenLock()
	refreshLock.Lock()
	defer refreshLock.Unlock()
	if forceRefresh {
		c.config.ExpireComponentAccessToken();
	}
	return WxFunction.GetComponentAccessToken()
}

//获取预授权码pre_auth_code
func (c *wechatConfig) GetPreAuthCode() (string, []error) {
	postData := struct {
		ComponentAppid string `json:"component_appid"`
	}{
		ComponentAppid: c.config.ComponentAppId(),
	}
	result := &struct {
		PreAuthCode string  `json:"pre_auth_code"`
		ExpiresIn   float64 `json:"expires_in"`
	}{}
	resp, _, errs := xrequest.New().Post(fmt.Sprintf(apiCreatePreAuthCode, WxFunction.GetComponentAccessToken())).SendStruct(postData).EndStruct(result)
	if errs != nil {
		log.Error("获取预授权码失败", errs)
		return "", errs
	}
	err := unmarshalResponseToJson(resp, result)
	if err != nil {
		log.Error("解析预授权码失败", err)
		return "", []error{err}
	}
	return result.PreAuthCode, nil
}

type ApiQueryAuth struct {
	AuthorizationInfo struct {
		AppId        string `json:"authorizer_appid"`
		AccessToken  string `json:"authorizer_access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"authorizer_refresh_token"`
		FuncInfo []struct {
			Funcscope struct {
				Id int64 `json:"id"`
			} `json:"funcscope_category"`
		} `json:"func_info"`
	} `json:"authorization_info"`
}

//使用授权码换取公众号或小程序的接口调用凭据和授权信息
func (c *wechatConfig) GetAuthorization(authCode string) (*ApiQueryAuth, []error) {
	postData := struct {
		ComponentAppid    string `json:"component_appid"`
		AuthorizationCode string `json:"authorization_code"`
	}{
		ComponentAppid:    c.config.ComponentAppId(),
		AuthorizationCode: authCode,
	}
	result := &ApiQueryAuth{}
	resp, _, errs := xrequest.New().Post(fmt.Sprintf(apiQueryAuth, WxFunction.GetComponentAccessToken())).SendStruct(postData).EndStruct(result)
	if errs != nil {
		log.Error("获取授权信息失败", errs)
		return nil, errs
	}
	err := unmarshalResponseToJson(resp, result)
	if err != nil {
		log.Error("授权信息解析失败", err)
		return nil, []error{err}
	}
	err = c.config.UpdataAccessToken(result.AuthorizationInfo.AccessToken, result.AuthorizationInfo.RefreshToken, result.AuthorizationInfo.ExpiresIn)
	if err != nil {
		log.Error("授权信息缓存失败", err)
		return nil, []error{err}
	}
	return result, nil
}

type ApiAuthorizerToken struct {
	AccessToken  string `json:"authorizer_access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"authorizer_refresh_token"`
}

//获取（刷新）授权公众号或小程序的接口调用凭据（令牌）
func (c *wechatConfig) GetAuthAccessToken() string {
	accessTokenLock := c.config.AccessTokenLock()
	accessTokenLock.Lock()
	defer accessTokenLock.Unlock()
	if c.config.IsAccessTokenExpired() {
		postData := struct {
			ComponentAppid         string `json:"component_appid"`
			AuthorizerAppid        string `json:"authorizer_appid"`
			AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
		}{
			ComponentAppid:         c.config.ComponentAppId(),
			AuthorizerAppid:        c.config.AppId(),
			AuthorizerRefreshToken: c.config.AccessTokenRefreshToken(),
		}
		result := &ApiAuthorizerToken{}
		resp, _, errs := xrequest.New().Post(fmt.Sprintf(apiAuthorizerToken, WxFunction.GetComponentAccessToken())).SendStruct(postData).EndStruct(result)
		if errs != nil {
			log.Error("获取公众号接口令牌失败", errs)
			return ""
		}
		err := unmarshalResponseToJson(resp, result)
		if err != nil {
			log.Error("公众号接口令牌解析失败", err)
			return ""
		}
		err = c.config.UpdataAccessToken(result.AccessToken, result.RefreshToken, result.ExpiresIn);
		if err != nil {
			log.Error("公众号接口令牌缓存失败: ", err)
			return ""
		}
	}
	return c.config.AccessToken()
}

//获取（刷新）授权公众号或小程序的接口调用凭据（令牌）
func (c *wechatConfig) RefreshAccessToken(forceRefresh bool) string {
	refreshLock := c.config.AccessTokenLock()
	refreshLock.Lock()
	defer refreshLock.Unlock()
	if forceRefresh {
		c.config.ExpireAccessToken();
	}
	return WxFunction.GetAuthAccessToken()
}

type ApiGetAuthorizerInfo struct {
	AuthorizerInfo struct {
		NickName string `json:"nick_name"`
		HeadImg  string `json:"head_img"`
		ServiceTypeInfo struct {
			Id int64 `json:"id"`
		}
		VerifyTypeInfo struct {
			Id int64 `json:"id"`
		}
		UserName string `json:"user_name"`
		Alias    string `json:"alias"`
	} `json:"authorizer_info"`
	QR string `json:"qrcode_url"`
	AuthorizationInfo struct {
		AppId string `json:"appid"`
		FuncInfo []struct {
			Funcscope struct {
				Id int64 `json:"id"`
			} `json:"funcscope_category"`
		} `json:"func_info"`
	} `json:"authorization_info"`
}

//获取授权方的帐号基本信息
func (c *wechatConfig) GetAuthorizerInfo() (*ApiGetAuthorizerInfo, []error) {
	postData := struct {
		ComponentAppid  string `json:"component_appid"`
		AuthorizerAppid string `json:"authorizer_appid"`
	}{
		ComponentAppid:  c.config.ComponentAppId(),
		AuthorizerAppid: c.config.AppId(),
	}

	result := &ApiGetAuthorizerInfo{}

	resp, _, errs := xrequest.New().Post(fmt.Sprintf(apiGetAuthorizerInfo, WxFunction.GetComponentAccessToken())).SendStruct(postData).EndStruct(result)
	if errs != nil {
		return nil, errs
	}
	err := unmarshalResponseToJson(resp, result)
	if err != nil {
		return nil, []error{err}
	}
	return result, nil
}

type ApiGetAuthorizerOption struct {
	AppId       string `json:"authorizer_appid"`
	OptionName  string `json:"option_name"`
	OptionValue string `json:"option_value"`
}

//获取授权方的选项设置信息
func (c *wechatConfig) GetAuthorizerOption(option string) (*ApiGetAuthorizerOption, []error) {
	postData := struct {
		ComponentAppid  string `json:"component_appid"`
		AuthorizerAppid string `json:"authorizer_appid"`
		OptionName      string `json:"option_name"`
	}{
		ComponentAppid:  c.config.ComponentAppId(),
		AuthorizerAppid: c.config.AppId(),
		OptionName:      option,
	}
	result := &ApiGetAuthorizerOption{}

	resp, _, errs := xrequest.New().Post(fmt.Sprintf(apiGetAuthorizerOption, WxFunction.GetComponentAccessToken())).SendStruct(postData).EndStruct(result)
	if errs != nil {
		log.Error("获取授权方选项设置信息失败", errs)
		return nil, errs
	}
	err := unmarshalResponseToJson(resp, result)
	if err != nil {
		log.Error("授权方选项设置信息解析失败", err)
		return nil, []error{err}
	}
	return result, nil
}

//设置授权方的选项信息
func (c *wechatConfig) SetAuthorizerOption(optionName, optionValue string) []error {
	postData := struct {
		ComponentAppid  string `json:"component_appid"`
		AuthorizerAppid string `json:"authorizer_appid"`
		OptionName      string `json:"option_name"`
		OptionValue     string `json:"option_value"`
	}{
		ComponentAppid:  c.config.ComponentAppId(),
		AuthorizerAppid: c.config.AppId(),
		OptionName:      optionName,
		OptionValue:     optionValue,
	}
	result := &ApiError{}

	resp, _, errs := xrequest.New().Post(fmt.Sprintf(apiSetAuthorizerOption, WxFunction.GetComponentAccessToken())).SendStruct(postData).EndStruct(result)
	if errs != nil {
		log.Error("设置授权方的选项信息失败", errs)
		return errs
	}
	return []error{unmarshalResponseToJson(resp, result)}
}

func (c *wechatConfig) CreateComponentLoginPageUrl(redirectUrl, preAuthCode string) string {
	u := url.QueryEscape(redirectUrl)
	return fmt.Sprintf(oauthUrl, c.config.ComponentAppId(), preAuthCode, u)
}
