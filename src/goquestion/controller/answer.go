package controller

import (
	"github.com/labstack/echo"
	"goquestion/service"
	"github.com/snluu/uuid"
	"strings"
	"time"
	"net/http"
	"goquestion/common"
	"fmt"
)

func DelPersonAndswer(c echo.Context) (err error) {
	b := service.DelPersonAndswer()
	if !b {
		c.JSON(http.StatusOK, common.Result{201, "清除失败", nil})
	}
	//所有人员列表
	nameList, b := service.GetNameList()
	Persons = nameList
	AlreadyPersons = []service.Person{}
	subjectList, b := service.GetSubjectAnswerList()
	Subjects = subjectList
	AlreadySubject = []service.Subject{}

	//nameList, b := service.GetNameList()
	//subjectList, b := service.GetSubjectAnswerList()
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
	if !b {
		c.JSON(http.StatusOK, common.Result{201, "获取所有人员列表失败", nil})
	}
	return c.JSON(http.StatusOK, common.Result{200, "清除成功", nil})
}

//保存
func SavePresonAndSubject(c echo.Context) (err error) {
	u := new(service.PersonAnswer)
	if err = c.Bind(u); err != nil {
		return err
	}
	u.Id = UUID()
	u.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	b := service.SavePresonAndSubject(*u)
	if !b {
		c.JSON(http.StatusOK, common.Result{201, "添加失败", nil})
	}
	return c.JSON(http.StatusOK, common.Result{200, "添加成功", nil})
}

//保存人员表
func SavePerson(c echo.Context) (err error) {
	u := new(service.Person)
	if err = c.Bind(u); err != nil {
		return err
	}
	name := u.Name
	id := UUID()
	now := time.Now().Format("2006-01-02 15:04:05")
	b := service.SavePerson(name, id, now)
	if !b {
		c.JSON(http.StatusOK, common.Result{201, "添加失败", nil})
	}
	p := new(service.Person)
	p.Id = id
	p.Name = name
	Persons = append(Persons, *p)
	fmt.Println(len(Persons))
	return c.JSON(http.StatusOK, common.Result{200, "添加成功", nil})
}

//删除人员
func DelPerson(c echo.Context) (err error) {
	u := new(service.Person)
	if err = c.Bind(u); err != nil {
		return err
	}
	id := u.Id
	person, e := service.DelPerson(id)
	if e != nil {
		return c.JSON(http.StatusOK, common.Result{201, "删除失败", nil})
	}
	AlreadyPersons = Remove(AlreadyPersons, person)
	Persons = Remove(Persons, person)
	fmt.Println("已答题人员数组中删除", len(AlreadyPersons), "++++++ 未答题数组中删除：", len(Persons))
	return c.JSON(http.StatusOK, common.Result{200, "删除成功", nil})
}

//修改名称接口
func UpdatePerson(c echo.Context) (err error) {
	u := new(service.Person)
	if err = c.Bind(u); err != nil {
		return err
	}
	persons , _ :=service.GetNameListById(u.Id)
	for i := 0 ; i < len(Persons) ; i++  {
		if strings.EqualFold(Persons[i].Id ,u.Id) {
			Persons[i].Name = u.Name
		}
	}
	for i := 0 ; i < len(AlreadyPersons) ; i++  {
		if strings.EqualFold(AlreadyPersons[i].Id ,u.Id) {
			AlreadyPersons[i].Name = u.Name
		}
	}
	b := service.UpdatePerson(u.Name, u.Id,persons[0].CreateTime)
	if !b {
		return c.JSON(http.StatusOK, common.Result{201, "修改失败", nil})
	}
	return c.JSON(http.StatusOK, common.Result{200, "修改成功", nil})
}

//获取人员列表接口
func GetNameList(c echo.Context) (err error) {
	nameList, b := service.GetNameList()
	if b {
		return c.JSON(http.StatusOK, common.Result{200, "获取成功", nameList})
	}
	return c.JSON(http.StatusOK, common.Result{201, "获取失败", nil})
}

//保存题目接口
func SaveSubjectAnswer(c echo.Context) (err error) {
	u := new(service.Subject)
	if err = c.Bind(u); err != nil {
		return err
	}
	u.Id = uuid.Rand().Hex()
	u.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	b := service.SaveSubjectAnswer(u)
	Subjects = append(Subjects, *u)
	SubjectT = append(SubjectT, *u)
	fmt.Println(len(Subjects))
	if !b {
		return c.JSON(http.StatusOK, common.Result{201, "添加失败", nil})
	}
	return c.JSON(http.StatusOK, common.Result{200, "添加成功", nil})
}

//修改题目接口
func UpdateSubjectAnswer(c echo.Context) (err error) {
	u := new(service.Subject)
	if err = c.Bind(u); err != nil {
		return err
	}
	//u.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	b := service.UpdateSubjectAnswer(u)
	for i := 0 ; i < len(Subjects) ; i++  {
		if strings.EqualFold(Subjects[i].Id ,u.Id) {
			Subjects[i].AnswerName = u.AnswerName
			Subjects[i].Question = u.Question
			Subjects[i].Answer = u.Answer
		}
	}
	for i := 0 ; i < len(SubjectT) ; i++  {
		if strings.EqualFold(SubjectT[i].Id ,u.Id) {
			SubjectT[i].AnswerName = u.AnswerName
			SubjectT[i].Question = u.Question
			SubjectT[i].Answer = u.Answer
		}
	}
	for i := 0 ; i < len(AlreadySubject) ; i++  {
		if strings.EqualFold(AlreadySubject[i].Id ,u.Id) {
			AlreadySubject[i].AnswerName = u.AnswerName
			AlreadySubject[i].Question = u.Question
			AlreadySubject[i].Answer = u.Answer
		}
	}

	if !b {
		return c.JSON(http.StatusOK, common.Result{201, "修改失败", nil})
	}
	return c.JSON(http.StatusOK, common.Result{200, "修改成功", nil})
}

//删除题目接口
func DelSubjectAnswer(c echo.Context) (err error) {
	id := c.QueryParam("id")
	//idArr := strings.Split(ids, ",")
	var b = true
	subject, i := service.GetSubjectAnswerById(id)
	b = service.DelSubjectAnswer(id)
	if !i {

	}
	AlreadySubject = RemoveSub(AlreadySubject, subject)
	Subjects = RemoveSub(Subjects, subject)
	if !b {
		return c.JSON(http.StatusOK, common.Result{201, "删除失败", nil})
	}
	return c.JSON(http.StatusOK, common.Result{200, "删除成功", nil})
}

//题目列表接口
func GetSubjectAnswerList(c echo.Context) (err error) {
	subjectAnswerList, b := service.GetSubjectAnswerList()
	if !b {
		return c.JSON(http.StatusOK, common.Result{201, "获取失败", nil})
	}
	return c.JSON(http.StatusOK, common.Result{200, "获取成功", subjectAnswerList})
}

/**
判断字符串是否为空
 */
func IsEmpty(str string) bool {
	if str != "" && !strings.EqualFold(str, "") {
		return false
	}
	return true
}

func UUID() string {
	return uuid.Rand().Hex()
}
