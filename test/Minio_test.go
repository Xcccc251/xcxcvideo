package test

import (
	"XcxcVideo/common/minIO"
	"XcxcVideo/service"
	"context"
	"fmt"
	"testing"
)

func TestMinioClient(t *testing.T) {
	path, err := service.MergeSlices("chunk-bc074d8f00118e5e304cd75019f6a102")
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("path:", path)

	//path := "C:\\Users\\86150\\GolandProjects\\XcXcVideo\\fileDir\\video\\chunk-bc074d8f00118e5e304cd75019f6a102-merged.mp4"
	//f, err := os.Open(path)
	//if err != nil {
	//	fmt.Println("open file err:", err.Error())
	//	return
	//}
	//minIO.UploadMP4("chunk-bc074d8f00118e5e304cd75019f6a102-merged.mp4", f)

	//	err := minIO.DelObject("606179f5-4908-447d-a71b-947f8129c07e.mp4")
	//
	//	if err != nil {
	//		fmt.Println(err)
	//	}
}
func TestListBuckets(t *testing.T) {
	minioClient := minIO.InitMinioClient()
	bucketInfos, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		fmt.Println("List Buckets err：", err.Error())
		return
	}
	for index, bucketInfo := range bucketInfos {
		fmt.Printf("List Bucket No {%d}----filename{%s}-----createTime{%s}\n", index+1, bucketInfo.Name, bucketInfo.CreationDate.Format("2006-01-02 15:04:05"))
	}
}
