package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	client := &http.Client{}
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
