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
	"go.mongodb.org/mongo-driver/bson"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getAllStudentGuides() ([]bson.M, error) {
	items, err := app.storage.GetAllStudentGuides()
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
