/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package model

import "time"

//ShibbolethAuth represents shibboleth auth entity
type ShibbolethAuth struct {
	Uin        string    `json:"uiucedu_uin" bson:"uiucedu_uin"`
	Email      string    `json:"email" bson:"email"`
	IsMemberOf *[]string `json:"uiucedu_is_member_of" bson:"uiucedu_is_member_of"`
}

//User represents user entity
type User struct {
	ID          string     `json:"id" bson:"_id"`
	ExternalID  string     `json:"external_id" bson:"external_id"`
	Email       string     `json:"email" bson:"email"`
	IsMemberOf  *[]string  `json:"is_member_of" bson:"is_member_of"`
	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`

	ClientID string `bson:"client_id"`
} // @name User
