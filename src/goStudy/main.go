package main

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}
func sayBayBay(w http.ResponseWriter,r *http.Request) {
	body ,err := ioutil.ReadAll(r.Body);
	if err != nil{

	}
	var data  = []byte("{\"code\": 1,\"msg\": \"failed\"}")
	w.Write(data);
	s := bytes.NewBuffer(body).String()
	fmt.Println("say bay bay  man !" , s);
}
func main(){
	//connectDB(); //连接数据库
	connectMysql(); //连接数据库
	http.HandleFunc("/", sayhelloName) //设置访问的路由s
	http.HandleFunc("/sayBayBay", sayBayBay) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func connectDB(){
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/csizg_card_manage?charset=utf8")//对应数据库的用户名和密码
	defer db.Close()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("success")
	}
	rows, err := db.Query("SELECT * ")
	if err != nil {
		panic(err)
		return
	}
	for rows.Next() {
		var name int
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		fmt.Println(name)
	}
}
func connectMysql(){
	dbw := DbWorker{
		Dsn: "root:@tcp(127.0.0.1:3306)/csizg_card_manage",
	}
	db, err := sql.Open("mysql",
		dbw.Dsn)
	defer db.Close();
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("启动成功");

}
type DbWorker struct {
	//mysql data source name
	Dsn string
}