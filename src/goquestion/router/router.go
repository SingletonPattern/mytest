package router
import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"goquestion/controller"
)


func Router(){
	e:=echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})
	e.POST("/users",controller.GetAll)

	e.GET("/save",controller.Save)
	e.GET("/ws", controller.Hello)
	e.GET("/query",controller.Query)
	e.GET("/find",controller.Find)
	e.GET("/select",controller.SelectPersion)  //点击停止随机时 将题目 答题人返回同事 同步到PC
	e.GET("/selectSubject",controller.SelectSubject)  //页面加载时调用接口
	e.GET("/selectSubjectAndAnswer",controller.SelectSubjectAndAnswer)   //手机端只显示题目答案接口 不推送消息
	e.GET("/endRand",controller.EndRand)  //取消按钮

	e.GET("/getAllSub",controller.GetAllSub)   //PC端点击开始按钮



	//保存人员接口
	e.POST("/savePerson",controller.SavePerson)
	//删除人员接口
	e.POST("/delPerson",controller.DelPerson)
	//修改人员名称接口
	e.POST("/updatePerson",controller.UpdatePerson)
	//获取列表接口
	e.GET("/getNameList",controller.GetNameList)
	//保存题目答案接口
	e.POST("/saveSubjectAnswer",controller.SaveSubjectAnswer)
	//修改题目答案接口
	e.POST("/updateSubjectAnswer",controller.UpdateSubjectAnswer)
	//获取题目答案人员接口
	e.GET("/getSubjectAnswerList",controller.GetSubjectAnswerList)
	//删除题目答案接口
	e.GET("/delSubjectAnswer",controller.DelSubjectAnswer)
	//保存用户答题信息
	e.POST("/savePresonAndSubject",controller.SavePresonAndSubject)
	//删除所有用户答题信息
	e.GET("/delPresonAndSubject",controller.DelPersonAndswer)


	//选择回答人员接口==》同时将通过websocket 将题目 和回答人 推送到 PC页面中

	//下一题 上一题 按钮接口
	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}

