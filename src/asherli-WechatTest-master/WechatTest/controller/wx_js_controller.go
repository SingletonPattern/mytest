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

	"weixin/utils"
	"weixin/wechat"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

/**
 * 获取JSSDK 参数
 * @param local_url 页面地址
 * @return JSSDK参数
 */
func GetJsapiSignature(c echo.Context) error {
	uri := c.QueryParam("local_url")
	if utils.IsEmptyStr(uri) {
		return nil
	}
	jsApiTicket := wechat.WxFunction.CreateJsapiSignature(uri)
	return c.JSON(http.StatusOK, jsApiTicket)
}

func GetJsapiSignatureForWexFive(c echo.Context) error {
	action:=c.QueryParam("action")
	if action=="getTicket" {
		ticket:=wechat.WxFunction.GetJsapiTicket()
		return c.String(http.StatusOK,ticket)
	}
	log.Error("已去除支付功能")
	return nil
}
