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
	"content/core/model"
	"go.mongodb.org/mongo-driver/bson"
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string
	GetStudentGuides(ids []string) ([]bson.M, error)
	GetStudentGuide(id string) (bson.M, error)
	CreateStudentGuide(item bson.M) (bson.M, error)
	UpdateStudentGuide(id string, item bson.M) (bson.M, error)
	DeleteStudentGuide(id string) error

	GetHealthLocations(ids []string) ([]bson.M, error)
	GetHealthLocation(id string) (bson.M, error)
	CreateHealthLocation(item bson.M) (bson.M, error)
	UpdateHealthLocation(id string, item bson.M) (bson.M, error)
	DeleteHealthLocation(id string) error

	GetContentItemsCategories() ([]string, error)
	GetContentItems(ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItem, error)
	GetContentItem(id string) (*model.ContentItem, error)
	CreateContentItem(item *model.ContentItem) (*model.ContentItem, error)
	UpdateContentItem(id string, item *model.ContentItem) (*model.ContentItem, error)
	DeleteContentItem(id string) error

	UploadImage(fileName string, filetype string, bytes []byte, path string, spec model.ImageSpec) (*string, error)
	GetProfileImage(userID string, imageType string) ([]byte, error)
	UploadProfileImage(userID string, filetype string, bytes []byte) error
	DeleteProfileImage(userID string) error

	GetTwitterPosts(userID string, twitterQueryParams string, force bool) (map[string]interface{}, error)
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

// Student Guides

func (s *servicesImpl) GetStudentGuides(ids []string) ([]bson.M, error) {
	return s.app.getStudentGuides(ids)
}

func (s *servicesImpl) CreateStudentGuide(item bson.M) (bson.M, error) {
	return s.app.createStudentGuide(item)
}

func (s *servicesImpl) GetStudentGuide(id string) (bson.M, error) {
	return s.app.getStudentGuide(id)
}

func (s *servicesImpl) UpdateStudentGuide(id string, item bson.M) (bson.M, error) {
	return s.app.updateStudentGuide(id, item)
}

func (s *servicesImpl) DeleteStudentGuide(id string) error {
	return s.app.deleteStudentGuide(id)
}

// Health Locations

func (s *servicesImpl) GetHealthLocations(ids []string) ([]bson.M, error) {
	return s.app.getHealthLocations(ids)
}

func (s *servicesImpl) CreateHealthLocation(item bson.M) (bson.M, error) {
	return s.app.createHealthLocation(item)
}

func (s *servicesImpl) GetHealthLocation(id string) (bson.M, error) {
	return s.app.getHealthLocation(id)
}

func (s *servicesImpl) UpdateHealthLocation(id string, item bson.M) (bson.M, error) {
	return s.app.updateHealthLocation(id, item)
}

func (s *servicesImpl) DeleteHealthLocation(id string) error {
	return s.app.deleteHealthLocation(id)
}

// Content Items

func (s *servicesImpl) GetContentItemsCategories() ([]string, error) {
	return s.app.getContentItemsCategories()
}

func (s *servicesImpl) GetContentItems(ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItem, error) {
	return s.app.getContentItems(ids, categoryList, offset, limit, order)
}

func (s *servicesImpl) GetContentItem(id string) (*model.ContentItem, error) {
	return s.app.getContentItem(id)
}

func (s *servicesImpl) CreateContentItem(item *model.ContentItem) (*model.ContentItem, error) {
	return s.app.createContentItem(item)
}

func (s *servicesImpl) UpdateContentItem(id string, item *model.ContentItem) (*model.ContentItem, error) {
	return s.app.updateContentItem(id, item)
}

func (s *servicesImpl) DeleteContentItem(id string) error {
	return s.app.deleteContentItem(id)
}

// Misc

func (s *servicesImpl) GetProfileImage(userID string, imageType string) ([]byte, error) {
	return s.app.getProfileImage(userID, imageType)
}

func (s *servicesImpl) UploadImage(fileName string, filetype string, bytes []byte, path string, spec model.ImageSpec) (*string, error) {
	return s.app.uploadImage(fileName, filetype, bytes, path, nil, spec)
}

func (s *servicesImpl) UploadProfileImage(userID string, filetype string, fileBytes []byte) error {
	return s.app.uploadProfileImage(userID, filetype, fileBytes)
}

func (s *servicesImpl) DeleteProfileImage(userID string) error {
	return s.app.deleteProfileImage(userID)
}

func (s *servicesImpl) GetTwitterPosts(userID string, twitterQueryParams string, force bool) (map[string]interface{}, error) {
	return s.app.getTwitterPosts(userID, twitterQueryParams, force)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetStudentGuides(ids []string) ([]bson.M, error)
	GetStudentGuide(id string) (bson.M, error)
	CreateStudentGuide(item bson.M) (bson.M, error)
	UpdateStudentGuide(id string, item bson.M) (bson.M, error)
	DeleteStudentGuide(id string) error

	GetHealthLocations(ids []string) ([]bson.M, error)
	GetHealthLocation(id string) (bson.M, error)
	CreateHealthLocation(item bson.M) (bson.M, error)
	UpdateHealthLocation(id string, item bson.M) (bson.M, error)
	DeleteHealthLocation(id string) error

	GetContentItemsCategories() ([]string, error)
	GetContentItems(ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItem, error)
	GetContentItem(id string) (*model.ContentItem, error)
	CreateContentItem(item *model.ContentItem) (*model.ContentItem, error)
	UpdateContentItem(id string, item *model.ContentItem) (*model.ContentItem, error)
	DeleteContentItem(id string) error
}
