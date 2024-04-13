package upload

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"opfile/common"
	"os"
	"path/filepath"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 只能上传文件
func UploadFile(ossCfg common.OssCfg, filePath string) (string, error) {
	u, _ := url.Parse(ossCfg.BatchURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossCfg.TmpSecretID,
			SecretKey:    ossCfg.TmpSecretKey,
			SessionToken: ossCfg.SessionToken,
		},
	})
	fleAbsPath, _ := filepath.Abs(filePath)
	fileName := filepath.Base(fleAbsPath)
	fmt.Println("fileName:", fileName)
	fmt.Println("fleAbsPath:", fleAbsPath)

	if is, err := IsFile(fleAbsPath); err != nil {
		return "", err
	} else {
		if is {
			// 从命令行获取参数
			key := "file/" + fileName //对象键（Key）是对象在存储桶中的唯一标识

			opt := &cos.MultiUploadOptions{ThreadPoolSize: 8, CheckPoint: true}
			resp, _, err := client.Object.Upload(
				context.Background(), key, fleAbsPath, opt,
			)
			if err != nil {
				return "", err
			}
			return resp.Key, nil
		} else {
			return "", errors.New("this is a folder not a file")
		}
	}
}

// 通过 tag 的方式，用户可以将请求参数或者请求头部放进签名中。
type URLToken struct {
	SessionToken string `url:"x-cos-security-token,omitempty" header:"-"`
}

func GetPresignedURL(ossCfg common.OssCfg, key string) {
	tak := ossCfg.TmpSecretID
	tsk := ossCfg.TmpSecretKey
	token := &URLToken{
		SessionToken: ossCfg.SessionToken,
	}

	u, _ := url.Parse(ossCfg.BatchURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{})
	ctx := context.Background()

	// 方法2 通过 tag 设置 x-cos-security-token
	// 获取预签名
	presignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodGet, key, tak, tsk, 10*time.Minute, token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("DownloadRUL: ", presignedURL.String())

}

// 在Linux下判断一个文件是否是文件夹
func IsFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if fileInfo.IsDir() {
		return false, nil
	} else {
		return true, nil
	}
}
