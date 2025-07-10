// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"content/core"
	"content/core/model"
	"content/driven/awsstorage"
	cacheadapter "content/driven/cache"
	corebb "content/driven/core"
	storage "content/driven/storage"
	"content/driven/twitter"
	driver "content/driver/web"
	"log"
	"os"
	"strconv"
	"strings"

	rokwireAuth "github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth"
	"github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/keys"
	"github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/sigauth"
	"github.com/rokwire/rokwire-building-block-sdk-go/utils/logging/logs"
)

var (
	// Version : version of this executable
	Version string
	// Build : build date of this executable
	Build string
)

func main() {
	if len(Version) == 0 {
		Version = "dev"
	}

	serviceID := "content"

	loggerOpts := logs.LoggerOpts{SuppressRequests: logs.NewStandardHealthCheckHTTPRequestProperties(serviceID + "/version")}
	logger := logs.NewLogger(serviceID, &loggerOpts)

	port := getEnvKey("CONTENT_PORT", true)

	//common
	coreBBHost := getEnvKey("CONTENT_CORE_BB_HOST", true)
	contentServiceURL := getEnvKey("CONTENT_SERVICE_URL", true)
	authService := rokwireAuth.Service{
		ServiceID:   serviceID,
		ServiceHost: contentServiceURL,
		FirstParty:  true,
		AuthBaseURL: coreBBHost,
	}
	serviceRegLoader, err := rokwireAuth.NewRemoteServiceRegLoader(&authService, []string{"auth"})
	if err != nil {
		log.Fatalf("Error initializing remote service registration loader: %v", err)
	}
	serviceRegManager, err := rokwireAuth.NewServiceRegManager(&authService, serviceRegLoader, true)
	if err != nil {
		log.Fatalf("Error initializing service registration manager: %v", err)
	}
	//end common

	//mongoDB adapter
	mongoDBAuth := getEnvKey("CONTENT_MONGO_AUTH", true)
	mongoDBName := getEnvKey("CONTENT_MONGO_DATABASE", true)
	mongoTimeout := getEnvKey("CONTENT_MONGO_TIMEOUT", false)
	storageAdapter := storage.NewStorageAdapter(mongoDBAuth, mongoDBName, mongoTimeout, logger)
	err = storageAdapter.Start()
	if err != nil {
		log.Fatal("Cannot start the mongoDB adapter - " + err.Error())
	}

	// S3 Adapter
	s3Bucket := getEnvKey("CONTENT_S3_BUCKET", true)
	s3ProfileImagesBucket := getEnvKey("CONTENT_S3_PROFILE_IMAGES_BUCKET", true)
	s3UsersAudiosBucket := getEnvKey("CONTENT_S3_USERS_AUDIOS_BUCKET", true)
	s3Region := getEnvKey("CONTENT_S3_REGION", true)
	awsAccessKeyID := getEnvKey("CONTENT_AWS_ACCESS_KEY_ID", true)
	awsSecretAccessKey := getEnvKey("CONTENT_AWS_SECRET_ACCESS_KEY", true)
	awsConfig := &model.AWSConfig{S3Bucket: s3Bucket,
		S3ProfileImagesBucket: s3ProfileImagesBucket,
		S3UsersAudiosBucket:   s3UsersAudiosBucket,
		S3Region:              s3Region, AWSAccessKeyID: awsAccessKeyID, AWSSecretAccessKey: awsSecretAccessKey}

	presignExpirationMinutesVal := getEnvKey("CONTENT_S3_REQUEST_PRESIGN_EXPIRATION_MINUTES", false)
	presignExpirationMinutes, err := strconv.Atoi(presignExpirationMinutesVal)
	if err != nil {
		logger.Warnf("error parsing S3 request presign expiration minutes: %s - applying default", err.Error())
	}
	awsAdapter := awsstorage.NewAWSStorageAdapter(awsConfig, presignExpirationMinutes)

	defaultCacheExpirationSeconds := getEnvKey("CONTENT_DEFAULT_CACHE_EXPIRATION_SECONDS", false)
	cacheAdapter := cacheadapter.NewCacheAdapter(defaultCacheExpirationSeconds)

	twitterFeedURL := getEnvKey("CONTENT_TWITTER_FEED_URL", true)
	twitterAccessToken := getEnvKey("CONTENT_TWITTER_ACCESS_TOKEN", true)
	twitterAdapter := twitter.NewTwitterAdapter(twitterFeedURL, twitterAccessToken)

	mtAppID := getEnvKey("CONTENT_MULTI_TENANCY_APP_ID", true)
	mtOrgID := getEnvKey("CONTENT_MULTI_TENANCY_ORG_ID", true)

	//core adapter
	serviceAccountID := getEnvKey("CONTENT_SERVICE_ACCOUNT_ID", false)
	authPrivKeyPemString := getEnvKey("CONTENT_PRIV_KEY", false)
	var authPrivKeyPem string
	if authPrivKeyPemString != "" {
		//make it to be a single line - AWS environemnt variable issue
		authPrivKeyPem = strings.ReplaceAll(authPrivKeyPemString, `\n`, "\n")
	} else {
		log.Fatalf("APPOINTMENTS_PRIV_KEY environment variable is not set")
	}
	alg := keys.RS256
	privKey, err := keys.NewPrivKey(alg, authPrivKeyPem)
	if err != nil {
		logger.Fatalf("Failed to parse auth priv key: %v", err)
	}
	signatureAuth, err := sigauth.NewSignatureAuth(privKey, serviceRegManager, false, true)
	if err != nil {
		log.Fatalf("Error initializing signature auth: %v", err)
	}
	serviceAccountLoader, err := rokwireAuth.NewRemoteServiceAccountLoader(&authService, serviceAccountID, signatureAuth)
	if err != nil {
		log.Fatalf("Error initializing remote service account loader: %v", err)
	}
	serviceAccountManager, err := rokwireAuth.NewServiceAccountManager(&authService, serviceAccountLoader)
	if err != nil {
		log.Fatalf("Error initializing service account manager: %v", err)
	}
	coreAdapter := corebb.NewCoreAdapter(coreBBHost, serviceAccountManager)

	// application
	application := core.NewApplication(Version, Build, storageAdapter, awsAdapter, twitterAdapter, cacheAdapter, mtAppID, mtOrgID, serviceID, coreAdapter, logger)
	application.Start()

	// web adapter
	webAdapter := driver.NewWebAdapter(contentServiceURL, port, application, serviceRegManager, logger)

	webAdapter.Start()
}

func getEnvKeyAsList(key string, required bool) []string {
	stringValue := getEnvKey(key, required)

	// it is comma separated format
	stringListValue := strings.Split(stringValue, ",")
	if len(stringListValue) == 0 && required {
		log.Fatalf("missing or empty env var: %s", key)
	}

	return stringListValue
}

func getEnvKey(key string, required bool) string {
	// get from the environment
	value, exist := os.LookupEnv(key)
	if !exist {
		if required {
			log.Fatal("No provided environment variable for " + key)
		} else {
			log.Printf("No provided environment variable for " + key)
		}
	}
	return value
}
