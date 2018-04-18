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
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/labstack/gommon/log"
)

// ValidateURL 验证 URL 以判断来源是否合法
func ValidateURL(token, timestamp, nonce, signature string) bool {
	tmpArr := []string{token, timestamp, nonce}
	sort.Strings(tmpArr)

	tmpStr := strings.Join(tmpArr, "")
	actual := fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))

	log.Debugf("%s %s", tmpArr, actual)
	return actual == signature
}


// Signature 对加密的报文计算签名
func Signature(token, timestamp, nonce, encrypt string) string {
	tmpArr := []string{token, timestamp, nonce, encrypt}
	sort.Strings(tmpArr)

	tmpStr := strings.Join(tmpArr, "")
	actual := fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))

	log.Debugf("%s %s", tmpArr, actual)
	return actual
}

// CheckSignature 验证加密的报文的签名
func CheckSignature(token, timestamp, nonce, encrypt, sign string) bool {
	return Signature(token, timestamp, nonce, encrypt) == sign
}

func Sha1Hex(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}