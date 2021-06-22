package rest

import (
	"content/core"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//AdminApisHandler handles the rest Admin APIs implementation
type AdminApisHandler struct {
	app *core.Application
}

// GetStudentGuides retrieves  all items
// @Description Retrieves  all items
// @Param ids query string false "Coma separated IDs of the desired records"
// @Tags Admin
// @ID AdminGetStudentGuides
// @Accept json
// @Success 200
// @Security AdminUserAuth
// @Router /admin/student_guides [get]
func (h AdminApisHandler) GetStudentGuides(w http.ResponseWriter, r *http.Request) {

	IDs := []string{}
	IDskeys, ok := r.URL.Query()["ids"]
	if ok && len(IDskeys[0]) > 0 {
		extIDs := IDskeys[0]
		IDs = strings.Split(extIDs, ",")
	}

	resData, err := h.app.Services.GetStudentGuides(IDs)
	if err != nil {
		log.Printf("Error on getting guide items by id - %s\n", err)
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
// @Tags Admin
// @ID AdminGetStudentGuide
// @Accept json
// @Produce json
// @Success 200
// @Security AdminUserAuth
// @Router /admin/student_guides/{id} [get]
func (h AdminApisHandler) GetStudentGuide(w http.ResponseWriter, r *http.Request) {
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
// @Tags Admin
// @ID AdminUpdateStudentGuide
// @Accept json
// @Produce json
// @Success 200
// @Security AdminUserAuth
// @Router /admin/student_guides/{id} [put]
func (h AdminApisHandler) UpdateStudentGuide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guideID := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a student guide - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item bson.M
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create student guide request data - %s\n", err.Error())
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
		log.Println("Error on marshal the updated student guide")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateStudentGuide retrieves  all items
// @Description Retrieves  all items
// @Tags Admin
// @ID AdminCreateStudentGuide
// @Accept json
// @Success 200
// @Security AdminUserAuth
// @Router /admin/student_guides [post]
func (h AdminApisHandler) CreateStudentGuide(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a student guide - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item bson.M
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create student guide request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateStudentGuide(item)
	if err != nil {
		log.Printf("Error on creating student guide: %s\n", err)
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

// DeleteStudentGuide Deletes a student guide with the specified id
// @Description Deletes a student guide with the specified id
// @Tags Admin
// @ID AdminDeleteStudentGuide
// @Success 200
// @Security AdminUserAuth
// @Router /admin/student_guides/{id} [delete]
func (h AdminApisHandler) DeleteStudentGuide(w http.ResponseWriter, r *http.Request) {
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
