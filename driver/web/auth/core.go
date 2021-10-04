package web

import (
	"content/core"
	"content/core/model"
	"github.com/rokmetro/auth-library/authservice"
	"github.com/rokmetro/auth-library/tokenauth"
	"log"
	"net/http"
)

// CoreAuth implementation
type CoreAuth struct {
	app                *core.Application
	tokenAuth          *tokenauth.TokenAuth
	coreAuthPrivateKey *string
}

// NewCoreAuth creates new CoreAuth
func NewCoreAuth(app *core.Application, config model.Config) *CoreAuth {

	serviceLoader := authservice.NewRemoteServiceRegLoader(config.CoreServiceRegLoaderURL, []string{"core"})
	authService, err := authservice.NewAuthService("content", config.ContentServiceURL, serviceLoader)
	if err != nil {
		log.Fatalf("Error initializing auth service: %v", err)
	}
	tokenAuth, err := tokenauth.NewTokenAuth(true, authService, nil, nil)
	if err != nil {
		log.Fatalf("Error intitializing token auth: %v", err)
	}

	auth := CoreAuth{app: app, tokenAuth: tokenAuth, coreAuthPrivateKey: &config.CoreAuthPrivateKey}
	return &auth
}

// Check checks the request contains a valid Core access token
func (ca CoreAuth) Check(w http.ResponseWriter, r *http.Request) (bool, *tokenauth.Claims) {
	claims, err := ca.tokenAuth.CheckRequestTokens(r)
	if err != nil {
		log.Printf("error validate token: %s", err)
		return false, nil
	}

	if claims != nil {
		if claims.Valid() == nil {
			return true, claims
		}
	}

	return false, nil
}
