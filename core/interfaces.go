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

package core

import (
	"content/core/model"
	"content/driven/storage"
	"io"

	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
	"go.mongodb.org/mongo-driver/bson"
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string
	GetStudentGuides(appID string, orgID string, ids []string) ([]bson.M, error)
	GetStudentGuide(appID string, orgID string, id string) (bson.M, error)
	CreateStudentGuide(appID string, orgID string, item bson.M) (bson.M, error)
	UpdateStudentGuide(appID string, orgID string, id string, item bson.M) (bson.M, error)
	DeleteStudentGuide(appID string, orgID string, id string) error

	GetHealthLocations(appID string, orgID string, ids []string) ([]bson.M, error)
	GetHealthLocation(appID string, orgID string, id string) (bson.M, error)
	CreateHealthLocation(appID string, orgID string, item bson.M) (bson.M, error)
	UpdateHealthLocation(appID string, orgID string, id string, item bson.M) (bson.M, error)
	DeleteHealthLocation(appID string, orgID string, id string) error

	//allApps says if the data is associated with the current app or it is for all the apps within the organization
	GetContentItemsCategories(allApps bool, appID string, orgID string) ([]string, error)
	GetContentItems(allApps bool, appID string, orgID string, ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItemResponse, error)
	GetContentItem(allApps bool, appID string, orgID string, id string) (*model.ContentItemResponse, error)
	CreateContentItem(allApps bool, appID string, orgID string, category string, data interface{}) (*model.ContentItem, error)
	UpdateContentItem(allApps bool, appID string, orgID string, id string, category string, data interface{}) (*model.ContentItem, error)
	UpdateContentItemData(allApps bool, appID string, orgID string, id string, category string, data interface{}) (*model.ContentItem, error)
	DeleteContentItem(allApps bool, appID string, orgID string, id string) error
	DeleteContentItemByCategory(allApps bool, appID string, orgID string, id string, category string) error

	UploadImage(imageBytes []byte, path string, spec model.ImageSpec) (*string, error)
	GetProfileImage(userID string, imageType string) ([]byte, error)
	UploadProfileImage(userID string, bytes []byte) error
	DeleteProfileImage(userID string) error

	GetTwitterPosts(userID string, twitterQueryParams string, force bool) (map[string]interface{}, error)

	CreateDataContentItem(claims *tokenauth.Claims, item *model.DataContentItem) (*model.DataContentItem, error)
	GetDataContentItem(claims *tokenauth.Claims, key string) (*model.DataContentItem, error)
	UpdateDataContentItem(claims *tokenauth.Claims, item *model.DataContentItem) (*model.DataContentItem, error)
	DeleteDataContentItem(claims *tokenauth.Claims, key string) error
	GetDataContentItems(claims *tokenauth.Claims, category string) ([]*model.DataContentItem, error)

	CreateCategory(claims *tokenauth.Claims, item *model.Category) (*model.Category, error)
	GetCategory(appID string, orgID string, id string) (*model.Category, error)
	UpdateCategory(appID string, orgID string, item *model.Category) (*model.Category, error)
	DeleteCategory(appID string, orgID string, id string) error

	UploadFileContentItem(file io.Reader, claims *tokenauth.Claims, fileName string, category string) (*string, error)
	GetFileContentItem(claims *tokenauth.Claims, fileName string, category string) ([]byte, error)
	DeleteFileContentItem(claims *tokenauth.Claims, fileName string, category string) error
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

// Student Guides

func (s *servicesImpl) GetStudentGuides(appID string, orgID string, ids []string) ([]bson.M, error) {
	return s.app.getStudentGuides(appID, orgID, ids)
}

func (s *servicesImpl) CreateStudentGuide(appID string, orgID string, item bson.M) (bson.M, error) {
	return s.app.createStudentGuide(appID, orgID, item)
}

func (s *servicesImpl) GetStudentGuide(appID string, orgID string, id string) (bson.M, error) {
	return s.app.getStudentGuide(appID, orgID, id)
}

func (s *servicesImpl) UpdateStudentGuide(appID string, orgID string, id string, item bson.M) (bson.M, error) {
	return s.app.updateStudentGuide(appID, orgID, id, item)
}

func (s *servicesImpl) DeleteStudentGuide(appID string, orgID string, id string) error {
	return s.app.deleteStudentGuide(appID, orgID, id)
}

// Health Locations

func (s *servicesImpl) GetHealthLocations(appID string, orgID string, ids []string) ([]bson.M, error) {
	return s.app.getHealthLocations(appID, orgID, ids)
}

func (s *servicesImpl) CreateHealthLocation(appID string, orgID string, item bson.M) (bson.M, error) {
	return s.app.createHealthLocation(appID, orgID, item)
}

func (s *servicesImpl) GetHealthLocation(appID string, orgID string, id string) (bson.M, error) {
	return s.app.getHealthLocation(appID, orgID, id)
}

func (s *servicesImpl) UpdateHealthLocation(appID string, orgID string, id string, item bson.M) (bson.M, error) {
	return s.app.updateHealthLocation(appID, orgID, id, item)
}

func (s *servicesImpl) DeleteHealthLocation(appID string, orgID string, id string) error {
	return s.app.deleteHealthLocation(appID, orgID, id)
}

// Content Items

func (s *servicesImpl) GetContentItemsCategories(allApps bool, appID string, orgID string) ([]string, error) {
	return s.app.getContentItemsCategories(allApps, appID, orgID)
}

func (s *servicesImpl) GetContentItems(allApps bool, appID string, orgID string, ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItemResponse, error) {
	return s.app.getContentItems(allApps, appID, orgID, ids, categoryList, offset, limit, order)
}

func (s *servicesImpl) GetContentItem(allApps bool, appID string, orgID string, id string) (*model.ContentItemResponse, error) {
	return s.app.getContentItem(allApps, appID, orgID, id)
}

func (s *servicesImpl) CreateContentItem(allApps bool, appID string, orgID string, category string, data interface{}) (*model.ContentItem, error) {
	return s.app.createContentItem(allApps, appID, orgID, category, data)
}

func (s *servicesImpl) UpdateContentItem(allApps bool, appID string, orgID string, id string, category string, data interface{}) (*model.ContentItem, error) {
	return s.app.updateContentItem(allApps, appID, orgID, id, category, data)
}

func (s *servicesImpl) UpdateContentItemData(allApps bool, appID string, orgID string, id string, category string, data interface{}) (*model.ContentItem, error) {
	return s.app.updateContentItemData(allApps, appID, orgID, id, category, data)
}

func (s *servicesImpl) DeleteContentItem(allApps bool, appID string, orgID string, id string) error {
	return s.app.deleteContentItem(allApps, appID, orgID, id)
}

func (s *servicesImpl) DeleteContentItemByCategory(allApps bool, appID string, orgID string, id string, category string) error {
	return s.app.deleteContentItemByCategory(allApps, appID, orgID, id, category)
}

// Misc

func (s *servicesImpl) GetProfileImage(userID string, imageType string) ([]byte, error) {
	return s.app.getProfileImage(userID, imageType)
}

func (s *servicesImpl) UploadImage(imageBytes []byte, path string, spec model.ImageSpec) (*string, error) {
	return s.app.uploadImage(imageBytes, path, nil, spec)
}

func (s *servicesImpl) UploadProfileImage(userID string, fileBytes []byte) error {
	return s.app.uploadProfileImage(userID, fileBytes)
}

func (s *servicesImpl) DeleteProfileImage(userID string) error {
	return s.app.deleteProfileImage(userID)
}

func (s *servicesImpl) GetTwitterPosts(userID string, twitterQueryParams string, force bool) (map[string]interface{}, error) {
	return s.app.getTwitterPosts(userID, twitterQueryParams, force)
}

func (s *servicesImpl) CreateDataContentItem(claims *tokenauth.Claims, item *model.DataContentItem) (*model.DataContentItem, error) {
	return s.app.createDataContentItem(claims, item)
}

func (s *servicesImpl) GetDataContentItem(claims *tokenauth.Claims, key string) (*model.DataContentItem, error) {
	return s.app.getDataContentItem(claims, key)
}

func (s *servicesImpl) GetDataContentItems(claims *tokenauth.Claims, category string) ([]*model.DataContentItem, error) {
	return s.app.getDataContentItems(claims, category)
}

func (s *servicesImpl) UpdateDataContentItem(claims *tokenauth.Claims, item *model.DataContentItem) (*model.DataContentItem, error) {
	return s.app.updateDataContentItem(claims, item)
}

func (s *servicesImpl) DeleteDataContentItem(claims *tokenauth.Claims, key string) error {
	return s.app.deleteDataContentItem(claims, key)
}

func (s *servicesImpl) CreateCategory(claims *tokenauth.Claims, item *model.Category) (*model.Category, error) {
	return s.app.createCategory(claims, item)
}

func (s *servicesImpl) GetCategory(appID string, orgID string, id string) (*model.Category, error) {
	return s.app.getCategory(&appID, orgID, id)
}

func (s *servicesImpl) UpdateCategory(appID string, orgID string, item *model.Category) (*model.Category, error) {
	return s.app.updateCategory(&appID, orgID, item)
}

func (s *servicesImpl) DeleteCategory(appID string, orgID string, id string) error {
	return s.app.deleteCategory(&appID, orgID, id)
}

func (s *servicesImpl) UploadFileContentItem(file io.Reader, claims *tokenauth.Claims, fileName string, category string) (*string, error) {
	return s.app.uploadFileContentItem(file, claims, fileName, category)
}

func (s *servicesImpl) GetFileContentItem(claims *tokenauth.Claims, fileName string, category string) ([]byte, error) {
	return s.app.getFileContentItem(claims, fileName, category)
}

func (s *servicesImpl) DeleteFileContentItem(claims *tokenauth.Claims, fileName string, category string) error {
	return s.app.deleteFileContentItem(claims, fileName, category)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	PerformTransaction(func(context storage.TransactionContext) error) error

	GetStudentGuides(appID string, orgID string, ids []string) ([]bson.M, error)
	GetStudentGuide(appID string, orgID string, id string) (bson.M, error)
	CreateStudentGuide(appID string, orgID string, item bson.M) (bson.M, error)
	UpdateStudentGuide(appID string, orgID string, id string, item bson.M) (bson.M, error)
	DeleteStudentGuide(appID string, orgID string, id string) error

	GetHealthLocations(appID string, orgID string, ids []string) ([]bson.M, error)
	GetHealthLocation(appID string, orgID string, id string) (bson.M, error)
	CreateHealthLocation(appID string, orgID string, item bson.M) (bson.M, error)
	UpdateHealthLocation(appID string, orgID string, id string, item bson.M) (bson.M, error)
	DeleteHealthLocation(appID string, orgID string, id string) error

	GetContentItemsCategories(appID *string, orgID string) ([]string, error)
	FindContentItems(appID *string, orgID string, ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItem, error)
	GetContentItems(appID *string, orgID string, ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItemResponse, error)
	GetContentItem(appID *string, orgID string, id string) (*model.ContentItemResponse, error)
	CreateContentItem(item model.ContentItem) (*model.ContentItem, error)
	UpdateContentItem(appID *string, orgID string, id string, category string, data interface{}) (*model.ContentItem, error)
	DeleteContentItem(appID *string, orgID string, id string) error
	SaveContentItem(item model.ContentItem) error

	//Used for multi-tenancy for already exisiting data.
	//To be removed when this is applied to all environments.
	FindAllContentItems(context storage.TransactionContext) ([]model.ContentItemResponse, error)
	StoreMultiTenancyData(context storage.TransactionContext, appID string, orgID string) error
	///

	CreateDataContentItem(item *model.DataContentItem) (*model.DataContentItem, error)
	FindDataContentItem(appID *string, orgID string, key string, context storage.TransactionContext) (*model.DataContentItem, error)
	UpdateDataContentItem(appID *string, orgID string, item *model.DataContentItem) (*model.DataContentItem, error)
	DeleteDataContentItem(appID *string, orgID string, key string, context storage.TransactionContext) error
	FindDataContentItems(appID *string, orgID string, key string) ([]*model.DataContentItem, error)

	CreateCategory(item *model.Category) (*model.Category, error)
	FindCategory(appID *string, orgID string, id string, context storage.TransactionContext) (*model.Category, error)
	UpdateCategory(appID *string, orgID string, item *model.Category) (*model.Category, error)
	DeleteCategory(appID *string, orgID string, key string) error
}
