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
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"

	"weixin/utils"
	"weixin/utils/xrequest"

	"github.com/labstack/gommon/log"
)

const (
	jssdkJsapiTicket        = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	oauth2Authorization     = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&component_appid=%s#wechat_redirect"
	oauth2AccessToken       = "https://api.weixin.qq.com/sns/oauth2/component/access_token?appid=%s&code=%s&grant_type=authorization_code&component_appid=%s&component_access_token=%s"
	oauth2UserInfo          = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=%s"
	materialGetTemporaryURL = "https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
	materialBatchGetNewsURL = "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=%s"
)

type JsApiTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

func (c *wechatConfig) GetJsapiTicket() string {
	jsApiTicketLock := c.config.JsApiTicketLock()
	jsApiTicketLock.Lock()
	defer jsApiTicketLock.Unlock()
	if c.config.IsJsApiTicketExpired() {
		result := &JsApiTicket{}
		resp, _, errs := xrequest.New().Get(fmt.Sprintf(jssdkJsapiTicket, WxFunction.GetAuthAccessToken())).EndStruct(result)
		if errs != nil {
			log.Error("获取JsApiTicket失败: ", errs)
			return ""
		}
		err := unmarshalResponseToJson(resp, result)
		if err != nil {
			log.Error("JsApiTicket数据解析失败: ", err)
			return ""
		}
		err = c.config.UpdateJsApiTicket(result.Ticket, result.ExpiresIn);
		if err != nil {
			log.Error("JsApiTicket缓存失败: ", err)
			return ""
		}
	}
	return c.config.JsApiTicket()
}

func (c *wechatConfig) RefreshJsapiTicket(forceRefresh bool) string {
	refreshLock := c.config.JsApiTicketLock()
	refreshLock.Lock()
	defer refreshLock.Unlock()
	if forceRefresh {
		c.config.ExpireJsApiTicket();
	}
	return WxFunction.GetJsapiTicket()
}

type JsApiSignature struct {
	Appid     string `json:"appid"`
	NonceStr  string `json:"noncestr"`
	Timestamp int64  `json:"timestamp"`
	Url       string `json:"url"`
	Signature string `json:"signature"`
}

func (c *wechatConfig) CreateJsapiSignature(uri string) *JsApiSignature {
	timestamp := GetCurrTs()
	nonceStr := RandomStr(16)
	jsapiTicket := WxFunction.GetJsapiTicket()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", jsapiTicket, nonceStr, timestamp, uri)
	sigStr := Sha1Hex(str)
	return &JsApiSignature{
		Appid:     c.config.AppId(),
		NonceStr:  nonceStr,
		Timestamp: timestamp,
		Url:       uri,
		Signature: sigStr,
	}
}

func (c *wechatConfig) BuildOauth2AuthorizationUrl(redirectUrl, scope, state string) string {
	return fmt.Sprintf(oauth2Authorization, c.config.AppId(), url.QueryEscape(redirectUrl), scope, state, c.config.ComponentAppId())
}

type OAuth2AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"union_id"`
}

func (c *wechatConfig) GetOAuth2AccessToken(code string) (*OAuth2AccessToken, []error) {
	result := &OAuth2AccessToken{}
	resp, _, errs := xrequest.New().Get(fmt.Sprintf(oauth2AccessToken, c.config.AppId(), code, c.config.ComponentAppId(), WxFunction.GetComponentAccessToken())).EndStruct(result)
	if errs != nil {
		log.Error("获取用户授权令牌失败: ", errs)
		return nil, errs
	}
	err := unmarshalResponseToJson(resp, result)
	if err != nil {
		log.Error("用户授权令牌数据解析失败: ", err)
		return nil, []error{err}
	}
	return result, nil
}

//TODO 暂不实现该功能因尚不明确用户AccessToken是否固定通用
//func (c *wechatConfig) RefreshOAuth2AccessToken(code string) (*OAuth2AccessToken, []error) {
//	result := &OAuth2AccessToken{}
//	resp, _, errs := xrequest.New().Get(fmt.Sprintf(oauth2AccessToken, c.config.AppId(), code, c.config.ComponentAppId(), function.GetComponentAccessToken())).EndStruct(&result)
//	if errs != nil {
//		log.Error("获取用户授权令牌失败: ", errs)
//		return nil, errs
//	}
//	err := unmarshalResponseToJson(resp, result)
//	if err != nil {
//		log.Error("用户授权令牌数据解析失败: ", err)
//		return nil, []error{err}
//	}
//	return result, nil
//}

// Lang 国家地区语言版本
type Lang string

// 微信支持的语言
const (
	LangZHCN = "zh_CN" // 简体
	LangZHTW = "zh_TW" // 繁体
	LangEN   = "en"    // 英语
)

type UserInfo struct {
	Subscribe     int    `json:"subscribe"`
	OpenId        string `json:"openid"`
	NickName      string `json:"nickname"`
	Sex           int    `json:"sex"`
	Language      Lang   `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	HeadImgURL    string `json:"headimgurl"`
	SubscribeTime int    `json:"subscribe_time"`
	UnionId       string `json:"unionid"`
	Remark        string `json:"remark"`
	GroupId       int    `json:"groupid"`
}

func (c *wechatConfig) OAuth2GwtUserInfo(oAuth2AccessToken *OAuth2AccessToken, lang string) (*UserInfo, []error) {
	result := &UserInfo{}
	if utils.IsEmptyStr(lang) {
		lang = LangZHCN
	}
	resp, _, errs := xrequest.New().Get(fmt.Sprintf(oauth2UserInfo, oAuth2AccessToken.AccessToken, oAuth2AccessToken.Openid, lang)).EndStruct(result)
	if errs != nil {
		log.Error("获取用户失败: ", errs)
		return nil, errs
	}
	err := unmarshalResponseToJson(resp, result)
	if err != nil {
		log.Error("用户数据解析失败: ", err)
		return nil, []error{err}
	}
	return result, nil
}

func (c *wechatConfig) MaterialImageOrVoiceDownload(mediaId string) (*os.File, error) {
	response, _, errs := xrequest.New().Get(fmt.Sprintf(materialGetTemporaryURL, WxFunction.GetAuthAccessToken(), mediaId)).End()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	//contentTypeHeader := response.Header.Get("Content-Type")
	//if contentTypeHeader == "text/plain" {
	//	buf := make([]byte, 1024)
	//	response.Body.Read(buf)
	//	err := &ApiError{}
	//	return nil, json.Unmarshal(buf, err)
	//}
	//
	//var err error
	//params := make(map[string]string)
	//fileName := ""
	//if cd := response.Header.Get("Content-Disposition"); cd == "" {
	//	return nil, errors.New("missing Content-Disposition header")
	//} else if _, params, err = mime.ParseMediaType(cd); err != nil {
	//	return nil, fmt.Errorf("parse Content-Disposition header fail: %s", err.Error())
	//} else if fileName = params["filename"]; fileName == "" {
	//	return nil, errors.New("no filename in Content-Disposition header")
	//}


	contentTypeHeader := response.Header.Get("Content-Type")
	if contentTypeHeader == "text/plain" {
		buf := make([]byte, 1024)
		response.Body.Read(buf)
		err := &ApiError{}
		return nil, json.Unmarshal(buf, err)
	}
	contentDispositionHeader := response.Header.Get("Content-disposition")
	if utils.IsEmptyStr(contentDispositionHeader) {
		return nil, &ApiError{ErrMsg: "无法获取文件名"}
	}
	r, _ := regexp.Compile(".*filename=\"(.*)\"")
	fileName := r.FindStringSubmatch(contentDispositionHeader)[1]
	if utils.IsEmptyStr(fileName) {
		return nil, &ApiError{ErrMsg: "无法获取文件名"}
	}

	f, err := os.OpenFile(os.TempDir()+"/"+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	f.Seek(stat.Size(), 0)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	io.Copy(f, response.Body)
	return f, nil

}

// Article 永久图文素材
type Article struct {
	Title            string `json:"title"`              // 标题
	ThumbMediaId     string `json:"thumb_media_id"`     // 图文消息的封面图片素材id（必须是永久mediaID）
	Author           string `json:"author"`             // 作者
	Digest           string `json:"digest"`             // 图文消息的摘要，仅有单图文消息才有摘要，多图文此处为空
	ShowCoverPic     int    `json:"show_cover_pic"`     // 是否显示封面，0为false，即不显示，1为true，即显示
	Content          string `json:"content"`            // 图文消息的具体内容，支持HTML标签，必须少于2万字符，小于1M，且此处会去除JS
	URL              string `json:"url"`                // 图文页的URL
	ContentSourceURL string `json:"content_source_url"` // 图文消息的原文地址，即点击“阅读原文”后的URL
}

// NewsList 素材列表
type NewsList struct {
	TotalCount int64 `json:"total_count"` // 该类型的素材的总数
	ItemCount  int64 `json:"item_count"`  // 本次调用获取的素材的数量
	Item []struct {
		MediaId    string `json:"media_id"`
		UpdateTime int64 `json:"update_time"` // 这篇图文消息素材的最后更新时间
		Name       string `json:"name"`        // 文件名称
		URL        string `json:"url"`         // 图文页的URL，或者，当获取的列表是图片素材列表时，该字段是图片的URL
		Content struct {
			NewsItem []Article `json:"news_item"`
		} `json:"content"`                     // 本次调用获取的素材的数量
	} `json:"item"`                        // 多图文消息会在此处有多篇文章
}

type BatchGetNewsRequest struct {
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
	Count  int64  `json:"count"`
}

func (c *wechatConfig) MaterialFileBatchGet(mediaType string, offset int64, count int64) (*NewsList, []error) {
	result := &NewsList{}
	postData := &BatchGetNewsRequest{
		Type:   mediaType,
		Offset: offset,
		Count:  count,
	}
	resp, _, errs := xrequest.New().Post(fmt.Sprintf(materialBatchGetNewsURL, WxFunction.GetAuthAccessToken())).SendStruct(postData).EndStruct(result)
	if len(errs) > 0 {
		return nil, errs
	}

	err := unmarshalResponseToJson(resp, result)
	if err != nil {
		log.Error("素材列表解析失败: ", err)
		return nil, []error{err}
	}
	return result, nil
}
