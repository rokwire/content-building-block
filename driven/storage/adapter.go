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
	"bytes"
	"content/core/interfaces"
	"content/core/model"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rokwire/rokwire-building-block-sdk-go/utils/errors"
	"github.com/rokwire/rokwire-building-block-sdk-go/utils/logging/logs"
	"github.com/rokwire/rokwire-building-block-sdk-go/utils/logging/logutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Adapter implements the Storage interface
type Adapter struct {
	db      *database
	context mongo.SessionContext
}

// Start starts the storage
func (sa *Adapter) Start() error {
	err := sa.db.start()
	return err
}

// PerformTransaction performs a transaction
func (sa *Adapter) PerformTransaction(transaction func(storage interfaces.Storage) error) error {
	// transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		adapter := sa.withContext(sessionContext)

		err := transaction(adapter)
		if err != nil {
			if wrappedErr, ok := err.(interface {
				Internal() error
			}); ok && wrappedErr.Internal() != nil {
				return nil, wrappedErr.Internal()
			}
			return nil, err
		}

		return nil, nil
	}

	session, err := sa.db.dbClient.StartSession()
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionStart, "mongo session", nil, err)
	}
	context := context.Background()
	defer session.EndSession(context)

	_, err = session.WithTransaction(context, callback)
	if err != nil {
		return errors.WrapErrorAction("performing", logutils.TypeTransaction, nil, err)
	}
	return nil
}

// GetStudentGuides retrieves all content items
func (sa *Adapter) GetStudentGuides(appID string, orgID string, ids []string) ([]bson.M, error) {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID}}
	if len(ids) > 0 {
		filter = bson.D{
			primitive.E{Key: "_id", Value: bson.M{"$in": ids}},
		}
	}

	var result []bson.M
	err := sa.db.studentGuides.Find(sa.context, filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateStudentGuide creates a new student guide record
func (sa *Adapter) CreateStudentGuide(appID string, orgID string, item bson.M) (bson.M, error) {

	id := item["_id"]
	if id == nil {
		item["_id"] = uuid.NewString()
	}
	item["app_id"] = appID
	item["org_id"] = orgID

	_, err := sa.db.studentGuides.InsertOne(sa.context, &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetStudentGuide retrieves a student guide record by id
func (sa *Adapter) GetStudentGuide(appID string, orgID string, id string) (bson.M, error) {

	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	var result []bson.M
	err := sa.db.studentGuides.Find(sa.context, filter, &result, nil)
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
func (sa *Adapter) UpdateStudentGuide(appID string, orgID string, id string, item bson.M) (bson.M, error) {
	jsonID := item["_id"]
	if jsonID == nil && jsonID != id {
		return nil, fmt.Errorf("attempt to override another object")
	}

	item["app_id"] = appID
	item["org_id"] = orgID

	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	err := sa.db.studentGuides.ReplaceOne(sa.context, filter, item, nil)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteStudentGuide deletes a student guide record with the desired id
func (sa *Adapter) DeleteStudentGuide(appID string, orgID string, id string) error {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	result, err := sa.db.studentGuides.DeleteOne(sa.context, filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("result is nil for resource item with id %s", id)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return fmt.Errorf("error occured while deleting a resource item with id %s", id)
	}
	return nil

}

//// Health locations

// GetHealthLocations retrieves all content items
func (sa *Adapter) GetHealthLocations(appID string, orgID string, ids []string) ([]bson.M, error) {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID}}
	if len(ids) > 0 {
		filter = bson.D{
			primitive.E{Key: "_id", Value: bson.M{"$in": ids}},
		}
	}

	var result []bson.M
	err := sa.db.healthLocations.Find(sa.context, filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateHealthLocation creates a new health location record
func (sa *Adapter) CreateHealthLocation(appID string, orgID string, item bson.M) (bson.M, error) {

	id := item["_id"]
	if id == nil {
		item["_id"] = uuid.NewString()
	}
	item["app_id"] = appID
	item["org_id"] = orgID

	_, err := sa.db.healthLocations.InsertOne(sa.context, &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetHealthLocation retrieves a health location record by id
func (sa *Adapter) GetHealthLocation(appID string, orgID string, id string) (bson.M, error) {

	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	var result []bson.M
	err := sa.db.healthLocations.Find(sa.context, filter, &result, nil)
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
func (sa *Adapter) UpdateHealthLocation(appID string, orgID string, id string, item bson.M) (bson.M, error) {
	jsonID := item["_id"]
	if jsonID == nil && jsonID != id {
		return nil, fmt.Errorf("attempt to override another object")
	}
	item["app_id"] = appID
	item["org_id"] = orgID

	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	err := sa.db.healthLocations.ReplaceOne(sa.context, filter, item, nil)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteHealthLocation deletes a health location record with the desired id
func (sa *Adapter) DeleteHealthLocation(appID string, orgID string, id string) error {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	result, err := sa.db.healthLocations.DeleteOne(sa.context, filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("result is nil for resource item with id %s", id)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return fmt.Errorf("error occured while deleting a resource item with id %s", id)
	}
	return nil

}

// Content Items

type getContentItemsCategoriesData struct {
	CategoryName string `json:"_id" bson:"_id"`
}

// GetContentItemsCategories  retrieve all content item categories
func (sa *Adapter) GetContentItemsCategories(appID *string, orgID string) ([]string, error) {
	pipeline := primitive.A{
		bson.M{"$match": bson.M{"app_id": appID, "org_id": orgID}},
		bson.M{"$group": bson.M{"_id": "$category"}},
	}
	var data []getContentItemsCategoriesData
	categories := []string{}

	err := sa.db.contentItems.Aggregate(sa.context, pipeline, &data, &options.AggregateOptions{})
	if err != nil {
		return nil, err
	}
	if data != nil && len(data) > 0 {
		for _, dataEntry := range data {
			categories = append(categories, dataEntry.CategoryName)
		}
	}

	return categories, nil
}

// FindContentItems finds content items
func (sa *Adapter) FindContentItems(appID *string, orgID string, ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItem, error) {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID}}
	if len(ids) > 0 {
		filter = append(filter, primitive.E{Key: "_id", Value: bson.M{"$in": ids}})
	}
	if categoryList != nil && len(categoryList) > 0 {
		filter = append(filter, primitive.E{Key: "category", Value: bson.M{"$in": categoryList}})
	}

	findOptions := options.Find()
	if order != nil && "desc" == *order {
		findOptions.SetSort(bson.M{"date_created": -1})
	} else {
		findOptions.SetSort(bson.M{"date_created": 1})
	}
	if limit != nil {
		findOptions.SetLimit(*limit)
	}
	if offset != nil {
		findOptions.SetSkip(*offset)
	}

	var result []model.ContentItem
	err := sa.db.contentItems.Find(sa.context, filter, &result, findOptions)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetContentItems retrieves all content items
func (sa *Adapter) GetContentItems(appID *string, orgID string, ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItemResponse, error) {

	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID}}
	if appID != nil {
		filter = append(filter, primitive.E{Key: "app_id", Value: appID})
	}
	if len(ids) > 0 {
		filter = append(filter, primitive.E{Key: "_id", Value: bson.M{"$in": ids}})
	}
	if categoryList != nil && len(categoryList) > 0 {
		filter = append(filter, primitive.E{Key: "category", Value: bson.M{"$in": categoryList}})
	}

	findOptions := options.Find()
	if order != nil && "desc" == *order {
		findOptions.SetSort(bson.M{"date_created": -1})
	} else {
		findOptions.SetSort(bson.M{"date_created": 1})
	}
	if limit != nil {
		findOptions.SetLimit(*limit)
	}
	if offset != nil {
		findOptions.SetSkip(*offset)
	}

	var result []model.ContentItemResponse
	err := sa.db.contentItems.Find(sa.context, filter, &result, findOptions)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateContentItem creates a new content item record
func (sa *Adapter) CreateContentItem(item model.ContentItem) (*model.ContentItem, error) {
	_, err := sa.db.contentItems.InsertOne(sa.context, &item)
	if err != nil {
		log.Printf("error create content item: %s", err)
		return nil, err
	}
	return &item, nil
}

// GetContentItem retrieves a content item record by id
func (sa *Adapter) GetContentItem(appID *string, orgID string, id string) (*model.ContentItemResponse, error) {

	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	var result []model.ContentItemResponse
	err := sa.db.contentItems.Find(sa.context, filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		log.Printf("content item with id: %s is not found", id)
		return nil, fmt.Errorf("content item with id: %s is not found", id)
	}
	return &result[0], nil

}

// UpdateContentItem updates a content item record
func (sa *Adapter) UpdateContentItem(appID *string, orgID string, id string,
	category string, data interface{}) (*model.ContentItem, error) {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "category", Value: category},
			primitive.E{Key: "data", Value: data},
			primitive.E{Key: "date_updated", Value: time.Now().UTC()},
		}},
	}
	_, err := sa.db.contentItems.UpdateOne(sa.context, filter, update, nil)
	if err != nil {
		log.Printf("error updating content item: %s", err)
		return nil, err
	}

	//get it to return the updated object
	var result []model.ContentItem
	err = sa.db.contentItems.Find(sa.context, filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		log.Printf("content item with id: %s is not found", id)
		return nil, fmt.Errorf("content item with id: %s is not found", id)
	}
	return &result[0], nil
}

// DeleteContentItem deletes a content item record with the desired id
func (sa *Adapter) DeleteContentItem(appID *string, orgID string, id string) error {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id}}
	result, err := sa.db.contentItems.DeleteOne(sa.context, filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("result is nil for resource item with id %s", id)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return fmt.Errorf("error occured while deleting a resource item with id %s", id)
	}
	return nil
}

// SaveContentItem saves content item
func (sa *Adapter) SaveContentItem(item model.ContentItem) error {
	filter := bson.D{primitive.E{Key: "org_id", Value: item.OrgID},
		primitive.E{Key: "_id", Value: item.ID}}
	if item.AppID != nil {
		filter = append(filter, primitive.E{Key: "app_id", Value: item.AppID})
	}

	opts := options.Replace().SetUpsert(true)
	err := sa.db.contentItems.ReplaceOne(sa.context, filter, item, opts)
	if err != nil {
		return err
	}
	return nil
}

// FindAllContentItems finds all content items
func (sa *Adapter) FindAllContentItems() ([]model.ContentItemResponse, error) {
	filter := bson.D{}
	var result []model.ContentItemResponse
	err := sa.db.contentItems.Find(sa.context, filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateDataContentItem creates a data content item
func (sa *Adapter) CreateDataContentItem(item *model.DataContentItem) (*model.DataContentItem, error) {

	_, err := sa.db.dataContentItems.InsertOne(sa.context, &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// FindDataContentItem gets a data content item
func (sa *Adapter) FindDataContentItem(appID *string, orgID string, key string) (*model.DataContentItem, error) {

	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "key", Value: key}}

	var result *model.DataContentItem
	err := sa.db.dataContentItems.FindOne(sa.context, filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindDataContentItems gets multiple data content items
func (sa *Adapter) FindDataContentItems(appID *string, orgID string, category string) ([]*model.DataContentItem, error) {
	var filter bson.D
	if len(category) > 0 {
		filter = bson.D{primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
			primitive.E{Key: "category", Value: category}}
	} else {
		filter = bson.D{primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID}}
	}

	var result []*model.DataContentItem
	err := sa.db.dataContentItems.Find(sa.context, filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateDataContentItem updates a data content item
func (sa *Adapter) UpdateDataContentItem(appID *string, orgID string, item *model.DataContentItem) (*model.DataContentItem, error) {

	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "key", Value: item.Key}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "category", Value: item.Category},
			primitive.E{Key: "data", Value: item.Data},
			primitive.E{Key: "date_updated", Value: time.Now().UTC()},
		}},
	}
	_, err := sa.db.dataContentItems.UpdateOne(sa.context, filter, update, nil)
	if err != nil {
		log.Printf("error updating data content item: %s", err)
		return nil, err
	}

	return item, nil
}

// DeleteDataContentItem deletes a data content item
func (sa *Adapter) DeleteDataContentItem(appID *string, orgID string, key string) error {

	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "key", Value: key}}

	result, err := sa.db.dataContentItems.DeleteOne(sa.context, filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("result is nil for data content item with key %s", key)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return fmt.Errorf("error occured while deleting a data content item with key %s", key)
	}
	return nil
}

// CreateCategory created a new category
func (sa *Adapter) CreateCategory(item *model.Category) (*model.Category, error) {

	_, err := sa.db.categories.InsertOne(sa.context, &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// FindCategory fins a category
func (sa *Adapter) FindCategory(appID *string, orgID string, name string) (*model.Category, error) {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "name", Value: name}}

	var result *model.Category
	err := sa.db.categories.FindOne(sa.context, filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateCategory updates a  category
func (sa *Adapter) UpdateCategory(appID *string, orgID string, item *model.Category) (*model.Category, error) {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: item.ID}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "name", Value: item.Name},
			primitive.E{Key: "permissions", Value: item.Permissions},
			primitive.E{Key: "date_updated", Value: time.Now().UTC()},
		}},
	}
	_, err := sa.db.categories.UpdateOne(sa.context, filter, update, nil)
	if err != nil {
		log.Printf("error updating category: %s", err)
		return nil, err
	}

	return item, nil
}

// DeleteCategory deletes a category
func (sa *Adapter) DeleteCategory(appID *string, orgID string, name string) error {
	filter := bson.D{primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "name", Value: name}}

	result, err := sa.db.categories.DeleteOne(sa.context, filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("result is nil for cateogry with id %s", name)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return fmt.Errorf("error occured while deleting a category with id %s", name)
	}
	return nil
}

// StoreMultiTenancyData stores multi-tenancy to already exisiting data in the collections
func (sa *Adapter) StoreMultiTenancyData(appID string, orgID string) error {

	filter := bson.D{}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
		}},
	}
	//content items
	_, err := sa.db.contentItems.UpdateMany(sa.context, filter, update, nil)
	if err != nil {
		return err
	}
	//health locations
	_, err = sa.db.healthLocations.UpdateMany(sa.context, filter, update, nil)
	if err != nil {
		return err
	}
	//student guides
	_, err = sa.db.studentGuides.UpdateMany(sa.context, filter, update, nil)
	if err != nil {
		return err
	}

	return nil
}

// CreateMetaData creates meta_data object
func (sa *Adapter) CreateMetaData(key string, value map[string]interface{}) (*model.MetaData, error) {
	now := time.Now()
	id, _ := uuid.NewUUID()
	item := model.MetaData{
		ID:          id.String(),
		Key:         key,
		Value:       value,
		DateCreated: now,
	}

	_, err := sa.db.metaData.InsertOne(sa.context, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// FindMetaData find meta_data object
func (sa *Adapter) FindMetaData(key *string) (*model.MetaData, error) {
	filter := bson.D{primitive.E{Key: "key", Value: key}}

	var result *model.MetaData
	err := sa.db.metaData.FindOne(sa.context, filter, &result, nil)
	if err != nil {
		return nil, nil
	}
	return result, nil
}

// DeleteMetaData deletes meta_data object
func (sa *Adapter) DeleteMetaData(key string) error {
	filter := bson.D{primitive.E{Key: "key", Value: key}}

	result, err := sa.db.metaData.DeleteOne(sa.context, filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("result is nil for meta_data with key %s", key)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return fmt.Errorf("error occured while deleting a meta_data with key %s", key)
	}
	return nil
}

// UpdateMetaData updates a  metaData
func (sa *Adapter) UpdateMetaData(item *model.MetaData, value map[string]interface{}) (*model.MetaData, error) {
	filter := bson.D{
		primitive.E{Key: "key", Value: item.Key}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "value", Value: value},
			primitive.E{Key: "date_updated", Value: time.Now().UTC()},
		}},
	}
	_, err := sa.db.metaData.UpdateOne(sa.context, filter, update, nil)
	if err != nil {
		log.Printf("error updating category: %s", err)
		return nil, err
	}

	return item, nil
}

func (sa *Adapter) SyncDepartmentAttributes(appID *string, orgID string, contentItemID string) error {
	units, err := fetchUniversityUnits()
	if err != nil {
		return err
	}

	// 1) load attributes doc (the one you pasted)
	filter := bson.M{"_id": contentItemID, "app_id": appID, "org_id": orgID}
	var doc bson.M
	if err := sa.db.contentItems.FindOne(sa.context, filter, &doc, nil); err != nil {
		return err
	}

	// 2) load snapshot (separate doc)
	snapFilter := bson.M{"category": "uiuc_units_snapshot", "app_id": appID, "org_id": orgID}
	var snapDoc bson.M
	err = sa.db.contentItems.FindOne(sa.context, snapFilter, &snapDoc, nil)
	if err != nil {
		// first run: just save snapshot and exit (no renames yet)
		return sa.saveUnitsSnapshot(appID, orgID, units)
	}

	oldUnits := extractUnitsSnapshot(snapDoc)
	oldByID := unitsByID(oldUnits)
	newByID := unitsByID(units)

	// 3) detect renames by ID and apply alias updates
	for id, oldU := range oldByID {
		newU, ok := newByID[id]
		if !ok {
			continue
		}

		deptRenamed := oldU.Name != newU.Name
		colRenamed := oldU.CollegeName != newU.CollegeName

		if !deptRenamed && !colRenamed {
			continue
		}

		// rename department entry: {label:new, value:old} and keep requirements.college as OLD
		if deptRenamed {
			if err := sa.renameDepartmentValue(filter, oldU.CollegeName, oldU.Name, newU.Name, newU.CollegeName); err != nil {
				return err
			}
		}

		// rename college display: update group for all departments with requirements.college == oldCollege
		if colRenamed {
			if err := sa.renameCollegeGroupForDepartments(filter, oldU.CollegeName, newU.CollegeName); err != nil {
				return err
			}
		}
	}

	// 4) add new departments (IDs present in new but not in old)
	for id, newU := range newByID {
		if _, ok := oldByID[id]; ok {
			continue
		}
		if err := sa.addDepartmentValue(filter, newU); err != nil {
			return err
		}
	}

	// 5) update snapshot at end
	return sa.saveUnitsSnapshot(appID, orgID, units)
}

func fetchUniversityUnits() ([]model.UniversityUnit, error) {
	resp, err := http.Get("https://www.dmi.illinois.edu/ddd/mktextdirectory.asp")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("dmi status: %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// TSV reader
	r := csv.NewReader(bytes.NewReader(raw))
	r.Comma = '\t'
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("dmi: not enough rows")
	}

	header := records[0]
	colIndex := func(name string) (int, error) {
		for i, h := range header {
			if strings.TrimSpace(h) == name {
				return i, nil
			}
		}
		return -1, fmt.Errorf("missing column %q", name)
	}

	idxID, err := colIndex("Banner_Org")
	if err != nil {
		return nil, err
	}
	idxName, err := colIndex("deptFullName")
	if err != nil {
		return nil, err
	}
	idxCollege, err := colIndex("CollegeName")
	if err != nil {
		return nil, err
	}

	units := make([]model.UniversityUnit, 0, len(records)-1)
	for _, row := range records[1:] {
		maxIdx := idxID
		if idxName > maxIdx {
			maxIdx = idxName
		}
		if idxCollege > maxIdx {
			maxIdx = idxCollege
		}
		if len(row) <= maxIdx {
			continue
		}

		id := strings.TrimSpace(row[idxID])
		name := strings.TrimSpace(row[idxName])
		col := strings.TrimSpace(row[idxCollege])

		if id == "" || name == "" || col == "" {
			continue
		}

		units = append(units, model.UniversityUnit{
			ID: id, Name: name, CollegeName: col,
		})
	}
	return units, nil
}

func unitsByID(units []model.UniversityUnit) map[string]model.UniversityUnit {
	m := make(map[string]model.UniversityUnit, len(units))
	for _, u := range units {
		m[u.ID] = u
	}
	return m
}

func (sa *Adapter) renameDepartmentValue(
	filter bson.M,
	oldCollege string,
	oldLabel string,
	newLabel string,
	newCollege string,
) error {

	update := bson.M{
		"$set": bson.M{
			"data.attributes.$[attr].values.$[val].label": newLabel,
			"data.attributes.$[attr].values.$[val].value": oldLabel,
			"data.attributes.$[attr].values.$[val].group": newCollege, // display only
			"date_updated": time.Now().UTC(),
		},
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"attr.id": "department"},
			bson.M{
				"val.label":                oldLabel,
				"val.requirements.college": oldCollege,
			},
		},
	})

	_, err := sa.db.contentItems.UpdateOne(sa.context, filter, update, opts)
	return err
}
func (sa *Adapter) renameCollegeGroupForDepartments(filter bson.M, oldCollege, newCollege string) error {
	update := bson.M{
		"$set": bson.M{
			"data.attributes.$[attr].values.$[val].group": newCollege,
			"date_updated": time.Now().UTC(),
		},
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"attr.id": "department"},
			bson.M{"val.requirements.college": oldCollege},
		},
	})

	_, err := sa.db.contentItems.UpdateOne(sa.context, filter, update, opts)
	return err
}

func extractUnitsSnapshot(doc bson.M) []model.UniversityUnit {
	data, _ := doc["data"].(bson.M)
	if data == nil {
		return nil
	}
	rawUnits, _ := data["units"].(bson.A)
	if rawUnits == nil {
		return nil
	}

	out := make([]model.UniversityUnit, 0, len(rawUnits))
	for _, it := range rawUnits {
		m, _ := it.(bson.M)
		if m == nil {
			continue
		}
		out = append(out, model.UniversityUnit{
			ID:          fmt.Sprint(m["id"]),
			Name:        fmt.Sprint(m["name"]),
			CollegeName: fmt.Sprint(m["college"]),
		})
	}
	return out
}

func (sa *Adapter) saveUnitsSnapshot(appID *string, orgID string, units []model.UniversityUnit) error {
	dataUnits := make(bson.A, 0, len(units))
	for _, u := range units {
		dataUnits = append(dataUnits, bson.M{
			"id":      u.ID,
			"name":    u.Name,
			"college": u.CollegeName,
		})
	}

	update := bson.M{
		"$set": bson.M{
			"category":     "uiuc_units_snapshot",
			"app_id":       appID,
			"org_id":       orgID,
			"data":         bson.M{"units": dataUnits},
			"date_updated": time.Now().UTC(),
		},
		"$setOnInsert": bson.M{
			"date_created": time.Now().UTC(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := sa.db.contentItems.UpdateOne(sa.context,
		bson.M{"category": "uiuc_units_snapshot", "app_id": appID, "org_id": orgID},
		update, opts,
	)
	return err
}

func (sa *Adapter) addDepartmentValue(filter bson.M, u model.UniversityUnit) error {
	newVal := bson.M{
		"group": u.CollegeName,
		"label": u.Name,
		"requirements": bson.M{
			"college": u.CollegeName,
		},
	}

	update := bson.M{
		"$push": bson.M{
			"data.attributes.$[attr].values": newVal,
		},
		"$set": bson.M{
			"date_updated": time.Now().UTC(),
		},
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"attr.id": "department"},
		},
	})

	_, err := sa.db.contentItems.UpdateOne(sa.context, filter, update, opts)
	return err
}

func (sa *Adapter) abortTransaction(sessionContext mongo.SessionContext) {
	err := sessionContext.AbortTransaction(sessionContext)
	if err != nil {
		log.Printf("error aborting a transaction - %s", err)
	}
}

// NewStorageAdapter creates a new storage adapter instance
func NewStorageAdapter(mongoDBAuth string, mongoDBName string, mongoTimeout string, logger *logs.Logger) *Adapter {
	timeout, err := strconv.Atoi(mongoTimeout)
	if err != nil {
		log.Println("Set default timeout - 500")
		timeout = 500
	}
	timeoutMS := time.Millisecond * time.Duration(timeout)

	db := &database{mongoDBAuth: mongoDBAuth, mongoDBName: mongoDBName, mongoTimeout: timeoutMS, logger: logger}
	return &Adapter{db: db}
}

// Creates a new Adapter with provided context
func (sa *Adapter) withContext(context mongo.SessionContext) *Adapter {
	return &Adapter{db: sa.db, context: context}
}
