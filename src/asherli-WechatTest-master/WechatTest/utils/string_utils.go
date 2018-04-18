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

package utils

import (
	"strings"
)

//判断字符串是否为空
func IsEmptyStr(str string) bool {
	str = strings.TrimSpace(str)
	if str == "" || len(str) == 0 {
		return true
	}
	return false
}

//判断字符串是否不为空
func IsNotEmptyStr(str string) bool {
	return !IsEmptyStr(str)
}

//判断字符串数组是否存在空字符串元素
func IsAnyEmptyStr(str[] string) bool{
	if str==nil||len(str)<=0 {
		return true
	}

	for _,v:=range str{
		if IsEmptyStr(v) {
			return true
		}
	}
	return false
}

//判断字符串数组是否不存在空字符串元素
func IsNoneEmptyStr(str[] string)bool  {
	return !IsAnyEmptyStr(str)
}

//判断是否为空数组
func IsEmptyArray(array [] interface{}) bool {
	if array == nil || len(array) == 0 {
		return true
	}
	return false
}

//判断数组是否不为空
func IsNotEmptyArray(array [] interface{}) bool {
	return !IsEmptyArray(array);
}

//判断字符串数组是否为空
func IsEmptyStrArray(array [] string) bool{
	var paramSlice []interface{}
	for _, param := range array {
		paramSlice = append(paramSlice, param)
	}
	return IsEmptyArray(paramSlice)
}

//判断字符串数组是否不为空
func IsNotEmptyStrArray(array [] string) bool{
	var paramSlice []interface{}
	for _, param := range array {
		paramSlice = append(paramSlice, param)
	}
	return IsNotEmptyArray(paramSlice)
}
