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

package qiniu

import (
	"github.com/labstack/gommon/log"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var (
	accessKey = "HrfSUQLnauuNKiErBrP_lBGlPBTfWQlqZqhP76-a"
	secretKey = "nvmXyeEhKRCA2o7LAthnM37uy-vfdAmRsfCr7yLt"
	bucket    = "guoanjia"
)

func Fetch(resUrl, fileName string) string {
	mac := qbox.NewMac(accessKey, secretKey)

	cfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	fetchRet, err := bucketManager.Fetch(resUrl, bucket, fileName)
	if err != nil {
		log.Error("七牛云存储 fetch 异常", err)
		return ""
	} else {
		return fetchRet.Key
	}
}

func FetchWithoutKey(resUrl string) string {
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	fetchRet, err := bucketManager.FetchWithoutKey(resUrl, bucket)
	if err != nil {
		log.Error("七牛云存储 fetch 异常", err)
		return ""
	} else {
		return fetchRet.Key
	}
}
