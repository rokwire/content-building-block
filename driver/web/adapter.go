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

package web

import (
	"content/core"
	"content/core/model"
	"content/driver/web/rest"
	"content/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	httpSwagger "github.com/swaggo/http-swagger"
)

//Adapter entity
type Adapter struct {
	host          string
	port          string
	auth          *Auth
	authorization *casbin.Enforcer

	apisHandler      rest.ApisHandler
	adminApisHandler rest.AdminApisHandler

	app *core.Application
}

// @title Rokwire Content Building Block API
// @description Rokwire Content Building Block API Documentation.
// @version 1.1.10
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /content
// @schemes https

// @securityDefinitions.apikey RokwireAuth
// @in header
// @name ROKWIRE-API-KEY

// @securityDefinitions.apikey UserAuth
// @in header (add Bearer prefix to the Authorization value)
// @name Authorization

// @securityDefinitions.apikey AdminUserAuth
// @in header (add Bearer prefix to the Authorization value)
// @name Authorization

// @securityDefinitions.apikey AdminGroupAuth
// @in header
// @name GROUP

//Start starts the module
func (we Adapter) Start() {

	router := mux.NewRouter().StrictSlash(true)

	// handle apis
	contentRouter := router.PathPrefix("/content").Subrouter()
	contentRouter.PathPrefix("/doc/ui").Handler(we.serveDocUI())
	contentRouter.HandleFunc("/doc", we.serveDoc)
	contentRouter.HandleFunc("/version", we.wrapFunc(we.apisHandler.Version)).Methods("GET")

	contentRouter.HandleFunc("/profile_photo/{user-id}", we.coreAuthWrapFunc(we.apisHandler.GetProfilePhoto)).Methods("GET")
	contentRouter.HandleFunc("/profile_photo", we.coreAuthWrapFunc(we.apisHandler.GetUserProfilePhoto)).Methods("GET")
	contentRouter.HandleFunc("/profile_photo", we.coreAuthWrapFunc(we.apisHandler.StoreProfilePhoto)).Methods("POST")
	contentRouter.HandleFunc("/profile_photo", we.coreAuthWrapFunc(we.apisHandler.DeleteProfilePhoto)).Methods("DELETE")

	// handle student guide client apis
	contentRouter.HandleFunc("/student_guides", we.coreAuthWrapFunc(we.apisHandler.GetStudentGuides)).Methods("GET")
	contentRouter.HandleFunc("/student_guides/{id}", we.coreAuthWrapFunc(we.apisHandler.GetStudentGuide)).Methods("GET")
	contentRouter.HandleFunc("/health_locations", we.coreAuthWrapFunc(we.apisHandler.GetHealthLocations)).Methods("GET")
	contentRouter.HandleFunc("/health_locations/{id}", we.coreAuthWrapFunc(we.apisHandler.GetHealthLocation)).Methods("GET")
	contentRouter.HandleFunc("/content_items", we.coreAuthWrapFunc(we.apisHandler.GetContentItems)).Methods("GET")
	contentRouter.HandleFunc("/content_items/{id}", we.coreAuthWrapFunc(we.apisHandler.GetContentItem)).Methods("GET")
	contentRouter.HandleFunc("/content_item/categories", we.coreAuthWrapFunc(we.apisHandler.GetContentItemsCategories)).Methods("GET")
	contentRouter.HandleFunc("/image", we.coreAuthWrapFunc(we.apisHandler.UploadImage)).Methods("POST")
	contentRouter.HandleFunc("/twitter/users/{user_id}/tweets", we.coreAuthWrapFunc(we.apisHandler.GetTweeterPosts)).Methods("GET")

	// handle student guide admin apis
	adminSubRouter := contentRouter.PathPrefix("/admin").Subrouter()
	adminSubRouter.HandleFunc("/student_guides", we.coreAuthWrapFunc(we.adminApisHandler.GetStudentGuides)).Methods("GET")
	adminSubRouter.HandleFunc("/student_guides", we.coreAuthWrapFunc(we.adminApisHandler.CreateStudentGuide)).Methods("POST")
	adminSubRouter.HandleFunc("/student_guides/{id}", we.coreAuthWrapFunc(we.adminApisHandler.GetStudentGuide)).Methods("GET")
	adminSubRouter.HandleFunc("/student_guides/{id}", we.coreAuthWrapFunc(we.adminApisHandler.UpdateStudentGuide)).Methods("PUT")
	adminSubRouter.HandleFunc("/student_guides/{id}", we.coreAuthWrapFunc(we.adminApisHandler.DeleteStudentGuide)).Methods("DELETE")

	adminSubRouter.HandleFunc("/health_locations", we.coreAuthWrapFunc(we.adminApisHandler.GetHealthLocations)).Methods("GET")
	adminSubRouter.HandleFunc("/health_locations", we.coreAuthWrapFunc(we.adminApisHandler.CreateHealthLocation)).Methods("POST")
	adminSubRouter.HandleFunc("/health_location/{id}", we.coreAuthWrapFunc(we.adminApisHandler.GetHealthLocation)).Methods("GET")
	adminSubRouter.HandleFunc("/health_location/{id}", we.coreAuthWrapFunc(we.adminApisHandler.UpdateHealthLocation)).Methods("PUT")
	adminSubRouter.HandleFunc("/health_location/{id}", we.coreAuthWrapFunc(we.adminApisHandler.DeleteHealthLocation)).Methods("DELETE")

	adminSubRouter.HandleFunc("/content_items", we.coreAuthWrapFunc(we.adminApisHandler.GetContentItems)).Methods("GET")
	adminSubRouter.HandleFunc("/content_items", we.coreAuthWrapFunc(we.adminApisHandler.CreateContentItem)).Methods("POST")
	adminSubRouter.HandleFunc("/content_items/{id}", we.coreAuthWrapFunc(we.adminApisHandler.GetContentItem)).Methods("GET")
	adminSubRouter.HandleFunc("/content_items/{id}", we.coreAuthWrapFunc(we.adminApisHandler.UpdateContentItem)).Methods("PUT")
	adminSubRouter.HandleFunc("/content_items/{id}", we.coreAuthWrapFunc(we.adminApisHandler.DeleteContentItem)).Methods("DELETE")
	adminSubRouter.HandleFunc("/content_item/categories", we.coreAuthWrapFunc(we.adminApisHandler.GetContentItemsCategories)).Methods("GET")

	adminSubRouter.HandleFunc("/image", we.coreAuthWrapFunc(we.adminApisHandler.UploadImage)).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+we.port, router))
}

func (we Adapter) serveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("access-control-allow-origin", "*")
	http.ServeFile(w, r, "./docs/swagger.yaml")
}

func (we Adapter) serveDocUI() http.Handler {
	url := fmt.Sprintf("%s/content/doc", we.host)
	return httpSwagger.Handler(httpSwagger.URL(url))
}

func (we Adapter) wrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		handler(w, req)
	}
}

type coreAuthFunc = func(*tokenauth.Claims, http.ResponseWriter, *http.Request)

func (we Adapter) coreAuthWrapFunc(handler coreAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		obj := req.URL.Path // the resource that is going to be accessed.
		act := req.Method   // the operation that the user performs on the resource.

		coreAuth, claims := we.auth.coreAuth.Check(req)
		if coreAuth {
			permissions := strings.Split(claims.Permissions, ",")

			HasAccess := false
			for _, s := range permissions {
				HasAccess = we.authorization.Enforce(s, obj, act)
				if HasAccess {
					break
				}
			}
			if HasAccess {
				handler(claims, w, req)
				return
			}
			log.Printf("Access control error - Core Subject: %s is trying to apply %s operation for %s\n", claims.Subject, act, obj)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

// NewWebAdapter creates new WebAdapter instance
func NewWebAdapter(host string, port string, app *core.Application, config model.Config) Adapter {
	auth := NewAuth(app, config)
	authorization := casbin.NewEnforcer("driver/web/authorization_model.conf", "driver/web/authorization_policy.csv")

	apisHandler := rest.NewApisHandler(app)
	adminApisHandler := rest.NewAdminApisHandler(app)
	return Adapter{host: host, port: port, auth: auth, authorization: authorization, apisHandler: apisHandler, adminApisHandler: adminApisHandler, app: app}
}

// AppListener implements core.ApplicationListener interface
type AppListener struct {
	adapter *Adapter
}
