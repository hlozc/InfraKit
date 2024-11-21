package main

import (
	"minio_test/middleware"
	"os"
)

func main() {
	middleware.NewMinioClient()
	middleware.MakeBucket("testbucket")

	// middleware.FUploadObject("testbucket", "./lani.png", "newname")
	buf := make([]byte, 3600000)
	f, _ := os.Open("./lani.png")
	n, _ := f.Read(buf)
	middleware.UploadOjbect("testbucket", buf, int64(n), "lani.png")
}
