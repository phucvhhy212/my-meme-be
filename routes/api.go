package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"example.com/hello/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
// InitRoutes initializes the API routes
func InitRoutes(router *gin.Engine, db *gorm.DB) {
	// Route to create a new submission
	router.POST("/api/submissions", func(c *gin.Context) {
		createSubmission(c, db)
	})
	// Route to get a submission by ID
	router.GET("/api/submissions/:id", func(c *gin.Context) {
		getSubmission(c, db)
	})
	// Route to get a presigned URL for uploading a file to S3
	router.GET("/api/storage/presignedURL/:fileName", func(c *gin.Context) {
		getPresignedURL(c)
	})
	// Route to get the version of a file
	router.GET("/api/storage/getVersion/:fileName", func(c *gin.Context) {
		getVersion(c)
	})
}
// getVersion retrieves the version of a file from S3
func getVersion(c *gin.Context) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	// Create S3 client
	s3Client := s3.NewFromConfig(sdkConfig)
	// List object versions
	request, err := s3Client.ListObjectVersions(context.TODO(), &s3.ListObjectVersionsInput{
		Bucket: aws.String("patht"),
		Prefix: aws.String(c.Param("fileName")),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error getting object version": err.Error()})
	}
	// Return the version ID of the first version (latest version)
	c.JSON(http.StatusOK, request.Versions[0].VersionId)

}
// getPresignedURL generates a presigned URL for uploading a file to S3
func getPresignedURL(c *gin.Context) {
	// Load AWS SDK configuration
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	// Create S3 client
	s3Client := s3.NewFromConfig(sdkConfig)
	// Generate presigned URL
	presignClient := s3.NewPresignClient(s3Client)
	request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("patht"),
		Key:    aws.String(c.Param("fileName")),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error getting presigned URL": err.Error()})
	}
	// Return the presigned URL
	c.JSON(http.StatusOK, request.URL)
}
// createSubmission handles the creation of a new submission
func createSubmission(c *gin.Context, db *gorm.DB) {
	var request models.SubmissionCreateRequest
	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create new submission in the database
	if err := db.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error creating submission": err.Error()})
		return
	}
	// Return the ID of the created submission
	c.JSON(http.StatusOK, request.ID)
}

// getSubmission retrieves a submission by ID
func getSubmission(c *gin.Context, db *gorm.DB) {
	// Convert ID from string to integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var data models.Submissions
	data.ID = uint(id)
	// Find submission in the database
	if err := db.First(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// Return the submission data
	c.JSON(http.StatusOK, data)

}