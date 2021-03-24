package main

import (
	// "html/template"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var MyBucket string
var filepath string

//GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func ConnectAws() *session.Session {

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})

	if err != nil {
		panic(err)
	}

	return sess
}

func UploadDocument(c *gin.Context) {
	sess := c.MustGet("sess").(*session.Session)
	uploader := s3manager.NewUploader(sess)

	file, header, err := c.Request.FormFile("document")
	filename := header.Filename

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to upload file",
			"uploader": up,
		})
		return
	}
	filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath": filepath,
		"success":  "Document uploaded successfully",
	})
}

func getListOfFiles(c *gin.Context) {
	sess := c.MustGet("sess").(*session.Session)
	svc := s3.New(sess)
	documentList, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(MyBucket)})

	if err != nil {
		fmt.Println("Error in geting list of documents")
	}

	for _, item := range documentList.Contents {
		filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + *item.Key
		c.JSON(http.StatusOK, gin.H{
			"Name":          *item.Key,
			"LastModified":  *item.LastModified,
			"Size":          *item.Size,
			"Downlaod Link": filepath,
			"StorageClass":  *item.StorageClass,
		})
	}
}

func downloadSingleDocument(c *gin.Context) {
	sess := c.MustGet("sess").(*session.Session)
	fileName := c.Request.FormValue("document")
	svc := s3.New(sess)
	documentList, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(MyBucket)})
	//success := svc.DeleteObject()
	if err != nil {
		fmt.Println("Error in geting list of documents")
	}

	for _, item := range documentList.Contents {
		if *item.Key == fileName {
			filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + *item.Key
			c.JSON(http.StatusOK, gin.H{
				"Name":          *item.Key,
				"LastModified":  *item.LastModified,
				"Size":          *item.Size,
				"Downlaod Link": filepath,
				"StorageClass":  *item.StorageClass,
			})
			return
		}
	}
}

func deleteObject(c *gin.Context) {
	sess := c.MustGet("sess").(*session.Session)
	fileName := c.Request.FormValue("document")
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"Status": err,
		})
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"Status": err,
		})
	} else {
		c.JSON(http.StatusBadGateway, gin.H{
			"Status": "Successfully deleted",
		})
	}
}

func main() {
	LoadEnv()
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = GetEnvWithKey("AWS_REGION")
	MyBucket = GetEnvWithKey("BUCKET_NAME")

	sess := ConnectAws()
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sess", sess)
		c.Next()
	})

	router.POST("/upload", UploadDocument)
	router.GET("/list", getListOfFiles)
	router.POST("/download", downloadSingleDocument)
	router.POST("/delete", deleteObject)
	_ = router.Run(":4000")
}
