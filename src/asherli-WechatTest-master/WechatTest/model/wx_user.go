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

package model

import (
	"time"
)

type WxUser struct {
	Id          string    `xorm:"not null pk VARCHAR(36)"`
	UOpenid     string    `xorm:"not null unique VARCHAR(50)"`
	UNickname   string    `xorm:"VARCHAR(50)"`
	USex        string    `xorm:"VARCHAR(10)"`
	UProvince   string    `xorm:"VARCHAR(50)"`
	UCity       string    `xorm:"VARCHAR(50)"`
	UCountry    string    `xorm:"VARCHAR(50)"`
	UHeadimg    string    `xorm:"VARCHAR(255)"`
	UUnionid    string    `xorm:"VARCHAR(255)"`
	UCreateTime time.Time `xorm:"DATETIME"`
}
