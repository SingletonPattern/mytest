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

package controller

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	. "weixin/initconfig"
	"weixin/model"
	"weixin/qiniu"
	"weixin/utils"
	"weixin/wechat"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func GetCode(c echo.Context) error {
	redirectUrl := c.QueryParam("redirect_url")
	scope := c.QueryParam("scope")
	state := c.QueryParam("state")
	if utils.IsEmptyStr(redirectUrl) {
		//TODO 此处应跳转错误页
		return nil
	}

	redirectUrla:= c.Scheme() + "://" + c.Request().Host + ContextPath+"/user/getUserInfo?redirect_url=" + redirectUrl
	if utils.IsEmptyStr(scope) || strings.EqualFold(scope, "snsapi_userinfo") {
		scope = "snsapi_userinfo"
	} else {
		scope = "snsapi_base"
		redirectUrla = c.Scheme() + "://" + c.Request().Host + ContextPath+"/user/getUserOpenId?redirect_url=" + redirectUrl
	}
	if utils.IsEmptyStr(state) {
		state = wechat.RandomStr(16)
	}
	redirectUrl = wechat.WxFunction.BuildOauth2AuthorizationUrl(redirectUrla, scope, state)
	return c.Redirect(http.StatusFound, redirectUrl)
}

func GetUserInfo(c echo.Context) error {
	redirectUrl := c.QueryParam("redirect_url")
	code := c.QueryParam("code")

	if strings.Index(redirectUrl, "?") >= 0 {
		redirectUrl = redirectUrl + "&"
	} else {
		redirectUrl = redirectUrl + "?"
	}

	if utils.IsAnyEmptyStr([]string{redirectUrl, code}) {
		return c.Redirect(http.StatusFound, redirectUrl+"error=true")
	}

	oAuth2AccessToken, errs := wechat.WxFunction.GetOAuth2AccessToken(code)
	if len(errs) > 0 {
		log.Error(errs[0])
		return c.Redirect(http.StatusFound, redirectUrl+"error=true")
	}

	user := &model.WxUser{
		UOpenid: oAuth2AccessToken.Openid,
	}
	exist, err := MySQL.Exist(user)
	if err == nil && exist {
		has, err := MySQL.Get(user)
		if has && err == nil {
			redirectUrl = redirectUrl + "openid=" + user.UOpenid +
				"&headimgurl=" + user.UHeadimg +
				"&nickname=" + url.QueryEscape(user.UNickname) +
				"&city=" + url.QueryEscape(user.UCity) +
				"&province=" + url.QueryEscape(user.UProvince) +
				"&country=" + url.QueryEscape(user.UCountry) +
				"&sex=" + url.QueryEscape(user.USex) +
				"&error=false"
		log.Error(redirectUrl)
			return c.Redirect(http.StatusFound, redirectUrl)
		}
	}

	userInfo, errs := wechat.WxFunction.OAuth2GwtUserInfo(oAuth2AccessToken, wechat.LangZHCN)
	if len(errs) > 0 {
		return c.Redirect(http.StatusFound, redirectUrl+"error=true")
	}

	sex := ""
	switch userInfo.Sex {
	case 1:
		sex = "男性"
	case 2:
		sex = "女性"
	default:
		sex = "未知"
	}

	headImg := "https://img.guoanfamily.com/"+qiniu.Fetch(userInfo.HeadImgURL, utils.CreateUUID()+".jpg")

	saveUser := &model.WxUser{
		Id:          utils.CreateUUID(),
		USex:        sex,
		UCountry:    userInfo.Country,
		UCity:       userInfo.City,
		UNickname:   userInfo.NickName,
		UHeadimg:    headImg,
		UOpenid:     userInfo.OpenId,
		UUnionid:    userInfo.UnionId,
		UProvince:   userInfo.Province,
		UCreateTime: time.Now(),
	}
	_, err = MySQL.Insert(saveUser)
	if err != nil {
		log.Error(err)
		return c.Redirect(http.StatusFound, redirectUrl+"error=true")
	}
	redirectUrl = redirectUrl + "openid=" + userInfo.OpenId +
		"&headimgurl=" + headImg +
		"&nickname=" + url.QueryEscape(userInfo.NickName) +
		"&city=" + url.QueryEscape(userInfo.City) +
		"&province=" + url.QueryEscape(userInfo.Province) +
		"&country=" + url.QueryEscape(userInfo.Country) +
		"&sex=" + url.QueryEscape(sex) +
		"&error=false"
	log.Error(redirectUrl)
	return c.Redirect(http.StatusFound, redirectUrl)

}

func GetUserOpenId(c echo.Context) error {
	redirectUrl := c.QueryParam("redirect_url")
	code := c.QueryParam("code")

	if utils.IsAnyEmptyStr([]string{redirectUrl, code}) {
		return c.Redirect(http.StatusFound, redirectUrl+"error=true")
	}
	if strings.Index(redirectUrl, "?") >= 0 {
		redirectUrl = redirectUrl + "&"
	} else {
		redirectUrl = redirectUrl + "?"
	}
	oAuth2AccessToken, errs := wechat.WxFunction.GetOAuth2AccessToken(code)
	if len(errs) > 0 {
		return c.Redirect(http.StatusFound, redirectUrl+"error=true")
	}

	redirectUrl = redirectUrl + "openid=" + oAuth2AccessToken.Openid + "&error=false"
	log.Error(redirectUrl)
	return c.Redirect(http.StatusFound, redirectUrl)

}
