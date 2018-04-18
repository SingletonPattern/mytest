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
	"crypto/aes"
	"encoding/base64"
	"testing"
)

func TestAESCBCDecrypt(t *testing.T) {
	aesKey := []byte("0123456789abcdef0123456789abcdef")
	src := "我们都是好孩子"

	b64Enc := "MTp5u8m7i4zMqJFlSo1QBIUn+iASUNmd+Co9u0y4Y5w="
	enc, _ := base64.StdEncoding.DecodeString(b64Enc)

	actual, err := AESCBCDecrypt(enc, aesKey, aesKey[:aes.BlockSize])
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if string(actual) != src {
		t.Logf("expect %s, but get %s", src, actual)
		t.FailNow()
	}
}

func TestAESCBCEncryptAndDecrypt(t *testing.T) {
	aesKey := []byte("0123456789abcdef0123456789abcdef")
	src := "我们都是好孩子"

	enc, err := AESCBCEncrypt([]byte(src), aesKey, aesKey[:aes.BlockSize])
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("%s", base64.StdEncoding.EncodeToString(enc))

	actual, err := AESCBCDecrypt(enc, aesKey, aesKey[:aes.BlockSize])
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if string(actual) != src {
		t.Logf("expect %s, but get %s", src, actual)
		t.FailNow()
	}
}
