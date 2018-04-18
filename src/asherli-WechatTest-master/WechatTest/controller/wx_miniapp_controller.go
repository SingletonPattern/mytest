package controller

import (
	"net/http"

	"weixin/wechat"

	"github.com/labstack/echo"
)

func MiniAppLogin(c echo.Context) error{
	jscode:=c.QueryParam("jsCode")
	miniAppSessionKey,errs:=wechat.WxFunction.MiniAppLogin(jscode)
	if len(errs) > 0 {
		return c.JSON(http.StatusOK,&Result{
			Success:false,
			Msg:errs[0].Error(),
			Code:10001,
			Data:nil,
		})
	}
	return c.JSON(http.StatusOK,&Result{
		Success:true,
		Msg:"成功",
		Code:0,
		Data:miniAppSessionKey.Openid,
	})
}