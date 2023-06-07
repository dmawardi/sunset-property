package db

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ObjectRepository interface {
	// Upload file at filepath to object storage and set key to keyPath (eg. "test-folder/broke-world.txt").
	// Returns the ETag & file size of the uploaded file
	UploadFile(filePath string, keyPath string, isPublic bool) (string, int64, error)
}

type objectRepository struct {
	svc *s3.S3
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

	return objectRepository{svc}
}

// Upload file at filepath to object storage and set key to keyPath (eg. "test-folder/broke-world.txt").
// Returns the ETag & file size of the uploaded file
func (s objectRepository) UploadFile(filePath string, keyPath string, isPublic bool) (string, int64, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return "", 0, fmt.Errorf(`failed opening file: %w`, err)
	}
	defer file.Close()

	// Grab file info to get file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return "", 0, fmt.Errorf(`failed getting file info: %w`, err)
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
		Bucket:        aws.String(bucket),  // The path to the directory you want to upload the object to, starting with your Space name.
		Key:           aws.String(keyPath), // Object key, referenced whenever you want to access this file later.
		Body:          file,                // The object's contents.
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
		return "", 0, fmt.Errorf(`failed putting object in S3 storage: %w`, err)
	}

	// Extract eTag from response
	eTag := aws.StringValue(putOutput.ETag)
	fmt.Printf("Successfully uploaded file with ETag: %s\n", eTag)

	return eTag, fileSize, nil
}
