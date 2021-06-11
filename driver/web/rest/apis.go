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

package rest

import (
	"content/core"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
)

//ApisHandler handles the rest APIs implementation
type ApisHandler struct {
	app *core.Application
}

//Version gives the service version
// @Description Gives the service version.
// @ID Version
// @Produce plain
// @Success 200
// @Router /version [get]
func (h ApisHandler) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.app.Services.GetVersion()))
}

// CreateStudentGuide retrieves  all items
// @Description Retrieves  all items
// @Tags StudentGuides
// @ID CreateStudentGuide
// @Accept json
// @Success 200
// @Security AdminUserAuth
// @Router /student_guides [post]
func (h ApisHandler) CreateStudentGuide(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a ctest - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item bson.M
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create ctest request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateStudentGuide(item)
	if err != nil {
		log.Printf("Error on creatint item: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Println("Error on marshal the new item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetAllStudentGuides retrieves  all items
// @Description Retrieves  all items
// @Tags StudentGuides
// @ID getAllStudentGuides
// @Accept json
// @Success 200
// @Security RokwireAuth AdminUserAuth
// @Router /student_guides [get]
func (h ApisHandler) GetAllStudentGuides(w http.ResponseWriter, r *http.Request) {
	resData, err := h.app.Services.GetAllStudentGuides()
	if err != nil {
		log.Printf("Error on getting track items by external id - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal all student guides")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetStudentGuide Retrieves a student guide by id
// @Description Retrieves  all items
// @Tags StudentGuides
// @ID GetStudentGuide
// @Accept json
// @Produce json
// @Success 200
// @Security RokwireAuth AdminUserAuth
// @Router /student_guides/{id} [get]
func (h ApisHandler) GetStudentGuide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guideID := vars["id"]

	resData, err := h.app.Services.GetStudentGuide(guideID)
	if err != nil {
		log.Printf("Error on getting student guide id - %s\n %s", guideID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal the student guide")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateStudentGuide Updates a student guide with the specified id
// @Description Updates a student guide with the specified id
// @Tags StudentGuides
// @ID UpdateStudentGuide
// @Accept json
// @Produce json
// @Success 200
// @Security AdminUserAuth
// @Router /student_guides/{id} [put]
func (h ApisHandler) UpdateStudentGuide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guideID := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a ctest - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item bson.M
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create ctest request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateStudentGuide(guideID, item)
	if err != nil {
		log.Printf("Error on updating student guide with id - %s\n %s", guideID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal the new item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteStudentGuide Deletes a student guide with the specified id
// @Description Deletes a student guide with the specified id
// @Tags StudentGuides
// @ID DeleteStudentGuide
// @Success 200
// @Security AdminUserAuth
// @Router /student_guides/{id} [delete]
func (h ApisHandler) DeleteStudentGuide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guideID := vars["id"]

	err := h.app.Services.DeleteStudentGuide(guideID)
	if err != nil {
		log.Printf("Error on deleting student guide with id - %s\n %s", guideID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

//NewApisHandler creates new rest Handler instance
func NewApisHandler(app *core.Application) ApisHandler {
	return ApisHandler{app: app}
}
