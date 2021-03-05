package main

import (
	"Airstack/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"path/filepath"
)

var db = connectDB(DatabaseName)

func main() {
	err := db.AutoMigrate(&Resource{})
	if err != nil {
		panic(err.Error())
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.Cors()) // remove, if you need not.
	router.MaxMultipartMemory = MaxFileSize
	api := router.Group("/api")
	{
		api.GET("/ping", func(context *gin.Context) {
			context.String(http.StatusOK, "pong")
		})

		api.POST("/upload", upload)
		api.GET("/download/:code", download)
	}

	log.Println("Airstack start!")
	err = router.Run(Host)
	if err != nil {
		panic(err.Error())
	}
}

func upload(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, genErrJson(err.Error()))
		return
	}
	fileHash := getFileMD5(file)
	filename := filepath.Base(file.Filename)
	log.Printf("Hash: %s   Filename: %s\n", fileHash, filename)
	if err := context.SaveUploadedFile(file, StoragePath+fileHash); err != nil {
		context.JSON(http.StatusBadRequest, genErrJson(err.Error()))
		return
	}

	resource := Resource{Hash: fileHash, ResourceName: filename, Size: file.Size, ResourceType: 0}
	db.Create(&resource)

	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"pwd":  genCode(resource.ID),
	})
}

func download(context *gin.Context) {
	code := context.Param("code")
	id := decode(code)
	if id == -1 {
		context.JSON(http.StatusBadRequest, genErrJson("Download Code Error!"))
		return
	}
	var resource Resource
	result := db.First(&resource, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusBadRequest, genErrJson("The resource does not exist!"))
		return
	}
	context.FileAttachment(StoragePath+resource.Hash, resource.ResourceName)
	db.Delete(&Resource{}, id) // delete resource
	return
}


