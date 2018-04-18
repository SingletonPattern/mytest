package main

import (
	 "github.com/jmoiron/sqlx"
	 _ "github.com/go-sql-driver/mysql"
	 "goquestion/common"
	"fmt"
	 "goquestion/router"
	 //"github.com/go-redis/redis"
	"time"
	"reflect"
)


func init() {
	var err error
	common.Db,err = sqlx.Connect("mysql", "common:GuoanjiaCommon1@tcp(rm-2ze2y2k3554s14j47o.mysql.rds.aliyuncs.com)/guoanjia-common?charset=utf8&parseTime=True&loc=Local")
	//common.Db,err = sqlx.Connect("mysql", "root:Guoan2015@tcp(111.207.11.206:3306)/subject?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	common.Db.SetMaxIdleConns(10)
	common.Db.SetMaxOpenConns(100)
	common.Db.Ping()
	//redis init
	//common.Rds = redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})

	//pong, err := common.Rds.Ping().Result()
	//fmt.Println(pong, err)
}

type St struct {
	dt time.Time
	str string
}

func main(){
	var st St
	rt:= reflect.TypeOf(st.dt)
	//rk := rt.Kind()
	fmt.Println(rt.Name())
	router.Router()
}

