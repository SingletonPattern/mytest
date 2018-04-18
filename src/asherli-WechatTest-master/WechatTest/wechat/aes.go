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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/labstack/gommon/log"
)

// AESCBCEncrypt 采用 CBC 模式的 AES 加密
func AESCBCEncrypt(src, key, iv []byte) (enc []byte, err error) {
	log.Debugf("src: %s", src)
	src = PKCS7Padding(src, len(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(src, src)
	enc = src

	log.Debugf("enc: % x", enc)
	return enc, nil
}

// AESCBCDecrypt 采用 CBC 模式的 AES 解密
func AESCBCDecrypt(enc, key, iv []byte) (src []byte, err error) {
	log.Debugf("enc: % x", enc)
	if len(enc) < len(key) {
		return nil, fmt.Errorf("the length of encrypted message too short: %d", len(enc))
	}
	if len(enc)&(len(key)-1) != 0 { // or len(enc)%len(key) != 0
		return nil, fmt.Errorf("encrypted message is not a multiple of the key size(%d), the length is %d", len(key), len(enc))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(enc, enc)
	src = PKCS7UnPadding(enc)

	log.Debugf("src: %s", src)
	return src, nil
}

func PKCS7Padding(src []byte, k int) (padded []byte) {
	padLen := k - len(src)%k
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}

func PKCS7UnPadding(src []byte) (padded []byte) {
	padLen := int(src[len(src)-1])
	return src[:len(src)-padLen]
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}