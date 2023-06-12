package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type ObjectRepository interface {
	// Upload file at filepath to object storage and set key to keyPath (eg. "test-folder/broke-world.txt").
	// Returns the updated object key path, ETag & file size of the uploaded file
	UploadFile(filePath string, keyPath string, isPublic bool) (string, string, int64, error)
	// Downloads file at object key path from object storage and saves it to tmp folder with filename
	DownloadTempFile(objectKeyPath string, filename string) (string, error)
}

type objectRepository struct {
	svc        *s3.S3
	downloader *s3manager.Downloader
}

// Create new object service connection
func NewObjectService() ObjectRepository {

	key := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Create a new session
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(os.Getenv("ENDPOINT")), // DigitalOcean Spaces endpoint
		S3ForcePathStyle: aws.Bool(false),
		Region:           aws.String(os.Getenv("AWS_REGION")), // DigitalOcean Spaces region
	})

	if err != nil {
		log.Fatal(err)
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Create a new S3 downloader
	downloader := s3manager.NewDownloader(sess)

	return objectRepository{svc, downloader}
}

// Upload file at filepath to object storage and set key to keyPath (eg. "test-folder/broke-world.txt").
// Returns the new objectKeyPath, ETag & file size of the uploaded file
func (s objectRepository) UploadFile(filePath string, keyPath string, isPublic bool) (string, string, int64, error) {
	// Replace all spaces with underscores
	urlKeyPath := strings.ReplaceAll(keyPath, " ", "_")
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return "", "", 0, fmt.Errorf(`failed opening file: %w`, err)
	}
	defer file.Close()

	// Grab file info to get file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return "", "", 0, fmt.Errorf(`failed getting file info: %w`, err)
	}
	fileSize := fileInfo.Size()

	// Grab env variables
	bucket := os.Getenv("BUCKET")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Determine file permissions based on parameter
	var permissions string
	if isPublic {
		permissions = "public"
	} else {
		permissions = "private"
	}

	// Configure object properties
	object := s3.PutObjectInput{
		Bucket:        aws.String(bucket),     // The path to the directory you want to upload the object to, starting with your Space name.
		Key:           aws.String(urlKeyPath), // Object key, referenced whenever you want to access this file later.
		Body:          file,                   // The object's contents.
		ContentLength: &fileSize,
		ACL:           aws.String(permissions), // Defines Access-control List (ACL) permissions, such as private or public.
		Metadata: map[string]*string{ // Required. Defines metadata tags.
			"x-amz-meta-my-key": aws.String(secretKey),
		},
	}

	// Upload the file
	putOutput, err := s.svc.PutObject(&object)
	if err != nil {
		log.Fatal("Fatal error in putting ", err)
		return "", "", 0, fmt.Errorf(`failed putting object in S3 storage: %w`, err)
	}

	// Extract eTag from response
	eTag := aws.StringValue(putOutput.ETag)
	fmt.Printf("Successfully uploaded file with ETag: %s\n", eTag)

	return urlKeyPath, eTag, fileSize, nil
}

// Downloads file at object key path from object storage and saves it to tmp folder with filename
// Returns the file path of the downloaded file and error
func (s objectRepository) DownloadTempFile(objectKeyPath, fileName string) (string, error) {
	// Create a new file in tmp folder to write the downloaded object
	filePath := fmt.Sprintf("/tmp/%s", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Grab env variables
	bucket := os.Getenv("BUCKET")

	// Download the object from S3 and write it to the file
	_, err = s.downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKeyPath),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download object: %v", err)
	}

	return filePath, nil
}
