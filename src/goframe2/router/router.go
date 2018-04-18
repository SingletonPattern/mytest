package router
import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"goframe2/controller"
)


func Router(){
	e:=echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})
	e.POST("/users",controller.GetAll)
	// Start server
	e.Logger.Fatal(e.Start(":8082"))
}

