package main

import (
	"minio_test/middleware"
	// "os"
)

func main() {
	middleware.MakeBucket("testbucket")

	// 上传本地文件
	// middleware.FUploadObject("testbucket", "./lani.png", "newname2")

	// 从内存中传数据给 MinIO
	// buf := make([]byte, 3600000)
	// f, _ := os.Open("./lani.png")
	// n, _ := f.Read(buf)
	// middleware.UploadOjbect("testbucket", buf, int64(n), "lani.png")

	// 获取目标对象 URL
	// url, _ := middleware.ObjectURL("testbucket", "lani.png", 1*time.Hour)
	// fmt.Println(url)

	// 获取对象信息
	// info := middleware.ObjectInfo("testbucket", "lani.png")
	// info := middleware.ObjectInfo("testbucket", "lani.png")
	// if info != nil {
	// 	fmt.Printf("Object: %s\n", info.Key)
	// 	fmt.Printf("Size: %d bytes\n", info.Size)
	// 	fmt.Printf("Last Modified: %v\n", info.LastModified)
	// 	fmt.Printf("ETag: %s\n", info.ETag)
	// } else {
	// 	fmt.Println("Object Nil")
	// }

	middleware.DeleteObjects("testbucket", []string{"newname", "newname2"})
}
