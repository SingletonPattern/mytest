package controller

import (
	"net/http"

	"weixin/wechat"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type MaterialFileBatchRequest struct {
	MediaType string `json:"mediaType"`
	Offset    int64  `json:"offset"`
	Count     int64  `json:"count"`
}

func MaterialFileBatchGet(c echo.Context) error {
	request := &MaterialFileBatchRequest{}
	err := c.Bind(request)
	if err != nil {
		log.Error(err)
		return nil
	}
	list, errs := wechat.WxFunction.MaterialFileBatchGet(request.MediaType, request.Offset, request.Count)
	if len(errs) > 0 {
		log.Error(errs)
		return nil
	}
	return c.JSON(http.StatusOK, list)
}
