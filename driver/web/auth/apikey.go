package web

import (
	"fmt"
	"log"
	"net/http"
)

//APIKeysAuth entity
type APIKeysAuth struct {
	appKeys []string
}

// Check checks the request contains a valid ROKWIRE-API-KEY header
func (auth *APIKeysAuth) Check(r *http.Request) bool {
	apiKey := r.Header.Get("ROKWIRE-API-KEY")
	//check if there is api key in the header
	if len(apiKey) == 0 {
		//no key, so return 400
		log.Println(fmt.Sprintf("400 - Bad Request"))
		return false
	}

	//check if the api key is one of the listed
	appKeys := auth.appKeys
	exist := false
	for _, element := range appKeys {
		if element == apiKey {
			exist = true
			break
		}
	}
	if !exist {
		//not exist, so return 401
		log.Println(fmt.Sprintf("401 - Unauthorized for key %s", apiKey))
		return false
	}
	return true
}

// NewAPIKeysAuth creates new APIKeysAuth
func NewAPIKeysAuth(appKeys []string) *APIKeysAuth {
	auth := APIKeysAuth{appKeys}
	return &auth
}
