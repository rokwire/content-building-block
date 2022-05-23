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
	"bytes"
	"content/core/model"
	"fmt"
	"image"
	"image/gif"
	jpeg "image/jpeg"
	"image/png"
	"strings"

	"github.com/nfnt/resize"
	"go.mongodb.org/mongo-driver/bson"
)

func (app *Application) getVersion() string {
	return app.version
}

// Student guides

func (app *Application) getStudentGuides(appID string, orgID string, ids []string) ([]bson.M, error) {
	items, err := app.storage.GetStudentGuides(appID, orgID, ids)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) getStudentGuide(appID string, orgID string, id string) (bson.M, error) {
	item, err := app.storage.GetStudentGuide(appID, orgID, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (app *Application) createStudentGuide(appID string, orgID string, item bson.M) (bson.M, error) {
	items, err := app.storage.CreateStudentGuide(appID, orgID, item)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) updateStudentGuide(appID string, orgID string, id string, item bson.M) (bson.M, error) {
	items, err := app.storage.UpdateStudentGuide(appID, orgID, id, item)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) deleteStudentGuide(appID string, orgID string, id string) error {
	err := app.storage.DeleteStudentGuide(appID, orgID, id)
	return err
}

// Health Locations

func (app *Application) getHealthLocations(appID string, orgID string, ids []string) ([]bson.M, error) {
	items, err := app.storage.GetHealthLocations(appID, orgID, ids)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) getHealthLocation(appID string, orgID string, id string) (bson.M, error) {
	item, err := app.storage.GetHealthLocation(appID, orgID, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (app *Application) createHealthLocation(appID string, orgID string, item bson.M) (bson.M, error) {
	items, err := app.storage.CreateHealthLocation(appID, orgID, item)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) updateHealthLocation(appID string, orgID string, id string, item bson.M) (bson.M, error) {
	items, err := app.storage.UpdateHealthLocation(appID, orgID, id, item)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) deleteHealthLocation(appID string, orgID string, id string) error {
	err := app.storage.DeleteHealthLocation(appID, orgID, id)
	return err
}

// Content Items

func (app *Application) getContentItemsCategories() ([]string, error) {
	return app.storage.GetContentItemsCategories()
}

func (app *Application) getContentItems(ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItemResponse, error) {
	return app.storage.GetContentItems(ids, categoryList, offset, limit, order)
}

func (app *Application) getContentItem(id string) (*model.ContentItemResponse, error) {
	return app.storage.GetContentItem(id)
}

func (app *Application) createContentItem(item *model.ContentItem) (*model.ContentItem, error) {
	return app.storage.CreateContentItem(item)
}

func (app *Application) updateContentItem(id string, item *model.ContentItem) (*model.ContentItem, error) {
	return app.storage.UpdateContentItem(id, item)
}

func (app *Application) deleteContentItem(id string) error {
	return app.storage.DeleteContentItem(id)
}

// Misc

func (app *Application) uploadImage(fileName string, filetype string, bytes []byte, path string, preferredFileName *string, spec model.ImageSpec) (*string, error) {

	err := app.tempStorageAdapter.Save(fileName, filetype, bytes)
	if err != nil {
		return nil, fmt.Errorf("Unable to save file: %s", err)
	}

	inputFileName := fileName
	var outputFileName string
	if strings.Contains(fileName, ".webp") {
		outputFileName = fileName
	} else {
		outputFileName = fmt.Sprintf("%s.%s", strings.Split(fileName, ".")[0], "webp") //get the file name without the extension
	}

	defer app.tempStorageAdapter.Delete(inputFileName)
	defer app.tempStorageAdapter.Delete(outputFileName)

	err = app.webpAdapter.Convert(inputFileName, outputFileName, spec)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert to webp file: %s", err)
	}

	convertedFile, err := app.tempStorageAdapter.Read(outputFileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to read webp file: %s", err)
	}

	url, err := app.awsAdapter.CreateImage(convertedFile, path, preferredFileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to upload to S3: %s", err)
	}

	if url != nil {
		return url, nil
	}

	return nil, nil
}

func (app *Application) getProfileImage(userID string, imageType string) ([]byte, error) {
	return app.awsAdapter.LoadProfileImage(fmt.Sprintf("profile-images/%s-%s.webp", userID, imageType))
}

func (app *Application) uploadProfileImage(userID string, filetype string, fileBytes []byte) error {
	var err error
	var defaultImage image.Image
	var mediumImage *image.Image
	var smallImage *image.Image
	defaultFileNameJpeg := fmt.Sprintf("%s-default.jpeg", userID)
	mediumFileNameJpeg := fmt.Sprintf("%s-medium.jpeg", userID)
	smallFileNameJpeg := fmt.Sprintf("%s-small.jpeg", userID)
	defaultFileNameWebp := fmt.Sprintf("%s-default", userID)
	mediumFileNameWebp := fmt.Sprintf("%s-medium", userID)
	smallFileNameWebp := fmt.Sprintf("%s-small", userID)

	defer app.tempStorageAdapter.Delete(defaultFileNameJpeg)
	defer app.tempStorageAdapter.Delete(mediumFileNameJpeg)
	defer app.tempStorageAdapter.Delete(smallFileNameJpeg)

	if filetype == "image/jpeg" {
		defaultImage, err = jpeg.Decode(bytes.NewReader(fileBytes))
	} else if filetype == "image/png" {
		defaultImage, err = png.Decode(bytes.NewReader(fileBytes))
	} else if filetype == "image/gif" {
		defaultImage, err = gif.Decode(bytes.NewReader(fileBytes))
	}
	if err != nil {
		return fmt.Errorf("unable to read profile image file: %s", err)
	}

	err = app.tempStorageAdapter.Save(defaultFileNameJpeg, filetype, fileBytes)
	if err != nil {
		return fmt.Errorf("Unable to save file: %s", defaultFileNameJpeg)
	}

	bounds := defaultImage.Bounds()
	if bounds.Dx() > 512 || bounds.Dy() > 512 {
		image := resize.Resize(512, 0, defaultImage, resize.Lanczos3)
		mediumImage = &image
	} else {
		mediumImage = &defaultImage
	}
	mediumImageBuff := bytes.NewBuffer([]byte{})

	// Generate medium profile photo
	if filetype == "image/jpeg" {
		err = jpeg.Encode(mediumImageBuff, *mediumImage, nil)
	} else if filetype == "image/png" {
		err = png.Encode(mediumImageBuff, *mediumImage)
	} else if filetype == "image/gif" {
		err = gif.Encode(mediumImageBuff, *mediumImage, nil)
	}
	if err != nil {
		return fmt.Errorf("Unable to save file: %s", mediumFileNameJpeg)
	}
	err = app.tempStorageAdapter.Save(mediumFileNameJpeg, filetype, mediumImageBuff.Bytes())
	if err != nil {
		return fmt.Errorf("Unable to save file: %s", defaultFileNameJpeg)
	}

	if bounds.Dx() > 256 || bounds.Dy() > 256 {
		image := resize.Resize(256, 0, defaultImage, resize.Lanczos3)
		smallImage = &image
	} else {
		smallImage = &defaultImage
	}

	// Generate small profile photo
	smallImageBuff := bytes.NewBuffer([]byte{})
	if filetype == "image/jpeg" {
		err = jpeg.Encode(smallImageBuff, *smallImage, nil)
	} else if filetype == "image/png" {
		err = png.Encode(smallImageBuff, *smallImage)
	} else if filetype == "image/gif" {
		err = gif.Encode(smallImageBuff, *smallImage, nil)
	}
	if err != nil {
		return fmt.Errorf("Unable to save file: %s", smallFileNameJpeg)
	}
	err = app.tempStorageAdapter.Save(smallFileNameJpeg, filetype, smallImageBuff.Bytes())
	if err != nil {
		return fmt.Errorf("Unable to save file: %s. Error: %s", smallFileNameJpeg, err)
	}

	_, err = app.uploadProfileImageToAws(defaultFileNameWebp, filetype, fileBytes, "profile-images/", &defaultFileNameWebp, model.ImageSpec{})
	if err != nil {
		return fmt.Errorf("Unable to upload file: %s. Error: %s", defaultFileNameWebp, err)
	}
	_, err = app.uploadProfileImageToAws(mediumFileNameWebp, filetype, mediumImageBuff.Bytes(), "profile-images/", &mediumFileNameWebp, model.ImageSpec{})
	if err != nil {
		return fmt.Errorf("Unable to upload file: %s. Error: %s", mediumFileNameWebp, err)
	}
	_, err = app.uploadProfileImageToAws(smallFileNameWebp, filetype, smallImageBuff.Bytes(), "profile-images/", &smallFileNameWebp, model.ImageSpec{})
	if err != nil {
		return fmt.Errorf("Unable to upload file: %s. Error: %s", smallFileNameWebp, err)
	}

	return nil
}

func (app *Application) deleteProfileImage(userID string) error {
	err := app.awsAdapter.DeleteProfileImage(fmt.Sprintf("profile-images/%s-default.webp", userID))
	if err != nil {
		return err
	}
	err = app.awsAdapter.DeleteProfileImage(fmt.Sprintf("profile-images/%s-medium.webp", userID))
	if err != nil {
		return err
	}
	err = app.awsAdapter.DeleteProfileImage(fmt.Sprintf("profile-images/%s-small.webp", userID))
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) uploadProfileImageToAws(fileName string, filetype string, bytes []byte, path string, preferredFileName *string, spec model.ImageSpec) (*string, error) {

	err := app.tempStorageAdapter.Save(fileName, filetype, bytes)
	if err != nil {
		return nil, fmt.Errorf("Unable to save file: %s", err)
	}

	inputFileName := fileName
	var outputFileName string
	if strings.Contains(fileName, ".webp") {
		outputFileName = fileName
	} else {
		outputFileName = fmt.Sprintf("%s.%s", strings.Split(fileName, ".")[0], "webp") //get the file name without the extension
	}

	defer app.tempStorageAdapter.Delete(inputFileName)
	defer app.tempStorageAdapter.Delete(outputFileName)

	err = app.webpAdapter.Convert(inputFileName, outputFileName, spec)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert to webp file: %s", err)
	}

	convertedFile, err := app.tempStorageAdapter.Read(outputFileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to read webp file: %s", err)
	}

	url, err := app.awsAdapter.CreateProfileImage(convertedFile, path, preferredFileName)
	if err != nil {
		return nil, fmt.Errorf("Unable to upload to S3: %s", err)
	}

	if url != nil {
		return url, nil
	}

	return nil, nil
}

func (app *Application) getTwitterPosts(userID string, twitterQueryParams string, force bool) (map[string]interface{}, error) {
	var err error
	posts := app.cacheAdapter.GetTwitterPosts(userID, twitterQueryParams)
	if posts == nil || force {
		app.cacheLock.Lock()
		posts = app.cacheAdapter.GetTwitterPosts(userID, twitterQueryParams)
		if posts == nil || force {
			if force {
				app.cacheAdapter.ClearTwitterCacheForUser(userID)
			}
			posts, err = app.twitterAdapter.GetTwitterPosts(userID, twitterQueryParams)
			if err == nil {
				app.cacheAdapter.SetTwitterPosts(userID, twitterQueryParams, posts)
			} else {
				fmt.Printf("error feeding twitter: %s", err)
			}
		}
		app.cacheLock.Unlock()
	}
	return posts, err
}
