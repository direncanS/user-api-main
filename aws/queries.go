package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var bucketName = os.Getenv("S3_BUCKET")

func getObjectURL(objectKey string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, objectKey)
}

func GetVideoByFolder(folder string) ([]string, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		fmt.Println("Error while loading the aws config ", err)
		return nil, err
	}

	client := s3.NewFromConfig(config)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folder + "/"), // add "/" to ensure only objects inside the folder are listed
	}

	result, err := client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var list []string
	for _, obj := range result.Contents {
		objURL := getObjectURL(*obj.Key)
		list = append(list, objURL)
	}

	return list, nil
}
