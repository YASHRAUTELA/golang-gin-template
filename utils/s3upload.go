package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketBasics struct {
	S3Client *s3.Client
}

var S3_URL_FORMAT string = "https://%s.s3.amazonaws.com/%s"
var S3_REGION string = os.Getenv("AWS_S3_REGION")

func GetS3Url(fileName string) string {
	var S3_BUCKET_NAME string = os.Getenv("AWS_S3_BUCKET_NAME")
	return fmt.Sprintf(S3_URL_FORMAT, S3_BUCKET_NAME, fileName)
}

func NewBucketBasics() *BucketBasics {
	key := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(S3_REGION), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, "")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &BucketBasics{
		S3Client: s3.NewFromConfig(cfg),
	}
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func UploadFileToS3(objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	}
	defer file.Close()

	s3ClientBasics := NewBucketBasics()

	var S3_BUCKET_NAME string = os.Getenv("AWS_S3_BUCKET_NAME")
	_, err = s3ClientBasics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(S3_BUCKET_NAME),
		Key:    aws.String(objectKey),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			fileName, S3_BUCKET_NAME, objectKey, err)
	}
	return err
}
