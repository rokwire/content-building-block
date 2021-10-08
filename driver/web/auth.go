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
	web "content/driver/web/auth"
	"fmt"
	"log"
	"net/http"
)

// Auth handler
type Auth struct {
	apiKeysAuth    *web.APIKeysAuth
	shibbolethAuth *web.ShibbolethAuth
	coreAuth       *web.CoreAuth
}

func (auth *Auth) clientIDCheck(w http.ResponseWriter, r *http.Request) bool {
	clientID := r.Header.Get("APP")
	if len(clientID) == 0 {
		clientID = "edu.illinois.rokwire"
	}

	log.Println(fmt.Sprintf("400 - Bad Request"))
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
	return false
}

func (auth *Auth) apiKeyCheck(w http.ResponseWriter, r *http.Request) bool {
	return auth.apiKeysAuth.Check(r)
}

func (auth *Auth) shibbolethCheck(w http.ResponseWriter, r *http.Request) (bool, *model.ShibbolethToken) {
	return auth.shibbolethAuth.Check(r)
}

// NewAuth creates new auth handler
func NewAuth(app *core.Application, config model.Config) *Auth {
	apiKeysAuth := web.NewAPIKeysAuth(config.AppKeys)
	shibbolethAuth := web.NewShibbolethAuth(app, config)
	coreAuth := web.NewCoreAuth(app, config)

	auth := Auth{apiKeysAuth: apiKeysAuth, shibbolethAuth: shibbolethAuth, coreAuth: coreAuth}
	return &auth
}
