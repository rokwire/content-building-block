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

package model

import "time"

// ContentItemResponse is a workaround due to problem with data json & bson encode and decode with abstract type
type ContentItemResponse = map[string]interface{}

// ContentItem defines abstract data structure that would be used for any purpose
type ContentItem struct {
	ID          string      `json:"id" bson:"_id"`
	Category    string      `json:"category" bson:"category"`
	DateCreated time.Time   `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time  `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
	Data        interface{} `json:"data" bson:"data"` // could be eigther a primitive or nested json or array
	OrgID       string      `json:"org_id" bson:"org_id"`
	AppID       *string     `json:"app_id" bson:"app_id"`
} // @name ContentItem
