/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package awsstorage

import (
	"content/core/model"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
)

//Adapter implements the Storage interface
type Adapter struct {
	config *model.AWSConfig
}

//NewAWSStorageAdapter creates a new storage adapter instance
func NewAWSStorageAdapter(config *model.AWSConfig) *Adapter {
	//return &Adapter{S3Bucket: S3Bucket, S3Region: S3Region, AWSAccessKeyID: AWSAccessKeyID, AWSSecretAccessKey: AWSSecretAccessKey}
	return &Adapter{config: config}
}

// CreateImage uploads an image instance from a file and image type
func (a *Adapter) CreateImage(file *os.File, path string) (*string, error) {
	log.Println("Create image")

	s, err := a.createS3Session()
	if err != nil {
		log.Printf("Could not create S3 session")
		return nil, err
	}
	key := a.prepareKey(path)
	objectLocation, err := a.uploadFileToS3(s, file, key)
	if err != nil {
		log.Printf("Could not upload file")
		return nil, err
	}

	return &objectLocation, nil
}

func (a *Adapter) prepareKey(path string) string {
	fileName, _ := uuid.NewUUID() // add uuid for file name
	if strings.HasSuffix(path, "/") {
		return path + fmt.Sprintf("%s", fileName) + ".webp"
	}
	return path + "/" + fmt.Sprintf("%s", fileName) + ".webp"
}

func (a *Adapter) createS3Session() (*session.Session, error) {
	region := a.config.S3Region
	accessKeyID := a.config.AWSAccessKeyID
	secretAccessKey := a.config.AWSSecretAccessKey
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			""),
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return s, nil
}

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func (a *Adapter) uploadFileToS3(s *session.Session, file *os.File, key string) (string, error) {
	uploader := s3manager.NewUploader(s)
	bucket := a.config.S3Bucket
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Print(err)
		return "", err
	}
	return result.Location, err
}
