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

package main

import (
	"content/core"
	"content/core/model"
	"content/driven/awsstorage"
	cacheadapter "content/driven/cache"
	storage "content/driven/storage"
	"content/driven/tempstorage"
	"content/driven/twitter"
	"content/driven/webp"
	driver "content/driver/web"
	"log"
	"os"
	"strings"
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

	port := getEnvKey("CONTENT_PORT", true)

	//mongoDB adapter
	mongoDBAuth := getEnvKey("CONTENT_MONGO_AUTH", true)
	mongoDBName := getEnvKey("CONTENT_MONGO_DATABASE", true)
	mongoTimeout := getEnvKey("CONTENT_MONGO_TIMEOUT", false)
	storageAdapter := storage.NewStorageAdapter(mongoDBAuth, mongoDBName, mongoTimeout)
	err := storageAdapter.Start()
	if err != nil {
		log.Fatal("Cannot start the mongoDB adapter - " + err.Error())
	}

	// S3 Adapter
	s3Bucket := getEnvKey("S3_BUCKET", true)
	s3ProfileImagesBucket := getEnvKey("S3_PROFILE_IMAGES_BUCKET", true)
	s3Region := getEnvKey("S3_REGION", true)
	awsAccessKeyID := getEnvKey("AWS_ACCESS_KEY_ID", true)
	awsSecretAccessKey := getEnvKey("AWS_SECRET_ACCESS_KEY", true)
	awsConfig := &model.AWSConfig{S3Bucket: s3Bucket, S3ProfileImagesBucket: s3ProfileImagesBucket, S3Region: s3Region, AWSAccessKeyID: awsAccessKeyID, AWSSecretAccessKey: awsSecretAccessKey}
	awsAdapter := awsstorage.NewAWSStorageAdapter(awsConfig)

	tempStorageAdapter := tempstorage.NewTempStorageAdapter()

	webpAdapter := webp.NewWebpAdapter()

	defaultCacheExpirationSeconds := getEnvKey("DEFAULT_CACHE_EXPIRATION_SECONDS", false)
	cacheAdapter := cacheadapter.NewCacheAdapter(defaultCacheExpirationSeconds)

	twitterFeedURL := getEnvKey("TWITTER_FEED_URL", true)
	twitterAccessToken := getEnvKey("TWITTER_ACCESS_TOKEN", true)
	twitterAdapter := twitter.NewTwitterAdapter(twitterFeedURL, twitterAccessToken)

	// application
	application := core.NewApplication(Version, Build, storageAdapter, awsAdapter, tempStorageAdapter, webpAdapter, twitterAdapter, cacheAdapter)
	application.Start()

	// web adapter
	host := getEnvKey("CONTENT_HOST", true)
	oidcProvider := getEnvKey("CONTENT_OIDC_PROVIDER", true)
	oidcClientIDs := getEnvKeyAsList("CONTENT_OIDC_CLIENT_IDS", true)
	coreBBHost := getEnvKey("CORE_BB_HOST", true)
	contentServiceURL := getEnvKey("CONTENT_SERVICE_URL", true)

	config := model.Config{
		OidcProvider:      oidcProvider,
		OidcClientIDs:     oidcClientIDs,
		CoreBBHost:        coreBBHost,
		ContentServiceURL: contentServiceURL,
	}

	webAdapter := driver.NewWebAdapter(host, port, application, config)

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
