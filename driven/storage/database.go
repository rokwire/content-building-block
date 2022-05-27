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

package storage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	db       *mongo.Database
	dbClient *mongo.Client

	studentGuides   *collectionWrapper
	healthLocations *collectionWrapper
	contentItems    *collectionWrapper
}

func (m *database) start() error {

	log.Println("database -> start")

	//connect to the database
	clientOptions := options.Client().ApplyURI(m.mongoDBAuth)
	connectContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	client, err := mongo.Connect(connectContext, clientOptions)
	cancel()
	if err != nil {
		return err
	}

	//ping the database
	pingContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	err = client.Ping(pingContext, nil)
	cancel()
	if err != nil {
		return err
	}

	//apply checks
	db := client.Database(m.mongoDBName)

	studentGuides := &collectionWrapper{database: m, coll: db.Collection("student_guides")}
	if err != nil {
		return err
	}

	healthLocations := &collectionWrapper{database: m, coll: db.Collection("health_locations")}
	if err != nil {
		return err
	}

	contentItems := &collectionWrapper{database: m, coll: db.Collection("content_items")}
	if err != nil {
		return err
	}

	//asign the db, db client and the collections
	m.db = db
	m.dbClient = client

	m.studentGuides = studentGuides
	m.healthLocations = healthLocations

	err = m.applyContentItemsChecks(contentItems)
	if err != nil {
		log.Printf("error on applyContentItemsChecks: %s", err)
		return err
	}
	m.contentItems = contentItems

	return nil
}

func (m *database) applyContentItemsChecks(posts *collectionWrapper) error {
	log.Println("apply content_items checks.....")

	indexes, _ := posts.ListIndexes()
	indexMapping := map[string]interface{}{}
	if indexes != nil {

		for _, index := range indexes {
			name := index["name"].(string)
			indexMapping[name] = index
		}
	}
	if indexMapping["category_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "category", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}
	if indexMapping["date_created_1"] == nil {
		err := posts.AddIndex(
			bson.D{
				primitive.E{Key: "date_created", Value: 1},
			}, false)
		if err != nil {
			return err
		}
	}

	log.Println("content_items checks passed")
	return nil
}
