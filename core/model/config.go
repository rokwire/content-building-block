package model

// Config the main config structure
type Config struct {
	OidcProvider      string
	OidcClientIDs     []string
	CoreBBHost        string
	ContentServiceURL string
}
