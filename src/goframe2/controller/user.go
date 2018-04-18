package controller

import (
	"github.com/labstack/echo"
	"net/http"
	"goframe2/service"
)



func  GetAll(c echo.Context) error {

	return c.JSON(http.StatusCreated,service.Frist())
}
