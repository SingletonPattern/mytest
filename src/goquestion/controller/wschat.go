package controller

import (
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"fmt"
	"log"
	"os"
	"time"
	"net/http"
	"goquestion/service"
	"goquestion/common"
	"math/rand"
	"strconv"
	"encoding/json"
)

var (
	pwd, _         = os.Getwd()
	JSON           = websocket.JSON              // codec for JSON
	Message        = websocket.Message           // codec for string, []byte
	ActiveClients  = make(map[ClientConn]string) // map containing clients
	User           = make(map[string]string)
	Persons        = []service.Person{}  //所有未答题人员的数组
	AlreadyPersons = []service.Person{}  //已答题人员的数组
	Subjects       = []service.Subject{} //所有未答题
	SubjectT       = []service.Subject{} //所有未答题
	AlreadySubject = []service.Subject{} //已答题
)
// Initialize handlers and websocket handlers
func init() {
	User["aaa"] = "aaa"
	User["bbb"] = "bbb"
	User["test"] = "test"
	User["test2"] = "test2"
	User["test3"] = "test3"
}

//生成随机数
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}
func Remove(slice []service.Person, elems ...service.Person) []service.Person {
	isInElems := make(map[service.Person]bool)
	for _, elem := range elems {
		isInElems[elem] = true
	}
	w := 0
	for _, elem := range slice {
		if !isInElems[elem] {
			slice[w] = elem
			w += 1
		}
	}
	return slice[:w]
}
func RemoveSub(slice []service.Subject, elems ...service.Subject) []service.Subject {
	isInElems := make(map[service.Subject]bool)
	for _, elem := range elems {
		isInElems[elem] = true
	}
	w := 0
	for _, elem := range slice {
		if !isInElems[elem] {
			slice[w] = elem
			w += 1
		}
	}
	return slice[:w]
}
func SelectSubjectAndAnswer(c echo.Context) error {
	orderStr := c.QueryParam("order")
	order, i2 := strconv.Atoi(orderStr)
	if i2 != nil {
		c.JSON(http.StatusCreated, common.ResultMsg(201, "失败", nil))
	}
	if order <= 0 {
		order = 1
	}
	resMap := make(map[string]interface{})
	//查询总数
	total, e := service.TotalSubject()
	t, i2 := strconv.Atoi(total)
	if i2 != nil {
		c.JSON(http.StatusCreated, common.ResultMsg(201, "失败", nil))
	}
	if t < order {
		c.JSON(http.StatusCreated, common.ResultMsg(202, "已经没有题目", nil))
	}
	if e != nil {
		c.JSON(http.StatusCreated, common.ResultMsg(201, "失败", nil))
	}
	subjects, i := service.SelectSubject(order - 1)
	if ( i != nil) {
		c.JSON(http.StatusCreated, common.ResultMsg(201, "失败", nil))
	}
	resMap["totalCount"] = total
	subject := subjects[0]
	//name := ""
	subject.Name = ""
	resMap["subject"] = subject
	resMap["order"] = order + 1
	return c.JSON(http.StatusCreated, common.ResultMsg(200, "获取成功", resMap))
}

func EndRand(c echo.Context) error{
	for cs, na := range ActiveClients {
		if na != "" {
			msg := common.ResultMsg(10000, "停止滚动", nil)
			marshal, i3 := json.Marshal(msg)
			if i3 != nil {
			}
			s := string(marshal)
			if err := Message.Send(cs.websocket, s); err != nil {
				log.Println("Could not send message to ", cs.clientIP, err.Error())
			}
		}
	}
	return c.JSON(http.StatusOK, common.ResultMsg(200, "成功", nil))
}

//手机端选择题目
func SelectSubject(c echo.Context) error {
	SubjectT     = []service.Subject{}
	Persons        = []service.Person{}
	Subjects       = []service.Subject{}
	nameList, b := service.GetNameList()
	subjectList, b := service.GetSubjectAnswerList()
	//当有指定人时
	zhidingSubjectList  := []service.Subject{}
	for i:=0 ; i < len(subjectList) ;i ++{
		if !IsEmpty(subjectList[i].AnswerName){
			zhidingSubjectList = append(zhidingSubjectList, subjectList[i])
		}
	}
	if len(zhidingSubjectList) > 0{
		Persons = nameList
		Subjects = zhidingSubjectList
	}else{
		if b {
			fmt.Println(len(nameList),"获取所有人员列表")
			Persons = nameList
			Subjects = subjectList
			//service.Select()
		}
	}
	SubjectT = subjectList
	resMap := make(map[string]interface{})
	//subArr, _ := service.SelectSubject(order - 1)
	//查询总数
	total, _ := service.TotalSubject()
	t, _ := strconv.Atoi(total)
	subjects, _ := service.GetSubjectAnswerList() //题目数组
	resMap["alreadyTotal"] = len(AlreadySubject)
	resMap["totalCount"] = total
	resMap["subjectList"] = subjects
	//resMap["subject"] = subArr[0]
	nameList, b1 := service.GetNameList() //人名数组
	if !b1 {
		c.JSON(http.StatusCreated, common.ResultMsg(201, "获取名称列表失败", nil))
	}
	resMap["nameList"] = nameList

	for cs, na := range ActiveClients {
		if na != "" {
			//timestr:= time.Now().Format("2006-01-02 15:04:05")
			m := make(map[string]interface{})
			//m["Question"] = subj.Question
			m["total"] = t
			m["subjectList"] = subjects
			m["nameList"] = nameList
			m["alreadyTotal"] = len(AlreadySubject)
			//m["subject"] = subArr[0]
			msg := common.ResultMsg(10003, "返回题目列表和人员列表", m)
			marshal, i3 := json.Marshal(msg)
			if i3 != nil {
			}
			s := string(marshal)
			if err := Message.Send(cs.websocket, s); err != nil {
				log.Println("Could not send message to ", cs.clientIP, err.Error())
			}
		}
	}
	return c.JSON(http.StatusOK, common.ResultMsg(200, "成功", resMap))
}


func GetAllSub(c echo.Context) error {
	nameList, _ := service.GetNameList()
	subjectList, _ := service.GetSubjectAnswerList()
	resMap := make(map[string]interface{})
	resMap["nameList"] = nameList
	resMap["subjectList"] = subjectList
	return c.JSON(http.StatusOK, common.ResultMsg(200, "推送成功", resMap))
}
//手机端选择人员
func SelectPersion(c echo.Context) error {
	//subjectId := c.QueryParam("subjectId")
	//order := c.QueryParam("order")
	var person = service.Person{}
	if len(Subjects) <= 0 {
		c.JSON(http.StatusOK, common.ResultMsg(200, "题目已经回答完毕,请刷新题目重新开始", person))
	}
	intn := RandInt(1, len(Persons)) //获取随机数
	person = Persons[intn-1]
	Persons = Remove(Persons, person)
	AlreadyPersons = append(AlreadyPersons, person) //将答题人员放到已答题人员的数组中

	intns := RandInt(1, len(Subjects)) //获取随机数
	subject := Subjects[intns-1]
	Subjects = RemoveSub(Subjects, subject)
	SubjectT = RemoveSub(SubjectT, subject)
	AlreadySubject = append(AlreadySubject, subject)
	//nameList, b := service.GetNameList()
	//subjectList, b1 := service.GetSubjectAnswerList()
	//timestr:= time.Now().Format("2006-01-02 15:04:05")
	sub, i := service.GetSubjectAnswerById(subject.Id)
	anwerName := sub.AnswerName
	if !i {
	}
	//if !b || !b1{
	//	//查询出现错误
	//}
	//查询总数
	total, e := service.TotalSubject()
	t, i2 := strconv.Atoi(total)
	if e!= nil || i2 != nil {

	}
	for cs, na := range ActiveClients {
		if na != "" {

			m := make(map[string]interface{})

			if !IsEmpty(anwerName){
				m["Name"] = anwerName
			}else{
				m["Name"] = person.Name
			}
			m["Question"] = subject.Question
			m["nameList"] = Persons
			m["subjectList"] = SubjectT
			m["answer"] =  subject.Answer
			m["total"] = t
			m["alreadyTotal"] = len(AlreadySubject)
			msg := common.ResultMsg(10001, "推送问题和答题人员名称", m)
			marshal, i3 := json.Marshal(msg)
			if i3 != nil {
			}
			s := string(marshal)
			if err := Message.Send(cs.websocket, s); err != nil {
				log.Println("Could not send message to ", cs.clientIP, err.Error())
			}
		}
	}
	if len(Persons) <= 0 {
		Persons = AlreadyPersons
		AlreadyPersons = []service.Person{}
	}
	resMap := make(map[string]interface{})
	if !IsEmpty(anwerName){
		resMap["person"] = anwerName
		peoples, i2 := service.GetNameById(anwerName)
		if len(peoples) > 0{
			person = peoples[0]
		}
		if !i2 {
		}
	}else{
		resMap["person"] = person.Name
	}
	resMap["subject"] = subject
	resMap["nameList"] = Persons
	resMap["subjectList"] = SubjectT
	resMap["alreadyTotal"] = len(AlreadySubject)
	resMap["total"] = t
	//返回手机遥控器
	personandswer := new(service.PersonAnswer)
	personandswer.Id = UUID()
	personandswer.SubjectId = subject.Id
	personandswer.PersonId = person.Name
	personandswer.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	service.DelPresonAndSubject(*personandswer)
	service.SavePresonAndSubject(*personandswer)
	return c.JSON(http.StatusOK, common.ResultMsg(200, "推送成功", resMap))
}

// Client connection consists of the websocket and the client ip
type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

func Hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		var err error
		defer ws.Close()

		client := ws.Request().RemoteAddr
		log.Println("Client connected:", client)
		sockCli := ClientConn{ws, client}
		ActiveClients[sockCli] = "1"
		log.Println("Number of clients connected:", len(ActiveClients))

		for {
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				// If we cannot Read then the connection is closed
				log.Println("Websocket Disconnected waiting", err.Error())
				// remove the ws client conn from our active clients
				delete(ActiveClients, sockCli)
				log.Println("Number of clients still connected:", len(ActiveClients))
				return
			}
			ActiveClients[sockCli] = "a"
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
