package qiniu

import (
	"context"
	"encoding/base64"
	"os"
	"path"
	"strings"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var pipeline = "amr2mp3";

func Upload(file *os.File) (string, error) {

	key := path.Base(file.Name())

	putPolicy := storage.PutPolicy{}

	mac := qbox.NewMac(accessKey, secretKey)

	if strings.Index(key, ".amr") != -1 {
		key = strings.Replace(key, ".amr", ".mp3", -1)
		putPolicy = storage.PutPolicy{
			Scope:              bucket,
			PersistentOps:      "avthumb/mp3|saveas/" + base64.URLEncoding.EncodeToString([]byte(bucket+":"+key)),
			PersistentPipeline: pipeline,
			Expires:            3600,
		}
	} else {
		putPolicy = storage.PutPolicy{
			Scope: bucket,
		}
	}

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	uploadManager := storage.NewFormUploader(&cfg)

	ret := &storage.PutRet{}

	err := uploadManager.PutFile(context.Background(), ret, upToken, key, file.Name(), nil)
	if err != nil {
		return "", err
	}
	return ret.Key, nil
}
