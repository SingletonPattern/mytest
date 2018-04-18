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

package wechat

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	"weixin/utils/xrequest"
)

func unmarshalResponseToJson(res xrequest.Response, v interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	apiErr := &ApiError{}

	err = json.Unmarshal(b, apiErr)

	if err != nil {
		return err
	}
	if apiErr.IsError() {
		return apiErr
	}
	return json.Unmarshal(b, v)
}

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetCurrTs() int64 {
	return time.Now().Unix()
}