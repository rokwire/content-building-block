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

	//delete data timer
	dailyRenameTimer *time.Timer
	timerDone        chan bool
}

func newDepartmentRenameLogic(logger logs.Logger, storage interfaces.Storage) departmentRenameLogic {
	return departmentRenameLogic{logger: logger, storage: storage}
}

func (d departmentRenameLogic) start() {
	go d.setupTimerForRename()
}

func (d departmentRenameLogic) setupTimerForRename() {

	// cancel if active
	if d.dailyRenameTimer != nil {
		d.logger.Info("setupTimerForDepartmentRename -> there is active timer, so cancel it")

		// signal abort first, then stop timer
		d.timerDone <- true
		d.dailyRenameTimer.Stop()
	}

	// wait until it is the correct moment from the day
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		d.logger.Errorf("Error getting location:%s\n", err.Error())
	}

	now := time.Now().In(location)
	d.logger.Infof(
		"setupTimerForDepartmentRename -> now - hours:%d minutes:%d seconds:%d\n",
		now.Hour(), now.Minute(), now.Second(),
	)

	nowSecondsInDay := 60*60*now.Hour() + 60*now.Minute() + now.Second()
	desiredMoment := 10800 // 3 AM

	var durationInSeconds int
	d.logger.Infof(
		"setupTimerForDepartmentRename -> nowSecondsInDay:%d desiredMoment:%d\n",
		nowSecondsInDay, desiredMoment,
	)

	if nowSecondsInDay <= desiredMoment {
		d.logger.Infof("setupTimerForDepartmentRename -> not processed today, so the first process will be today")
		durationInSeconds = desiredMoment - nowSecondsInDay
	} else {
		d.logger.Infof("setupTimerForDepartmentRename -> already processed today, so the first process will be tomorrow")
		leftToday := 86400 - nowSecondsInDay
		durationInSeconds = leftToday + desiredMoment // time left today + desired moment tomorrow
	}
	//duration := time.Second * time.Duration(3)
	duration := time.Second * time.Duration(durationInSeconds)
	d.logger.Infof("setupTimerForDepartmentRename -> first call after %s", duration)

	d.dailyRenameTimer = time.NewTimer(duration)
	select {
	case <-d.dailyRenameTimer.C:
		d.logger.Info("setupTimerForDepartmentRename -> rename timer expired")
		d.dailyRenameTimer = nil

		d.process()
	case <-d.timerDone:
		// timer aborted
		d.logger.Info("setupTimerForDepartmentRename -> rename timer aborted")
		d.dailyRenameTimer = nil
	}
}

func (d departmentRenameLogic) process() {
	d.logger.Info("Deleting data process")

	//process work
	d.processRename()

	//generate new processing after 24 hours
	duration := time.Hour * 24
	d.logger.Infof("Rename data process -> next call after %s", duration)
	d.dailyRenameTimer = time.NewTimer(duration)
	select {
	case <-d.dailyRenameTimer.C:
		d.logger.Info("Rename data process -> timer expired")
		d.dailyRenameTimer = nil

		d.process()
	case <-d.timerDone:
		// timer aborted
		d.logger.Info("Rename data process -> timer aborted")
		d.dailyRenameTimer = nil
	}
}

func (d departmentRenameLogic) processRename() {
}
