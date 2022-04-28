package model

// Config the main config structure
type Config struct {
	AppKeys            []string
	OidcProvider       string
	OidcClientIDs      []string
	PhoneAuthSecret    string
	AuthKeys           string
	AuthIssuer         string
	CoreAuthPrivateKey string
	CoreBBHost         string
	ContentServiceURL  string
}
