package service

import (
	"goquestion/common"
	"github.com/jmoiron/sqlx"
	"log"
	"fmt"
)

//人员表
type Person struct {
	Id   string
	Name string
	CreateTime string
}

//题目表
type Subject struct {
	Id         string
	Order      int
	Question   string
	Answer     string
	CreateTime string
	AnswerName string
}

//已回答记录表
type PersonAnswer struct {
	Id         string
	SubjectId  string
	PersonId   string
	CreateTime string
}

type ResSubject struct{
	Id string
	Question string
	Answer   string
	Name   string
}

func DelPersonAndswer() (bool) {
	tx := common.Db.MustBegin()
	Finally(tx)
	sql := "DELETE FROM person_answer "
	tx.MustExec(sql)
	fmt.Println(sql ,"清空所有答题人员表")
	tx.Commit()
	b := true
	return b
}
func DelPresonAndSubject(personAndswer PersonAnswer) (bool) {
	tx := common.Db.MustBegin()
	Finally(tx)
	sql := "DELETE FROM person_answer WHERE person_id =  ? AND subject_id = ?"
	tx.MustExec(sql,personAndswer.PersonId, personAndswer.SubjectId)
	fmt.Println(sql ,"参数：",personAndswer.PersonId,personAndswer.SubjectId,"根据题目ID和答题人ID删除记录")
	tx.Commit()
	b := true
	return b
}

func SavePresonAndSubject(personAndswer PersonAnswer) (bool) {
	tx := common.Db.MustBegin()
	Finally(tx)
	sql := "INSERT INTO person_answer (id,subject_id,person_id,create_time) VALUES (?,?,?,?)"
	tx.MustExec(sql, personAndswer.Id, personAndswer.SubjectId, personAndswer.PersonId, personAndswer.CreateTime)
	tx.Commit()
	fmt.Println(sql ,"参数：",personAndswer.Id, personAndswer.SubjectId, personAndswer.PersonId, personAndswer.CreateTime,"保存答题信息")
	b := true
	return b
}

//手机端选择要答的题
func SelectSubjectById(id string) ([]ResSubject, error) {
	var subjects []ResSubject
	sql := " SELECT" +
		" s.id AS id , s.question AS question ,s.answer AS answer ,per.`name` AS `name`" +
		" FROM `subject` s" +
		" LEFT JOIN person_answer p ON s.id = p.subject_id" +
		" LEFT JOIN person per ON per.id = p.person_id where s.id = ? GROUP BY s.id  ORDER BY  p.create_time DESC  "
	e := Query(&subjects, sql, id)
	//e := common.Db.Select(&subjects, sql, id)
	fmt.Println(sql ,"参数：",id,"根据题目ID查询答题信息")
	if e != nil {
		return nil, e
	}
	return subjects, nil
}
//手机端选择要答的题
func SelectSubject(order int) ([]ResSubject, error) {
	subjects := []ResSubject{}
	sql := " SELECT  t.id AS id , t.question AS question ,t.answer AS answer ,t.`name` AS `name`  FROM (SELECT   " +
		"  s.id AS id , s.question AS question ,s.answer AS answer ,per.`name` AS `name` " +
		" FROM `subject` s  LEFT JOIN person_answer p ON s.id = p.subject_id  " +
		" LEFT JOIN person per ON per.id = p.person_id  ORDER BY  p.create_time DESC )  t  GROUP BY  t.id   LIMIT ?,1"
	e := Query(&subjects, sql, order)
	//e := common.Db.Select(&subjects, sql, order)
	fmt.Println(sql ,"参数：",order,"根据题目序号查询答题信息")
	if e != nil {
		return nil, e
	}
	return subjects, nil
}
type Total struct {
	TotalCount string
}

//查询题目总数
func TotalSubject() (string, error) {
	totals := []string{}
	sql := "SELECT count(*) FROM subject"
	//querybySql := QuerybySql(sql)
	//for i := 0 ;i< len(querybySql) ; i++  {
	//	querybySql[0]
	//}
	//e := Query(&totals, sql)
	e := common.Db.Select(&totals, sql)
	if e != nil {
		//return nil, e

	}
	return totals[0], nil
}

//添加人员接口
func SavePerson(name, id ,createTime string) (b bool) {
	tx := common.Db.MustBegin()
	Finally(tx)
	sql := "insert INTO person (id,`name`,create_time) values (?,?,?)"
	tx.MustExec(sql, id, name,createTime )
	tx.Commit()
	fmt.Println("添加人员成功",sql, id, name,createTime )
	b = true
	return b
}

//删除人员接口
func DelPerson(id string) (p Person,e error) {
	tx := common.Db.MustBegin()
	persons := []Person{}
	selectSql := "select id AS  id , `name` AS  `name`,SUBSTR(create_time FROM 1 FOR 19)  AS createTime from  person WHERE id = ?"
	e = Query(&persons, selectSql, id)
	//e := common.Db.Select(&persons, selectSql, id)
	if e != nil {
		return p, e
	}
	Finally(tx)
	tx.MustExec("DELETE FROM  person WHERE id = ?", id)
	tx.Commit()
	if len(persons) > 0{
		p = persons[0]
	}
	fmt.Println("删除人员成功","DELETE FROM  person WHERE id = ?", id)
	return p, nil
}

//修改人员名称
func UpdatePerson(name, id,createTime string) (b bool) {
	tx := common.Db.MustBegin()
	Finally(tx)
	tx.MustExec("UPDATE person SET `name` = ? ,create_time = ? WHERE id = ?", name,createTime,id)
	tx.Commit()
	b = true
	fmt.Println("修改人员成功","UPDATE person SET `name` = ? WHERE id = ?",  name, id)
	return b
}

//获取答题人员列表
func GetNameList() ([]Person, bool) {
	persons := []Person{}
	sql := "SELECT id AS id , `name` AS `name`  FROM person ORDER BY create_time ASC "
	e := Query(&persons, sql)
	//e := common.Db.Select(&persons, sql)
	if e != nil {
		return nil, false
	}
	return persons, true
}
//根据答题人Id查询数据
func GetNameListById(id string) ([]Person, bool) {
	persons := []Person{}
	sql := "SELECT id AS id , `name` AS `name`  ,SUBSTR(create_time FROM 1 FOR 19)  AS createTime FROM person where id = ? ORDER BY create_time ASC "
	e := Query(&persons, sql,id)
	//e := common.Db.Select(&persons, sql)
	if e != nil {
		return nil, false
	}
	return persons, true
}
//根据name获取答题人姓名
func GetNameById(name string) ([]Person, bool) {
	persons := []Person{}
	sql := "SELECT id AS id , `name` AS `name`  ,create_time AS createTime FROM person where `name` = ? ORDER BY create_time ASC "
	e := Query(&persons, sql,name)
	//e := common.Db.Select(&persons, sql)
	if e != nil {
		return nil, false
	}
	return persons, true
}

//添加题目接口
func SaveSubjectAnswer(subject *Subject) (b bool) {
	tx := common.Db.MustBegin()
	Finally(tx)
	sql := "insert INTO subject (id,`order`,`question`,`answer`,create_time,answer_name) values(?,?,?,?,?,?)"
	tx.MustExec(sql, subject.Id, subject.Order, subject.Question, subject.Answer, subject.CreateTime,subject.AnswerName)
	tx.Commit()
	b = true
	return b
}

//修改题目接口
func UpdateSubjectAnswer(subject *Subject) (b bool) {
	tx := common.Db.MustBegin()
	Finally(tx)
	sql := "UPDATE `subject` SET `order` = ?,`question` = ?, `answer` = ?,answer_name = ? WHERE id = ?"
	tx.MustExec(sql, subject.Order, subject.Question, subject.Answer,subject.AnswerName, subject.Id)
	tx.Commit()
	b = true
	return b
}

//删除题目题目接口
func DelSubjectAnswer(id string) (b bool) {
	tx := common.Db.MustBegin()
	Finally(tx)
	sql := "DELETE FROM `subject` WHERE id = ?"
	tx.MustExec(sql, id)
	tx.Commit()
	b = true
	return b
}

//获取题目列表接口
func GetSubjectAnswerList() ([]Subject, bool) {
	subjects := []Subject{}
	sql := "SELECT id ,`order`,`question`,`answer` , answer_name AS answerName FROM `subject` order by create_time ASC"
	e := Query(&subjects, sql)
	//e := common.Db.Select(&subjects, sql)
	if e != nil {
		return nil, false
	}
	return subjects, true
}

//查询单个题目信息
func GetSubjectAnswerById(id string) (Subject, bool) {
	subjects := []Subject{}
	sql := "SELECT id ,`question`,`answer` , answer_name AS answerName FROM `subject` where id = ? "
	e := Query(&subjects, sql,id)
	if e != nil{

	}

	return subjects[0], true
}

//最终执行的语句  要放到mustbegin后第一行
func Finally(tx *sqlx.Tx) {
	defer func() {
		if r := recover(); r != nil {
			log.Print("Recovered in testPanic2Error", r)
			tx.Rollback()
		}
	}()
}
