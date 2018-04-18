package main

import (
	_ "beeapi/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type myt struct {
	i int
	s string
}
var first int8 = 1
var myvalue *int8 = &first
func init() {
	orm.RegisterDataBase("default", "mysql", "root@tcp(172.16.4.12:3306)/test")
}

func main() {

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

