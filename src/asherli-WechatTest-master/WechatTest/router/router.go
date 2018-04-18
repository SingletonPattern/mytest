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

package router

import (
	"weixin/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

)

var (
	ContextPath  = "/openweixin"
)

func Router() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials:true,
		MaxAge:3600,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE,echo.OPTIONS},
	}))

	e.Any(ContextPath+"/service/event/authorize", func(context echo.Context) error {
		return controller.ReceiveTicket(context)
	})

	e.Any(ContextPath+"/service/:appid/callback", func(context echo.Context) error {
		return controller.Callback(context)
	})

	e.Any(ContextPath+"/service/authLoginPage", func(context echo.Context) error {
		return controller.AuthLoginPage(context)
	})

	e.Any(ContextPath+"/service/authorCallback", func(context echo.Context) error {
		return controller.AuthorCallback(context)
	})

	e.Any(ContextPath+"/jsapi/getJsapiSignature", func(context echo.Context) error {
		return controller.GetJsapiSignature(context)
	})

	e.Any(ContextPath+"/user/getCode", func(context echo.Context) error {
		return controller.GetCode(context)
	})

	e.Any(ContextPath+"/user/getUserInfo", func(context echo.Context) error {
		return controller.GetUserInfo(context)
	})

	e.Any(ContextPath+"/user/getUserOpenId", func(context echo.Context) error {
		return controller.GetUserOpenId(context)
	})

	e.Any(ContextPath+"/service/getMediaFile", func(context echo.Context) error {
		return controller.GetMediaFile(context)
	})

	e.Any(ContextPath+"/wexFive/wexFiveJsApi", func(context echo.Context) error {
		return controller.GetJsapiSignatureForWexFive(context)
	})
	e.Any(ContextPath+"/material/materialFileBatchGet", func(context echo.Context) error {
		return controller.MaterialFileBatchGet(context)
	})
	e.Any(ContextPath+"/miniapp/login", func(context echo.Context) error {
		return controller.MiniAppLogin(context)
	})

	e.Logger.Fatal(e.Start(":8083"))
}
