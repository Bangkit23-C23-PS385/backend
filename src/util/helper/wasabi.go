package helper

import (
	"mime/multipart"
	"os"
	"ta/backend/src/constant"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

func createSession() *session.Session {
	config := aws.Config{
		Region:           aws.String(constant.WasabiBucketRegion),
		Endpoint:         aws.String(constant.WasabiBucketEndpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			constant.WasabiAccessKey,
			constant.WasabiSecretKey,
			"",
		),
	}
	options := session.Options{
		Config: config,
	}
	sess := session.Must(session.NewSessionWithOptions(options))

	return sess
}

func CheckIfUserFolderExists(folderName string) (isExists bool, err error) {
	sess := createSession()

	s3Svc := s3.New(sess)
	resp, err := s3Svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  aws.String(constant.WasabiBucketName),
		Prefix:  aws.String(folderName),
		MaxKeys: aws.Int64(1),
	})
	if err != nil {
		err = errors.Wrap(err, "wasabi: check folder exists")
		return
	}
	if len(resp.Contents) > 0 {
		isExists = true
	}

	return
}

func CreateFolder(folderName string) (err error) {
	sess := createSession()
	_ = s3.New(sess)
	folderCreator := s3manager.NewUploader(sess)

	devNull, _ := os.Open(os.DevNull)

	_, err = folderCreator.Upload(&s3manager.UploadInput{
		Bucket: aws.String(constant.WasabiBucketName),
		Key:    aws.String(folderName),
		Body:   aws.ReadSeekCloser(devNull),
	})
	if err != nil {
		err = errors.Wrap(err, "wasabi: create folder")
		return
	}

	return
}

func UploadFile(fileName string, file multipart.File) (err error) {
	sess := createSession()
	s3Svc := s3.New(sess)

	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(constant.WasabiBucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		err = errors.Wrap(err, "wasabi: upload file")
	}

	return
}

func GenerateGetPresignedURL(fileName, contentType string) (url string, err error) {
	sess := createSession()
	s3svc := s3.New(sess)

	req, _ := s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(constant.WasabiBucketName),
		Key:                        aws.String(fileName),
		ResponseContentDisposition: aws.String("inline"),
		ResponseContentType:        aws.String(contentType),
	})

	url, err = req.Presign(24 * time.Hour)
	if err != nil {
		err = errors.Wrap(err, "wasabi: generate get presigned url")
	}

	return
}

func GenerateUploadPresignedURL(fileName string) (url string, err error) {
	sess := createSession()
	s3svc := s3.New(sess)

	req, _ := s3svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(constant.WasabiBucketName),
		Key:    aws.String(fileName),
	})

	url, err = req.Presign(6 * time.Hour)
	if err != nil {
		err = errors.Wrap(err, "wasabi: generate upload presigned url")
	}

	return
}
