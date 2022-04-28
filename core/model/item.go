package model

import "time"

// ContentItem defines abstract data structure that would be used for any purpose
type ContentItem struct {
	ID          string      `json:"id" bson:"_id"`
	Category    string      `json:"category" bson:"category"`
	DateCreated time.Time   `json:"date_created" bson:"date_created"`
	DateUpdated time.Time   `json:"date_updated" bson:"date_updated"`
	Data        interface{} `json:"data" bson:"data"`
} // @name ContentItem
