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
	"log"
	"net/http"

	"github.com/rokwire/core-auth-library-go/authorization"
	"github.com/rokwire/core-auth-library-go/authservice"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"github.com/rokwire/logging-library-go/errors"
	"github.com/rokwire/logging-library-go/logs"
	"github.com/rokwire/logging-library-go/logutils"
)

//Authorization is an interface for auth types
type Authorization interface {
	check(req *http.Request) (int, *tokenauth.Claims, error)
	start()
}

// Auth handler
type Auth struct {
	coreAuth *CoreAuth
}

// NewAuth creates new auth handler
func NewAuth(app *core.Application, config model.Config) *Auth {
	coreAuth := NewCoreAuth(app, config)

	auth := Auth{coreAuth: coreAuth}
	return &auth
}

// CoreAuth implementation
type CoreAuth struct {
	app       *core.Application
	tokenAuth *tokenauth.TokenAuth

	permissionsAuth *PermissionsAuth
	userAuth        *UserAuth
}

// NewCoreAuth creates new CoreAuth
func NewCoreAuth(app *core.Application, config model.Config) *CoreAuth {

	remoteConfig := authservice.RemoteAuthDataLoaderConfig{
		AuthServicesHost: config.CoreBBHost,
	}

	serviceLoader, err := authservice.NewRemoteAuthDataLoader(remoteConfig, []string{"core"}, logs.NewLogger("content", &logs.LoggerOpts{}))
	authService, err := authservice.NewAuthService("content", config.ContentServiceURL, serviceLoader)
	if err != nil {
		log.Fatalf("Error initializing auth service: %v", err)
	}

	adminPermissionAuth := authorization.NewCasbinStringAuthorization("driver/web/authorization_policy.csv")
	tokenAuth, err := tokenauth.NewTokenAuth(true, authService, adminPermissionAuth, nil)
	if err != nil {
		log.Fatalf("Error intitializing token auth: %v", err)
	}
	permissionsAuth := newPermissionsAuth(tokenAuth)
	usersAuth := newUserAuth(tokenAuth)

	auth := CoreAuth{app: app, tokenAuth: tokenAuth, permissionsAuth: permissionsAuth, userAuth: usersAuth}
	return &auth
}

//PermissionsAuth entity
//This enforces that the user has permissions matching the policy
type PermissionsAuth struct {
	tokenAuth *tokenauth.TokenAuth
}

func (a *PermissionsAuth) start() {}

func (a *PermissionsAuth) check(req *http.Request) (int, *tokenauth.Claims, error) {
	claims, err := a.tokenAuth.CheckRequestTokens(req)
	if err != nil {
		return http.StatusUnauthorized, nil, errors.WrapErrorAction("typeCheckServicesAuthRequestToken", logutils.TypeToken, nil, err)
	}

	if err == nil && claims != nil {
		err = a.tokenAuth.AuthorizeRequestPermissions(claims, req)
		if err != nil {
			return http.StatusForbidden, nil, errors.WrapErrorAction("check permissions", logutils.TypeRequest, nil, err)
		}
	}

	return http.StatusOK, claims, err
}

func newPermissionsAuth(tokenAuth *tokenauth.TokenAuth) *PermissionsAuth {
	permissionsAuth := PermissionsAuth{tokenAuth: tokenAuth}
	return &permissionsAuth
}

//UserAuth entity
// This enforces that the user is not anonymous
type UserAuth struct {
	tokenAuth *tokenauth.TokenAuth
}

func (a *UserAuth) start() {}

func (a *UserAuth) check(req *http.Request) (int, *tokenauth.Claims, error) {
	claims, err := a.tokenAuth.CheckRequestTokens(req)
	if err != nil {
		return http.StatusUnauthorized, nil, errors.WrapErrorAction("typeCheckServicesAuthRequestToken", logutils.TypeToken, nil, err)
	}

	if err == nil && claims != nil {
		if claims.Anonymous {
			return http.StatusForbidden, nil, errors.New("token must not be anonymous")
		}
	}

	return http.StatusOK, claims, err
}

func newUserAuth(tokenAuth *tokenauth.TokenAuth) *UserAuth {
	userAuth := UserAuth{tokenAuth: tokenAuth}
	return &userAuth
}
