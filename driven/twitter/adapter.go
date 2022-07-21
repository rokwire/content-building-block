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

package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Adapter struct
type Adapter struct {
	twitterFeedURL     string
	twitterAccessToken string
}

// NewTwitterAdapter crates new instance
func NewTwitterAdapter(twitterFeedURL string, twitterAccessToken string) *Adapter {
	return &Adapter{
		twitterFeedURL:     twitterFeedURL,
		twitterAccessToken: twitterAccessToken,
	}
}

// GetTwitterPosts converts an image
func (a *Adapter) GetTwitterPosts(userID string, twitterQueryParams string) (map[string]interface{}, error) {
	url := fmt.Sprintf(a.twitterFeedURL, userID)
	url += fmt.Sprintf("?%s", twitterQueryParams)

	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error creating Twitter request - %s", err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.twitterAccessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error loading Twitter data - %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("error with Twitter response code - %d", resp.StatusCode)
		return nil, fmt.Errorf("error with Twitter response code != 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading the body data for the loading Twitter data request - %s", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("error converting data for the loading Twitter data request - %s", err)
		return nil, err
	}

	return result, nil
}
