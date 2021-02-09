package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
)

func getFileMD5(fileHeader *multipart.FileHeader) string {
	var md5c = md5.New()
	file, _ := fileHeader.Open()
	defer file.Close()
	r := bufio.NewReader(file)
	if _, err := io.Copy(md5c, r); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(md5c.Sum(nil))
}

func genErrJson(msg string) gin.H {
	return gin.H{
		"code": 400,
		"msg":  msg,
	}
}
