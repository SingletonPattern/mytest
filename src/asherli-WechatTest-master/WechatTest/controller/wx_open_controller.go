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
	"encoding/base64"
	"encoding/xml"
	"net/http"
	"strings"

	"weixin/initconfig"
	"weixin/qiniu"
	"weixin/utils"
	"weixin/wechat"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

var (
	ContextPath = "/openweixin"
)

type XmlAuthorizeMessage struct {
	AppId                 string `xml:"AppId"`
	CreateTime            string `xml:"CreateTime"`
	InfoType              string `xml:"InfoType"`
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
}

type XmlEncryptMessage struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}

//出于安全考虑，在第三方平台创建审核通过后，微信服务器每隔10分钟会向第三方的消息接收地址推送一次component_verify_ticket，用于获取第三方平台接口调用凭据
func ReceiveTicket(c echo.Context) error {
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	msgSignature := c.QueryParam("msg_signature")
	encryptType := c.QueryParam("encrypt_type")

	//读取Body数据，注意request中Body只能读取一次
	buf := make([]byte, 1024)
	n, _ := c.Request().Body.Read(buf)
	bytBody := buf[0:n]
	requestBody := string(bytBody)

	log.Infof("接收微信请求：[signature=[%s], encType=[%s], msgSignature=[%s], timestamp=[%s], nonce=[%s], requestBody=[%s] ", signature, encryptType, msgSignature, timestamp, nonce, requestBody)
	//验证必要参数是否存在
	if utils.IsAnyEmptyStr([]string{signature, timestamp, nonce}) {
		log.Errorf("微信请求必要参数缺失")
		return nil
	}
	encryptMsg := &XmlEncryptMessage{}
	//绑定Body数据
	err := xml.Unmarshal(bytBody, encryptMsg)
	if err != nil {
		return nil
	}
	//判断消息加密类型，开放平台所有消息均为加密类型
	if !strings.EqualFold("aes", encryptType) || !wechat.CheckSignature(wechat.WxConfig.ComponentToken(), timestamp, nonce, encryptMsg.Encrypt, msgSignature) {
		log.Errorf("非法请求，可能属于伪造的请求！")
		return nil
	}
	//构造AesKey
	aesKey, _ := base64.StdEncoding.DecodeString(wechat.WxConfig.ComponentEncodingAesKey() + "=")
	//Wechat消息数据解密
	src, err := wechat.DecryptMsg(encryptMsg.Encrypt, aesKey, wechat.WxConfig.ComponentAppId())
	if err != nil {
		return nil
	}
	authorizeMsg := &XmlAuthorizeMessage{}
	//绑定解密数据
	err = xml.Unmarshal(src, authorizeMsg)
	if err != nil {
		return nil
	}
	//更新redis缓存VerifyTicket
	err = initconfig.WxConfig.UpdateComponentVerifyTicket(authorizeMsg.ComponentVerifyTicket)
	if err != nil {
		return nil
	}
	return c.String(http.StatusOK, "success")
}

//消息回调接口
func Callback(c echo.Context) error {
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	msgSignature := c.QueryParam("msg_signature")
	encryptType := c.QueryParam("encrypt_type")

	//读取Body数据，注意request中Body只能读取一次
	buf := make([]byte, 1024)
	n, _ := c.Request().Body.Read(buf)
	bytBody := buf[0:n]
	requestBody := string(bytBody)

	log.Infof("接收微信请求：[signature=[%s], encType=[%s], msgSignature=[%s], timestamp=[%s], nonce=[%s], requestBody=[%s] ", signature, encryptType, msgSignature, timestamp, nonce, requestBody)
	//验证必要参数是否存在
	if utils.IsAnyEmptyStr([]string{signature, timestamp, nonce}) {
		log.Errorf("微信请求必要参数缺失")
		return nil
	}
	encryptMsg := &XmlEncryptMessage{}
	//绑定Body数据
	err := xml.Unmarshal(bytBody, encryptMsg)
	if err != nil {
		return nil
	}
	//判断消息加密类型，开放平台所有消息均为加密类型
	if !strings.EqualFold("aes", encryptType) || !wechat.CheckSignature(wechat.WxConfig.ComponentToken(), timestamp, nonce, encryptMsg.Encrypt, msgSignature) {
		log.Errorf("非法请求，可能属于伪造的请求！")
		return nil
	}
	//构造AesKey
	aesKey, _ := base64.StdEncoding.DecodeString(wechat.WxConfig.ComponentEncodingAesKey() + "=")
	//Wechat消息数据解密
	src, err := wechat.DecryptMsg(encryptMsg.Encrypt, aesKey, wechat.WxConfig.ComponentAppId())
	if err != nil {
		return nil
	}
	messgae := &wechat.Message{}
	//绑定解密数据
	err = xml.Unmarshal(src, messgae)
	if err != nil {
		return nil
	}
	//Handler演示
	//EventUnsubscribeHandler = func(m *EventSubscribe) ReplyMsg {
	//	log.Debugf("%+v", m)
	//
	//	// echo message
	//	ret := &ReplyText{
	//		ToUserName:   m.FromUserName,
	//		FromUserName: m.ToUserName,
	//		CreateTime:   m.CreateTime,
	//		Content:      fmt.Sprintf("Event=%s, EventKey=%s, Ticket=%s", m.Event, m.EventKey, m.Ticket),
	//	}
	//	ret.SetMsgType(MsgTypeText)
	//
	//	log.Debugf("replay message: %+v", ret)
	//	return ret
	//
	//消息处理器，当前版本未实现Session及消息路由 ,当前未初始化任何有效消息处理器 默认返回nil
	responseMessage := wechat.HandleMessage(messgae)
	if nil == responseMessage {
		return nil
	}
	return c.XML(http.StatusOK, responseMessage)
}

func AuthLoginPage(c echo.Context) error {
	r := c.Request();
	redirectUrl := c.Scheme() + "://" + r.Host + ContextPath + "/service/authorCallback"
	log.Infof("授权页面回调地址:%s", redirectUrl)
	preAuthCode, errs := wechat.WxFunction.GetPreAuthCode()
	if len(errs) > 0 {
		return c.String(http.StatusOK, errs[0].Error())
	}
	redirectUrl = wechat.WxFunction.CreateComponentLoginPageUrl(redirectUrl, preAuthCode)
	return c.Redirect(http.StatusFound, redirectUrl)
}

func AuthorCallback(c echo.Context) error {
	authCode := c.QueryParam("auth_code")
	expiresIn := c.QueryParam("expires_in")
	if utils.IsAnyEmptyStr([]string{authCode, expiresIn}) {
		return c.String(http.StatusOK, "缺失必要参数")
	}
	_, errs := wechat.WxFunction.GetAuthorization(authCode)
	if len(errs) > 0 {
		return c.String(http.StatusOK, errs[0].Error())
	}
	return c.String(http.StatusOK, "公众号授权成功")
}

type Result struct {
	Success bool        `json:"success"`
	Code    int64       `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func GetMediaFile(c echo.Context) error {
	mediaId := c.QueryParam("mediaId")
	file, err := wechat.WxFunction.MaterialImageOrVoiceDownload(mediaId)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusOK, &Result{
			Success: false,
			Code:    201,
			Msg:     "获取微信临时素材失败",
		})
	}

	newFileName, err := qiniu.Upload(file)

	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusOK, &Result{
			Success: false,
			Code:    201,
			Msg:     "上传七牛存储失败",
		})
	}

	return c.JSON(http.StatusOK, &Result{
		Success: true,
		Code:    200,
		Msg:     "操作成功",
		Data:    newFileName,
	})

}
