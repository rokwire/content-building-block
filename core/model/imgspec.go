package model

// ImageSpec wrapper for image convertor that holds all settings
type ImageSpec struct {
	Height  int `json:"height"`
	Width   int `json:"width"`
	Quality int `json:"quality"`
}
