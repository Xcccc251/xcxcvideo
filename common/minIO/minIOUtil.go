package minIO

import (
	"XcxcVideo/common/helper"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"path"
)

var endpoint = "localhost:9000"
var accessKey = "zdBuAVgzmJQp33nlv05T"
var secretKey = "lvoxgTCdRlY4kxzntPU6oiRjssBXohZ2EQb4NUKL"
var bucketName = "xcxcaudio"
var minioClient = InitMinioClient()

func InitMinioClient() *minio.Client {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
	})
	if err != nil {

		log.Fatalf("Error creating MinIO client: %v", err)
	}
	return client
}

func UploadMP4(fileName string, file *os.File) (finalUrl string, err error) {
	ext := path.Ext(fileName)
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	objectName := helper.GetUUID() + ext
	contentType := "video/mp4"
	uploadInfo, err := minioClient.PutObject(context.Background(), bucketName, objectName, file, fileInfo.Size(),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, uploadInfo.Size)
	return "http://127.0.0.1:9000/xcxcaudio/" + objectName, nil
}
func DelObject(fileName string) error {
	err := minioClient.RemoveObject(context.Background(), bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}
func UploadImage() {}
