package rest

import (
	"content/core"
	"content/core/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//AdminApisHandler handles the rest Admin APIs implementation
type AdminApisHandler struct {
	app *core.Application
}

// GetStudentGuides Retrieves  all student guides
// @Description Retrieves  all student guides
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

	if resData == nil {
		resData = []bson.M{}
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

// CreateStudentGuide Creates a student guide item
// @Description Creates a student guide item
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

// DeleteStudentGuide Deletes a student guide item with the specified id
// @Description Deletes a student guide item with the specified id
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

// GetHealthLocations Retrieves  all health locations
// @Description Retrieves  all health locations
// @Param ids query string false "Coma separated IDs of the desired records"
// @Tags Admin
// @ID AdminGetHealthLocations
// @Accept json
// @Success 200
// @Security AdminUserAuth
// @Router /admin/health_locations [get]
func (h AdminApisHandler) GetHealthLocations(w http.ResponseWriter, r *http.Request) {

	IDs := []string{}
	IDskeys, ok := r.URL.Query()["ids"]
	if ok && len(IDskeys[0]) > 0 {
		extIDs := IDskeys[0]
		IDs = strings.Split(extIDs, ",")
	}

	resData, err := h.app.Services.GetHealthLocations(IDs)
	if err != nil {
		log.Printf("Error on health location items by id - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []bson.M{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal all health locations")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetHealthLocation Retrieves a health location by id
// @Description Retrieves a health location by id
// @Tags Admin
// @ID AdminGetHealthLocation
// @Accept json
// @Produce json
// @Success 200
// @Security AdminUserAuth
// @Router /admin/health_locations/{id} [get]
func (h AdminApisHandler) GetHealthLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	locationID := vars["id"]

	resData, err := h.app.Services.GetHealthLocation(locationID)
	if err != nil {
		log.Printf("Error on getting health location id - %s\n %s", locationID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal the health location")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateHealthLocation Updates a health location with the specified id
// @Description Updates a health location with the specified id
// @Tags Admin
// @ID AdminUpdateHealthLocation
// @Accept json
// @Produce json
// @Success 200
// @Security AdminUserAuth
// @Router /admin/health_locations/{id} [put]
func (h AdminApisHandler) UpdateHealthLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	locationID := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a health location - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item bson.M
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create health location request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateHealthLocation(locationID, item)
	if err != nil {
		log.Printf("Error on updating health location with id - %s\n %s", locationID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal the updated health location")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateHealthLocation Create a new health location
// @Description Create a new health location
// @Tags Admin
// @ID AdminCreateHealthLocation
// @Accept json
// @Success 200
// @Security AdminUserAuth
// @Router /admin/health_locations [post]
func (h AdminApisHandler) CreateHealthLocation(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a health location - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item bson.M
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create health location request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateHealthLocation(item)
	if err != nil {
		log.Printf("Error on creating health location: %s\n", err)
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

// DeleteHealthLocation Deletes a health location with the specified id
// @Description Deletes a health location with the specified id
// @Tags Admin
// @ID AdminDeleteHealthLocation
// @Success 200
// @Security AdminUserAuth
// @Router /admin/health_location/{id} [delete]
func (h AdminApisHandler) DeleteHealthLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	locationID := vars["id"]

	err := h.app.Services.DeleteHealthLocation(locationID)
	if err != nil {
		log.Printf("Error on deleting health location with id - %s\n %s", locationID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// uploadImageResponse wrapper
type uploadImageResponse struct {
	URL string `json:"url"`
} // @name uploadImageResponse

// UploadImage Uploads an image to AWS S3
// @Description Uploads an image to AWS S3
// @Tags Admin
// @ID AdminUploadImage
// @Param path body string true "path - path within the S3 bucket"
// @Param width body string false "width - width of the image to resize. If width and height are missing - then the new image will use the original size"
// @Param height body string false "height - height of the image to resize. If width and height are missing - then the new image will use the original size"
// @Param quality body string false "quality - quality of the image. Default: 90"
// @Param fileName body string false "fileName - the uploaded file name"
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} uploadImageResponse
// @Security AdminUserAuth
// @Router /admin/image [post]
func (h AdminApisHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	//validate the image type
	path := r.PostFormValue("path")
	if len(path) <= 0 {
		log.Print("Missing image path\n")
		http.Error(w, "missing 'path' form param", http.StatusBadRequest)
		return
	}

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
	file, fileHeader, err := r.FormFile("fileName")
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
	fileName := fileHeader.Filename
	url, err := h.app.Services.UploadImage(fileName, filetype, fileBytes, path, imgSpec)
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

type getContentItemsRequestBody struct {
	IDs        []string `json:"ids,omitempty"`        // List of IDs for the filter. Optional and may be null or missing.
	Categories []string `json:"categories,omitempty"` // List of Categories for the filter. Optional and may be null or missing.
} // @name getContentItemsRequestBody

// GetContentItems Retrieves  all content items. <b> The data element could be either a primitive or nested json or array.</b>
// @Description Retrieves  all content items.<b> The data element could be either a primitive or nested json or array.</b>
// @Tags Admin
// @ID AdminGetContentItems
// @Param offset query string false "offset"
// @Param limit query string false "limit - limit the result"
// @Param order query string false "order - Possible values: asc, desc. Default: desc"
// @Param data body getContentItemsRequestBody false "Optional - body json of the all items ids that need to be filtered. NOTE: Bad/broken json will be interpreted as an empty filter and the request will be proceeded further."
// @Accept json
// @Success 200 {array} model.ContentItem
// @Security AdminUserAuth
// @Router /admin/content_items [get]
func (h AdminApisHandler) GetContentItems(w http.ResponseWriter, r *http.Request) {

	var offset *int64
	offsets, ok := r.URL.Query()["offset"]
	if ok && len(offsets[0]) > 0 {
		val, err := strconv.ParseInt(offsets[0], 0, 64)
		if err == nil {
			offset = &val
		}
	}

	var limit *int64
	limits, ok := r.URL.Query()["limit"]
	if ok && len(limits[0]) > 0 {
		val, err := strconv.ParseInt(limits[0], 0, 64)
		if err == nil {
			limit = &val
		}
	}

	var order *string
	orders, ok := r.URL.Query()["order"]
	if ok && len(orders[0]) > 0 {
		order = &orders[0]
	}

	var body getContentItemsRequestBody
	bodyData, _ := ioutil.ReadAll(r.Body)
	if len(bodyData) > 0 {
		bodyErr := json.Unmarshal(bodyData, &body)
		if bodyErr != nil {
			log.Printf("Warning: bad getContentItemsRequestBody request: %s", bodyErr)
		}
	}

	resData, err := h.app.Services.GetContentItems(body.IDs, body.Categories, offset, limit, order)
	if err != nil {
		log.Printf("Error on cgetting content items - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.ContentItemResponse{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal all content items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetContentItem Retrieves a content item by id. <b> The data element could be either a primitive or nested json or array.</b>
// @Description Retrieves a content item by id. <b> The data element could be either a primitive or nested json or array.</b>
// @Tags Admin
// @ID AdminGetContentItem
// @Accept json
// @Produce json
// @Success 200 {object} model.ContentItem
// @Security AdminUserAuth
// @Router /admin/content_items/{id} [get]
func (h AdminApisHandler) GetContentItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetContentItem(id)
	if err != nil {
		log.Printf("Error on getting content item id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal the content item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateContentItem Updates a content item with the specified id. <b> The data element could be either a primitive or nested json or array.</b>
// @Description Updates a content item with the specified id. <b> The data element could be either a primitive or nested json or array.</b>
// @Tags Admin
// @ID AdminUpdateContentItem
// @Accept json
// @Produce json
// @Success 200 {object} model.ContentItem
// @Security AdminUserAuth
// @Router /admin/content_items/{id} [put]
func (h AdminApisHandler) UpdateContentItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a content item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.ContentItem
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create content item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.ID != id {
		log.Printf("Inconsistent attempt to update: path id is not equal to json ID")
		http.Error(w, "Inconsistent attempt to update: path id is not equal to json ID", http.StatusBadRequest)
		return
	}

	if len(item.Category) == 0 {
		log.Printf("Unable to update content item: Missing category")
		http.Error(w, "Unable to create content item: Missing category", http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateContentItem(id, &item)
	if err != nil {
		log.Printf("Error on updating content item with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal the updated content item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// createContentItemRequestBody Expected body while creating a new content item
type createContentItemRequestBody struct {
	Category string      `json:"category" bson:"category"`
	Data     interface{} `json:"data" bson:"data"`
} // @name createContentItemRequestBody

// CreateContentItem creates a new content item. <b> The data element could be either a primitive or nested json or array.</b>
// @Description Creates a new content item. <b> The data element could be either a primitive or nested json or array.</b>
// @Tags Admin
// @ID AdminCreateContentItem
// @Accept json
// @Success 200 {object} createContentItemRequestBody
// @Security AdminUserAuth
// @Router /admin/content_items [post]
func (h AdminApisHandler) CreateContentItem(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a content item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item createContentItemRequestBody
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create content item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(item.Category) == 0 {
		log.Printf("Unable to create content item: Missing category")
		http.Error(w, "Unable to create content item: Missing category", http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateContentItem(&model.ContentItem{
		Category: item.Category,
		Data:     item.Data,
	})
	if err != nil {
		log.Printf("Error on creating content item: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Println("Error on marshal the new content item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteContentItem Deletes a content item with the specified id
// @Description Deletes a content item with the specified id
// @Tags Admin
// @ID AdminDeleteContentItem
// @Success 200
// @Security AdminUserAuth
// @Router /admin/content_items/{id} [delete]
func (h AdminApisHandler) DeleteContentItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guideID := vars["id"]

	err := h.app.Services.DeleteContentItem(guideID)
	if err != nil {
		log.Printf("Error on deleting content item with id - %s\n %s", guideID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// GetContentItemsCategories Retrieves  all content item categories that have in the database
// @Description Retrieves  all content item categories that have in the database
// @Tags Admin
// @ID AdminGetContentItemsCategories
// @Success 200
// @Security AdminUserAuth
// @Router /admin/content_item/categories [get]
func (h AdminApisHandler) GetContentItemsCategories(w http.ResponseWriter, r *http.Request) {
	resData, err := h.app.Services.GetContentItemsCategories()
	if err != nil {
		log.Printf("Error on cgetting content items - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []string{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal all content items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
