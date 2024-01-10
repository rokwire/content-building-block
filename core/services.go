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
	"bytes"
	"content/core/model"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // Allow image.Decode to detect GIFs
	_ "image/jpeg" // Allow image.Decode to detect JPEGs
	_ "image/png"  // Allow image.Decode to detect PNGs
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
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

func (app *Application) getContentItemsCategories(allApps bool, appID string, orgID string) ([]string, error) {
	//logic
	var appIDParam *string
	if !allApps {
		appIDParam = &appID //associated with current app
	}
	return app.storage.GetContentItemsCategories(appIDParam, orgID)
}

func (app *Application) getContentItems(allApps bool, appID string, orgID string, ids []string, categoryList []string, offset *int64, limit *int64, order *string) ([]model.ContentItemResponse, error) {
	//logic
	var appIDParam *string
	if !allApps {
		appIDParam = &appID //associated with current app
	}
	return app.storage.GetContentItems(appIDParam, orgID, ids, categoryList, offset, limit, order)
}

func (app *Application) getContentItem(allApps bool, appID string, orgID string, id string) (*model.ContentItemResponse, error) {
	//logic
	var appIDParam *string
	if !allApps {
		appIDParam = &appID //associated with current app
	}
	return app.storage.GetContentItem(appIDParam, orgID, id)
}

func (app *Application) createContentItem(allApps bool, appID string, orgID string, category string, data interface{}) (*model.ContentItem, error) {
	//logic
	var appIDParam *string
	if !allApps {
		appIDParam = &appID //associated with current app
	}
	cItem := model.ContentItem{ID: uuid.NewString(), Category: category, DateCreated: time.Now().UTC(),
		Data: data, OrgID: orgID, AppID: appIDParam}
	return app.storage.CreateContentItem(cItem)
}

func (app *Application) updateContentItem(allApps bool, appID string, orgID string, id string, category string, data interface{}) (*model.ContentItem, error) {
	//logic
	var appIDParam *string
	if !allApps {
		appIDParam = &appID //associated with current app
	}

	//update
	item, err := app.storage.UpdateContentItem(appIDParam, orgID, id, category, data)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (app *Application) updateContentItemData(allApps bool, appID string, orgID string, id string, category string, data interface{}) (*model.ContentItem, error) {
	//logic
	var appIDParam *string
	if !allApps {
		appIDParam = &appID //associated with current app
	}

	//find the item
	items, err := app.storage.FindContentItems(appIDParam, orgID, []string{id}, []string{category}, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(items) != 1 {
		return nil, errors.New("not found")
	}
	item := items[0]

	//update the data
	item.Data = data
	now := time.Now()
	item.DateUpdated = &now

	//save it
	err = app.storage.SaveContentItem(item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (app *Application) deleteContentItem(allApps bool, appID string, orgID string, id string) error {
	//logic
	var appIDParam *string
	if !allApps {
		appIDParam = &appID //associated with current app
	}
	return app.storage.DeleteContentItem(appIDParam, orgID, id)
}

func (app *Application) deleteContentItemByCategory(allApps bool, appID string, orgID string, id string, category string) error {
	//logic
	var appIDParam *string
	if !allApps {
		appIDParam = &appID //associated with current app
	}

	//find the item
	items, err := app.storage.FindContentItems(appIDParam, orgID, []string{id}, []string{category}, nil, nil, nil)
	if err != nil {
		return err
	}
	if len(items) != 1 {
		return errors.New("not found")
	}

	//delete it
	err = app.storage.DeleteContentItem(appIDParam, orgID, id)
	if err != nil {
		return err
	}

	return nil
}

// Misc

func (app *Application) uploadImage(imageBytes []byte, path string, preferredFileName *string, spec model.ImageSpec) (*string, error) {
	image, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("Error decoding image: %s", err)
	}

	if spec.Height > 0 || spec.Width > 0 {
		image = resize.Resize(uint(spec.Width), uint(spec.Height), image, resize.Lanczos3)
	}

	quality := 75
	if spec.Quality > 0 {
		quality = spec.Quality
	}

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(quality))
	if err != nil {
		return nil, fmt.Errorf("Error creating webp encoder options: %s", err)
	}

	var output bytes.Buffer
	if err := webp.Encode(&output, image, options); err != nil {
		return nil, fmt.Errorf("Error encoding webp: %s", err)
	}

	url, err := app.awsAdapter.CreateImage(&output, path, preferredFileName)
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

func (app *Application) uploadProfileImage(userID string, imageBytes []byte) error {
	var mediumImage image.Image
	var smallImage image.Image

	defaultFileNameWebp := fmt.Sprintf("%s-default", userID)
	mediumFileNameWebp := fmt.Sprintf("%s-medium", userID)
	smallFileNameWebp := fmt.Sprintf("%s-small", userID)

	defaultImage, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return fmt.Errorf("Error decoding image: %s", err)
	}

	bounds := defaultImage.Bounds()
	if bounds.Dx() > 512 || bounds.Dy() > 512 {
		image := resize.Resize(512, 0, defaultImage, resize.Lanczos3)
		mediumImage = image
	} else {
		mediumImage = defaultImage
	}

	if bounds.Dx() > 256 || bounds.Dy() > 256 {
		image := resize.Resize(256, 0, defaultImage, resize.Lanczos3)
		smallImage = image
	} else {
		smallImage = defaultImage
	}

	_, err = app.uploadProfileImageToAws(defaultImage, defaultFileNameWebp, "profile-images/", model.ImageSpec{})
	if err != nil {
		return fmt.Errorf("Unable to upload de file: %s. Error: %s", defaultFileNameWebp, err)
	}
	_, err = app.uploadProfileImageToAws(mediumImage, mediumFileNameWebp, "profile-images/", model.ImageSpec{})
	if err != nil {
		return fmt.Errorf("Unable to upload file: %s. Error: %s", mediumFileNameWebp, err)
	}
	_, err = app.uploadProfileImageToAws(smallImage, smallFileNameWebp, "profile-images/", model.ImageSpec{})
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

func (app *Application) uploadProfileImageToAws(image image.Image, filename string, path string, spec model.ImageSpec) (*string, error) {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return nil, fmt.Errorf("Error creating webp encoder options: %s", err)
	}

	var output bytes.Buffer
	if err := webp.Encode(&output, image, options); err != nil {
		return nil, fmt.Errorf("Error encoding webp: %s", err)
	}

	url, err := app.awsAdapter.CreateProfileImage(&output, path, &filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to upload to S3: %s", err)
	}

	if url != nil {
		return url, nil
	}

	return nil, nil
}

func (app *Application) uploadVoiceRecord(userID string, bytes []byte) error {
	//TODO
	return errors.New("not implemented")
}

func (app *Application) getVoiceRecord(userID string) ([]byte, error) {
	//TODO
	return nil, errors.New("not implemented")
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
