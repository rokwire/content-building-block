package model

// AWSConfig wrapper for all S3 configuration keys
type AWSConfig struct {
	S3Bucket              string
	S3ProfileImagesBucket string
	S3Region              string
	AWSAccessKeyID        string
	AWSSecretAccessKey    string
}
