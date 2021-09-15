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
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getStudentGuides(ids []string) ([]bson.M, error) {
	items, err := app.storage.GetStudentGuides(ids)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) getStudentGuide(id string) (bson.M, error) {
	item, err := app.storage.GetStudentGuide(id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (app *Application) createStudentGuide(item bson.M) (bson.M, error) {
	items, err := app.storage.CreateStudentGuide(item)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) updateStudentGuide(id string, item bson.M) (bson.M, error) {
	items, err := app.storage.UpdateStudentGuide(id, item)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) deleteStudentGuide(id string) error {
	err := app.storage.DeleteStudentGuide(id)
	return err
}

func (app *Application) uploadImage(fileName string, filetype string, bytes []byte, path string, spec model.ImageSpec) (bson.M, error) {

	err := app.tempStorageAdapter.Save(fileName, filetype, bytes)
	if err != nil {
		return nil, fmt.Errorf("Unable to save file: %s", err)
	}

	inputFileName := fileName
	outputFileName := fmt.Sprintf("%s.%s", strings.Split(fileName, ".")[0], "webp") //get the file name without the extension
	err = app.webpAdapter.Convert(inputFileName, outputFileName, spec)
	if err != nil {
		app.tempStorageAdapter.Delete(inputFileName)
		app.tempStorageAdapter.Delete(outputFileName)
		return nil, fmt.Errorf("Unable to convert to webp file: %s", err)
	}

	convertedFile, err := app.tempStorageAdapter.Read(outputFileName)
	if err != nil {
		app.tempStorageAdapter.Delete(inputFileName)
		app.tempStorageAdapter.Delete(outputFileName)
		return nil, fmt.Errorf("Unable to read webp file: %s", err)
	}

	url, err := app.awsAdapter.CreateImage(convertedFile, path)
	if err != nil {
		app.tempStorageAdapter.Delete(inputFileName)
		app.tempStorageAdapter.Delete(outputFileName)
		return nil, fmt.Errorf("Unable to upload to S3: %s", err)
	}

	app.tempStorageAdapter.Delete(inputFileName)
	app.tempStorageAdapter.Delete(outputFileName)

	if url != nil {
		return bson.M{"url": url}, nil
	}

	return nil, nil
}

func (app *Application) getTwitterPosts(count int, force bool) (map[string]interface{}, error) {
	var err error
	posts := app.cacheAdapter.GetTwitterPosts(count)
	if posts == nil || force {
		app.cacheLock.Lock()
		posts = app.cacheAdapter.GetTwitterPosts(count)
		if posts == nil || force {
			posts, err = app.twitterAdapter.GetTwitterPosts(count)
			if err == nil {
				app.cacheAdapter.SetTwitterPosts(count, posts)
			}
			app.cacheLock.Unlock()
		}
	}
	return posts, err
}
