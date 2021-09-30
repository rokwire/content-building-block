package model

// Config the main config structure
type Config struct {
	AppKeys                 []string
	OidcProvider            string
	OidcAppClientID         string
	AdminAppClientID        string
	WebAppClientID          string
	PhoneAuthSecret         string
	AuthKeys                string
	AuthIssuer              string
	CoreAuthPrivateKey      string
	CoreServiceRegLoaderURL string
	ContentServiceURL       string
}
