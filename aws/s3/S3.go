package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mikel-at-tatari/tatari-dev-db/aws/config"
	"github.com/mikel-at-tatari/tatari-dev-db/log"
	"os"
)

type S3Object struct {
	Bucket string
	Key    string
}

func Download(s3Object S3Object, file *os.File) error {
	svc := s3.NewFromConfig(*config.AWSConfig)
	downloader := manager.NewDownloader(svc)

	numBytes, err := downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(s3Object.Bucket),
		Key:    aws.String(s3Object.Key),
	})

	if err != nil {
		log.Error(err)
		return err
	}

	log.WithFields(
		log.Fields{
			"file":     file.Name(),
			"numBytes": numBytes,
		},
	).Debug("Finished Downloading")

	return nil
}

func Upload(s3Object S3Object, f *string) error {

	file, err := os.Open(*f)

	if err != nil {
		return err
	}

	defer file.Close()

	svc := s3.NewFromConfig(*config.AWSConfig)

	uploader := manager.NewUploader(svc)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3Object.Bucket),
		Key:    aws.String(s3Object.Key),
		Body:   file,
	})

	if err != nil {
		log.Error(err)
		return err
	}

	log.WithFields(log.Fields{"file": file.Name(), "key": *result.Key}).Debug("Finished Uploading")

	return nil
}
