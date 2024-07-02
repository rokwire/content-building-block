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

package rest

import (
	"content/core"
	"content/core/model"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
)

// TPsApisHandler handles the rest TPs APIs implementation
type TPsApisHandler struct {
	app *core.Application
}

// UploadImage Uploads an image to the Rokwire system
func (h TPsApisHandler) UploadImage(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	path := "tps-images"

	heightParam := intPostValueFromString(r.PostFormValue("height"))
	widthParam := intPostValueFromString(r.PostFormValue("width"))
	qualityParam := intPostValueFromString(r.PostFormValue("quality"))
	imgSpec := model.ImageSpec{Height: heightParam, Width: widthParam, Quality: qualityParam}

	// validate file size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Print("File is too big\n")
		http.Error(w, "File is too big", http.StatusBadRequest)
		return
	}

	// parse and validate file and post parameters
	file, _, err := r.FormFile("fileName")
	if err != nil {
		log.Print("Invalid file\n")
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Print("Invalid file\n")
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	// check file type, detectcontenttype only needs the first 512 bytes
	filetype := http.DetectContentType(fileBytes)
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "image/webp":
		break
	default:
		log.Print("Invalid file type\n")
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// pass the file to be processed by the use case handler
	url, err := h.app.Services.UploadImage(fileBytes, path, imgSpec)
	if err != nil {
		log.Printf("Error converting image: %s\n", err)
		http.Error(w, "Error converting image", http.StatusInternalServerError)
		return
	}

	jsonData := map[string]string{"url": *url}
	jsonBynaryData, err := json.Marshal(jsonData)
	if err != nil {
		log.Println("Error on marshal s3 location data")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBynaryData)
}
