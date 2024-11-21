package middleware

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var location string = "" // 暂时用不到
var minioUrl string = ""
var minioPort int = 9000
var minioAccessKey string = ""
var minioSecretKey string = ""

var mimeTypeMap = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".tiff": "image/tiff",
	".mp4":  "video/mp4",
	".avi":  "video/x-msvideo",
	".mkv":  "video/x-matroska",
	".mov":  "video/quicktime",
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".pdf":  "application/pdf",
	".txt":  "text/plain",
	".html": "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".json": "application/json",
	".zip":  "application/zip",
	".rar":  "application/x-rar-compressed",
	".7z":   "application/x-7z-compressed",
}

func NewMinioClient() {
	useSSL := false

	var err error
	minioClient, err = minio.New(fmt.Sprintf("%s:%d", minioUrl, minioPort), &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: useSSL,
	})

	if err != nil || minioClient == nil {
		log.Fatalln(err)
		os.Exit(0)
	}

	log.Printf("Minio Client Init Success") // minioClient is now set up
}

// 判断存储桶是否存在
func BucketExists(bucketName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err == nil {
		log.Println(bucketName, " bucket exists")
		return exists, nil
	} else {
		log.Println(bucketName, " bucket not exists, Error: ", err)
		return exists, err
	}
}

func MakeBucket(bucketName string) error {
	if exists, err := BucketExists(bucketName); exists {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	return err
}

// 根据内容来判断文件 mime 类型
func detectContentType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(filePath, " open fail")
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		fmt.Println(filePath, " read fail")
		return "", err
	}

	return http.DetectContentType(buf), nil
}

func getContentType(filePath string) (string, error) {
	// 获取文件拓展名
	ext := filepath.Ext(filePath)
	if ext == "" {
		return detectContentType(filePath)
	}

	if mimeType, exists := mimeTypeMap[ext]; exists {
		return mimeType, nil
	}
	return "application/octet-stream", nil // default
}

// 上传对象（内存级）
func UploadOjbect(bucketName string, data []byte, n int64, objName string) bool {
	reader := bytes.NewReader(data)

	contentType, _ := getContentType(objName)
	info, err := minioClient.PutObject(context.Background(), bucketName, objName, reader, n, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		log.Fatalln(objName, " upload fail")
		return false
	}
	log.Printf("Successfully uploaded %s of size %d\n", objName, info.Size)
	return true
}

func FUploadObject(bucketName string, filePath string, objName string) bool {
	contentType, _ := getContentType(filePath)
	info, err := minioClient.FPutObject(context.Background(), bucketName, objName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		log.Fatalln(filePath, " upload fail, err: ", err)
		return false
	}

	log.Printf("Successfully uploaded %s of size %d\n", objName, info.Size)
	return true
}

// 获取该用户指定对象的 obj 资源的 url (私有的也可以直接通过这个 URL 下载)
func PresignedObjectURL(bucketName, objName string, expiry time.Duration) (string, error) {
	url, err := minioClient.PresignedGetObject(context.Background(), bucketName, objName, expiry, nil)
	if err != nil {
		log.Println("Get object url fail, reason: ", err)
		return "", err
	}

	return url.String(), nil
}

func ObjectInfo(bucketName, objName string) *minio.ObjectInfo {
	info, err := minioClient.StatObject(context.Background(), bucketName, objName, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil
		}
		log.Println("StatObject fail, Reason: ", err)
		return nil
	}

	return &info
}

func DeleteObject(bucketName, objName string) error {
	err := minioClient.RemoveObject(context.Background(), bucketName, objName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println("Delete object fail, Reason: ", err)
		return err
	}

	return nil
}

// 通过管道的方式来删除，和 MinIO 的交互是并发的，并且逐个传递对象，内存开销不会暴增
func DeleteObjects(bucketName string, objNames []string) {
	objCh := make(chan minio.ObjectInfo)

	// 启动一个协程，将要删除的 obj 逐个放到管道里面
	go func() {
		defer close(objCh)

		for _, objName := range objNames {
			objCh <- minio.ObjectInfo{Key: objName}
		}
	}()

	for err := range minioClient.RemoveObjects(context.Background(), bucketName, objCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			fmt.Println(err.ObjectName, " Delete Fail, Reason: ", err.Err)
		}
	}
}

func init() {
	NewMinioClient()
}
