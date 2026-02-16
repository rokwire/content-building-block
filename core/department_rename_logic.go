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
package core

import (
	"content/core/interfaces"
	"time"

	"github.com/rokwire/rokwire-building-block-sdk-go/utils/logging/logs"
)

type departmentRenameLogic struct {
	logger  logs.Logger
	storage interfaces.Storage

	appID         *string
	orgID         string
	contentItemID string

	timer *time.Timer
	done  chan bool
}

func newDepartmentRenameLogic(logger logs.Logger, storage interfaces.Storage, appID *string, orgID string, contentItemID string) departmentRenameLogic {

	return departmentRenameLogic{
		logger:        logger,
		storage:       storage,
		appID:         appID,
		orgID:         orgID,
		contentItemID: contentItemID,
		done:          make(chan bool),
	}
}

func (d departmentRenameLogic) start() {
	go d.setupTimerForRename()
}

func (d departmentRenameLogic) setupTimerForRename() {

	location, _ := time.LoadLocation("America/Chicago")
	now := time.Now().In(location)

	nowSeconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
	desired := 3 * 3600 // 3 AM

	var wait int
	if nowSeconds <= desired {
		wait = desired - nowSeconds
	} else {
		wait = (86400 - nowSeconds) + desired
	}

	d.timer = time.NewTimer(time.Duration(wait) * time.Second)

	select {
	case <-d.timer.C:
		d.process()
	case <-d.done:
		d.timer.Stop()
	}
}

func (d departmentRenameLogic) process() {

	d.logger.Info("Department sync started")

	err := d.storage.SyncDepartmentAttributes(d.appID, d.orgID, d.contentItemID)
	if err != nil {
		d.logger.Errorf("Department sync error: %v", err)
	}

	d.timer = time.NewTimer(24 * time.Hour)

	select {
	case <-d.timer.C:
		d.process()
	case <-d.done:
		d.timer.Stop()
	}
}
