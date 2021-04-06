package utils

/***
*上传类封装
***/

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"
)

type Upload struct {
	OssAdminEndPoint  string
	OssAdminBucket    string
	OssAdminAccessKey string
	OssAdminSecret    string
	OssResourceUrl    string
}


//批量上传文件
func (up *Upload) UploadImages(projectTag, filePath string, c *http.Request, keys []string,ossUrl string) (map[string]string, error) {
	tmp := make(map[string]string, 0)
	for _, key := range keys {
		name, data, err := up.readFormFile(c, key)
		if err != nil {
			return nil, err
		}
		downloadUrl, err := up.BatchUploadImageToAliyun(projectTag, filePath, name, data,ossUrl)
		if err != nil {
			return nil, err
		}
		tmp[key] = downloadUrl
	}
	return tmp, nil
}

//单个上传
func (up *Upload) UploadImagesOne(projectTag, filePath string, c *http.Request, keys string,ossUrl string) (string, error) {

	    name, data, err := up.readFormFile(c, keys)
		if err != nil {
			if err == http.ErrMissingFile {
				return "",nil
			}
			return "", err
		}
		downloadUrl, err := up.BatchUploadImageToAliyun(projectTag, filePath, name, data,ossUrl)
		if err != nil {
			return "", err
		}
		
     return downloadUrl,nil
}
//读取form中的上传文件图片或者视频
func (up *Upload) readFormFile(c *http.Request, key string) (string, []byte, error) {
	_,formFile, err := c.FormFile(key)
	if err != nil {
		return "", nil, err
	}
	buf := make([]byte, formFile.Size)
	file, err := formFile.Open()
	if err != nil {
		return "", nil, err
	}
	defer file.Close()
	_, err = file.Read(buf)
	if err != nil {
		return "", nil, err
	}
	return formFile.Filename, buf, nil
}


//上传文件到阿里云
func (up *Upload) BatchUploadImageToAliyun(projectTag, filePath, name string, data []byte,ossUrl string, options ...oss.Option) (string, error) {

	client, err := oss.New(up.OssAdminEndPoint, up.OssAdminAccessKey, up.OssAdminSecret)
	if err != nil {
		return "", err
	}

	var bucketName string = projectTag + up.OssAdminBucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	filePath = up.genPath(filePath, name)

	if err := bucket.PutObject(filePath, bytes.NewReader(data), options...); err != nil {
		return "", err
	}

	if ossUrl == "" {
		ossUrl = up.OssResourceUrl
	}
	
	return ossUrl + "/" + filePath, nil
}

//格式化路径
func  (up *Upload) genPath(filePath string, name string) string {
	now := time.Now()
	y, m, _ := now.Date()
	items := strings.Split(name, ".")
	n := len(items)
	r := rand.Uint32()
	strNow := now.Format("20060102150405") + fmt.Sprint(r)
	return path.Join(filePath, fmt.Sprint(y), fmt.Sprintf("%02d", m), strNow+"."+items[n-1])
}

func NewUpload(endPoint,bucket,accessKey,secret,resourceUrl string) *Upload {
	return &Upload{
		OssAdminEndPoint:endPoint,
		OssAdminBucket:bucket,
		OssAdminAccessKey:accessKey,
		OssAdminSecret:secret,
		OssResourceUrl:resourceUrl,
	}
}
