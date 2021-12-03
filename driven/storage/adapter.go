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

package storage

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"time"
)

// Adapter implements the Storage interface
type Adapter struct {
	db *database
}

// Start starts the storage
func (sa *Adapter) Start() error {
	err := sa.db.start()
	return err
}

// NewStorageAdapter creates a new storage adapter instance
func NewStorageAdapter(mongoDBAuth string, mongoDBName string, mongoTimeout string) *Adapter {
	timeout, err := strconv.Atoi(mongoTimeout)
	if err != nil {
		log.Println("Set default timeout - 500")
		timeout = 500
	}
	timeoutMS := time.Millisecond * time.Duration(timeout)

	db := &database{mongoDBAuth: mongoDBAuth, mongoDBName: mongoDBName, mongoTimeout: timeoutMS}
	return &Adapter{db: db}
}

// GetStudentGuides retrieves all content items
func (sa *Adapter) GetStudentGuides(ids []string) ([]bson.M, error) {
	filter := bson.D{}
	if len(ids) > 0 {
		filter = bson.D{
			primitive.E{Key: "_id", Value: bson.M{"$in": ids}},
		}
	}

	var result []bson.M
	err := sa.db.studentGuides.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateStudentGuide creates a new student guide record
func (sa *Adapter) CreateStudentGuide(item bson.M) (bson.M, error) {

	id := item["_id"]
	if id == nil {
		item["_id"] = uuid.NewString()
	}

	_, err := sa.db.studentGuides.InsertOne(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetStudentGuide retrieves a student guide record by id
func (sa *Adapter) GetStudentGuide(id string) (bson.M, error) {

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	var result []bson.M
	err := sa.db.studentGuides.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, fmt.Errorf("student guide with id: %s is not found", id)
	}
	return result[0], nil

}

// UpdateStudentGuide updates a student guide record
func (sa *Adapter) UpdateStudentGuide(id string, item bson.M) (bson.M, error) {

	jsonID := item["_id"]
	if jsonID == nil && jsonID != id {
		return nil, fmt.Errorf("attempt to override another object")
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := sa.db.studentGuides.ReplaceOne(filter, item, nil)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteStudentGuide deletes a student guide record with the desired id
func (sa *Adapter) DeleteStudentGuide(id string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	result, err := sa.db.studentGuides.DeleteOne(filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("result is nil for resource item with id " + id)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return fmt.Errorf("error occured while deleting a resource item with id " + id)
	}
	return nil

}

//// Health locations

// GetHealthLocations retrieves all content items
func (sa *Adapter) GetHealthLocations(ids []string) ([]bson.M, error) {
	filter := bson.D{}
	if len(ids) > 0 {
		filter = bson.D{
			primitive.E{Key: "_id", Value: bson.M{"$in": ids}},
		}
	}

	var result []bson.M
	err := sa.db.healthLocations.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateHealthLocation creates a new health location record
func (sa *Adapter) CreateHealthLocation(item bson.M) (bson.M, error) {

	id := item["_id"]
	if id == nil {
		item["_id"] = uuid.NewString()
	}

	_, err := sa.db.healthLocations.InsertOne(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetHealthLocation retrieves a health location record by id
func (sa *Adapter) GetHealthLocation(id string) (bson.M, error) {

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	var result []bson.M
	err := sa.db.healthLocations.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, fmt.Errorf("student guide with id: %s is not found", id)
	}
	return result[0], nil

}

// UpdateHealthLocation updates a health location record
func (sa *Adapter) UpdateHealthLocation(id string, item bson.M) (bson.M, error) {

	jsonID := item["_id"]
	if jsonID == nil && jsonID != id {
		return nil, fmt.Errorf("attempt to override another object")
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := sa.db.healthLocations.ReplaceOne(filter, item, nil)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteHealthLocation deletes a health location record with the desired id
func (sa *Adapter) DeleteHealthLocation(id string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	result, err := sa.db.healthLocations.DeleteOne(filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("result is nil for resource item with id " + id)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return fmt.Errorf("error occured while deleting a resource item with id " + id)
	}
	return nil

}

// Event

func (m *database) onDataChanged(changeDoc map[string]interface{}) {
	if changeDoc == nil {
		return
	}
	log.Printf("onDataChanged: %+v\n", changeDoc)
	ns := changeDoc["ns"]
	if ns == nil {
		return
	}
	nsMap := ns.(map[string]interface{})
	coll := nsMap["coll"]

	if "configs" == coll {
		log.Println("configs collection changed")
	} else {
		log.Println("other collection changed")
	}
}
