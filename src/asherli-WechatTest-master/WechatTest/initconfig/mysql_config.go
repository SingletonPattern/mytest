/*
 *    Copyright 2016-2018 Li ZongZe
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package initconfig

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
)

var MySQL *xorm.Engine


func init() {
	var err error
	//MySQL, err = xorm.NewEngine(core.MYSQL, "root:@tcp(112.126.88.206:3306)/guoanjiawx?parseTime=true&charset=utf8")
	//MySQL, err = xorm.NewEngine(core.MYSQL, "root:Wang_2015@tcp(127.0.0.1:3306)/guoanjiawx?parseTime=true&charset=utf8")
	MySQL, err = xorm.NewEngine(core.MYSQL, "root:@tcp(127.0.0.1:3306)/guoanjiawx?parseTime=true&charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}

	MySQL.ShowSQL(true)

	logfile, err := os.Create("sql.log")
	if err != nil {
		fmt.Println(err)
		return
	}

	MySQL.SetLogger(xorm.NewSimpleLogger(logfile))

	MySQL.SetMaxIdleConns(10)
	MySQL.SetMaxOpenConns(50)
	err = MySQL.Ping();
	if err != nil {
		fmt.Println(err)
		return
	}
}
