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

package core

import (
	"content/driven/awsstorage"
	cacheadapter "content/driven/cache"
	"content/driven/storage"
	"content/driven/tempstorage"
	"content/driven/twitter"
	"content/driven/webp"
	"log"
	"sync"
)

//Application represents the core application code based on hexagonal architecture
type Application struct {
	version string
	build   string

	cacheLock *sync.Mutex

	Services Services //expose to the drivers adapters

	storage            Storage
	awsAdapter         *awsstorage.Adapter
	tempStorageAdapter *tempstorage.Adapter
	webpAdapter        *webp.Adapter
	twitterAdapter     *twitter.Adapter
	cacheAdapter       *cacheadapter.CacheAdapter
}

// Start starts the core part of the application
func (app *Application) Start() {
	err := app.storeMultiTenancyData()
	if err != nil {
		log.Fatalf("error initializing multi-tenancy data: %s", err.Error())
	}
}

//as the service starts supporting multi-tenancy we need to add the needed multi-tenancy fields for the existing data,
func (app *Application) storeMultiTenancyData() error {
	//in transaction
	transaction := func(context storage.TransactionContext) error {

		sg, err := app.storage.GetStudentGuides(nil)
		if err != nil {
			return err
		}
		log.Println(sg)

		return nil
	}

	err := app.storage.PerformTransaction(transaction)
	if err != nil {
		log.Printf("error performing transaction for multi tenancy")
		return err
	}
	return nil
}

// NewApplication creates new Application
func NewApplication(version string, build string, storage Storage, awsAdapter *awsstorage.Adapter, tempStorageAdapter *tempstorage.Adapter, webpAdapter *webp.Adapter, twitterAdapter *twitter.Adapter, cacheadapter *cacheadapter.CacheAdapter) *Application {
	cacheLock := &sync.Mutex{}
	application := Application{version: version, build: build, cacheLock: cacheLock, storage: storage, awsAdapter: awsAdapter, tempStorageAdapter: tempStorageAdapter, webpAdapter: webpAdapter, twitterAdapter: twitterAdapter, cacheAdapter: cacheadapter}

	// add the drivers ports/interfaces
	application.Services = &servicesImpl{app: &application}

	return &application
}
