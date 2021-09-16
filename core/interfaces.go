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

	UploadImage(fileName string, filetype string, bytes []byte, path string, spec model.ImageSpec) (bson.M, error)
	GetTwitterPosts(userID string, twitterQueryParams string, force bool) (map[string]interface{}, error)
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

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

func (s *servicesImpl) UploadImage(fileName string, filetype string, bytes []byte, path string, spec model.ImageSpec) (bson.M, error) {
	return s.app.uploadImage(fileName, filetype, bytes, path, spec)
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
}
