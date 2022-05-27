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

package web

import (
	"content/core"
	"content/core/model"
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"gopkg.in/ericchiang/go-oidc.v2"
)

type tokenData struct {
	UIuceduUIN        *string   `json:"uiucedu_uin"`
	Sub               *string   `json:"sub"`
	Audience          *string   `json:"aud"`
	Email             *string   `json:"email"`
	UIuceduIsMemberOf *[]string `json:"uiucedu_is_member_of"`
}

func (d *tokenData) HasClientID(clientIDs []string) bool {
	if d.Audience != nil && len(clientIDs) > 0 {
		for _, clientID := range clientIDs {
			if strings.EqualFold(*d.Audience, clientID) {
				return true
			}
		}
	}
	return false
}

// ShibbolethAuth entity
type ShibbolethAuth struct {
	app           *core.Application
	clientIDs     []string
	tokenVerifier *oidc.IDTokenVerifier
}

// Check checks the request contains a valid OIDC token from Shibboleth
func (auth *ShibbolethAuth) Check(r *http.Request) (bool, *model.ShibbolethToken) {
	//1. Get the token from the request
	rawIDToken, err := auth.getIDToken(r)
	if err != nil {
		return false, nil
	}

	//3. Validate the token
	idToken, err := auth.verify(*rawIDToken)
	if err != nil {
		log.Printf("error validating token - %s\n", err)
		return false, nil
	}

	//4. Get the user data from the token
	var tokenData tokenData
	if err := idToken.Claims(&tokenData); err != nil {
		log.Printf("error getting user data from token - %s\n", err)
		return false, nil
	}

	if !tokenData.HasClientID(auth.clientIDs) {
		log.Printf("error - Aud (%s) is not permitted %s\n", *tokenData.Audience, err)
		return false, nil
	}

	// we must have UIuceduUIN
	if tokenData.UIuceduUIN == nil {
		log.Printf("error - missing uiuceuin data in the token - %s\n", err)
		return false, nil
	}

	shibboAuth := &model.ShibbolethToken{Uin: *tokenData.UIuceduUIN, Email: *tokenData.Email,
		IsMemberOf: tokenData.UIuceduIsMemberOf}

	return true, shibboAuth
}

//gets the token from the request - as cookie or as Authorization header.
//returns the id token and its type - mobile or web. If the token is taken by the cookie it is web otherwise it is mobile
func (auth *ShibbolethAuth) getIDToken(r *http.Request) (*string, error) {
	//1. Check if there is a cookie
	cookie, err := r.Cookie("rwa-at-data")
	if err == nil && cookie != nil && len(cookie.Value) > 0 {
		//there is a cookie
		return &cookie.Value, nil
	}

	//2. Check if there is a token in the Authorization header
	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) <= 0 {
		return nil, errors.New("error getting Authorization header")
	}
	splitAuthorization := strings.Fields(authorizationHeader)
	if len(splitAuthorization) != 2 {
		return nil, errors.New("error processing the Authorization header")
	}
	// expected - Bearer 1234
	if splitAuthorization[0] != "Bearer" {
		return nil, errors.New("error processing the Authorization header")
	}
	rawIDToken := splitAuthorization[1]
	return &rawIDToken, nil
}

func (auth *ShibbolethAuth) verify(rawIDToken string) (*oidc.IDToken, error) {
	log.Println("ShibbolethToken -> token")
	return auth.tokenVerifier.Verify(context.Background(), rawIDToken)
}

// NewShibbolethAuth creates ShibbolethAuth instance
func NewShibbolethAuth(app *core.Application, config model.Config) *ShibbolethAuth {
	provider, err := oidc.NewProvider(context.Background(), config.OidcProvider)
	if err != nil {
		log.Fatalln(err)
	}

	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})

	return &ShibbolethAuth{
		app:           app,
		clientIDs:     config.OidcClientIDs,
		tokenVerifier: verifier,
	}
}
