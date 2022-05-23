/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TapdBug struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ID           uint64 `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	EpicKey      string
	Title        string `json:"name" gorm:"type:varchar(255)"`
	Description  string
	WorkspaceID  uint64          `json:"workspace_id,string"`
	Created      *helper.CSTTime `json:"created"`
	Modified     *helper.CSTTime `json:"modified" gorm:"index"`
	Status       string          `json:"status"`
	Cc           string          `json:"cc"`
	Begin        *helper.CSTTime `json:"begin"`
	Due          *helper.CSTTime `json:"due"`
	Priority     string          `json:"priority"`
	IterationID  uint64          `json:"iteration_id,string"`
	Source       string          `json:"source"`
	Module       string          `json:"module"`
	ReleaseID    uint64          `json:"release_id,string"`
	CreatedFrom  string          `json:"created_from"`
	Feature      string          `json:"feature"`
	common.NoPKModel

	Severity         string          `json:"severity"`
	Reporter         string          `json:"reporter"`
	Resolved         *helper.CSTTime `json:"resolved"`
	Closed           *helper.CSTTime `json:"closed"`
	Lastmodify       string          `json:"lastmodify"`
	Auditer          string          `json:"auditer"`
	De               string          `json:"De" gorm:"comment:developer;type:varchar(255)"`
	Fixer            string          `json:"fixer"`
	VersionTest      string          `json:"version_test"`
	VersionReport    string          `json:"version_report"`
	VersionClose     string          `json:"version_close"`
	VersionFix       string          `json:"version_fix"`
	BaselineFind     string          `json:"baseline_find"`
	BaselineJoin     string          `json:"baseline_join"`
	BaselineClose    string          `json:"baseline_close"`
	BaselineTest     string          `json:"baseline_test"`
	Sourcephase      string          `json:"sourcephase"`
	Te               string          `json:"te"`
	CurrentOwner     string          `json:"current_owner"`
	Resolution       string          `json:"resolution"`
	Originphase      string          `json:"originphase"`
	Confirmer        string          `json:"confirmer"`
	Participator     string          `json:"participator"`
	Closer           string          `json:"closer"`
	Platform         string          `json:"platform"`
	Os               string          `json:"os"`
	Testtype         string          `json:"testtype"`
	Testphase        string          `json:"testphase"`
	Frequency        string          `json:"frequency"`
	RegressionNumber string          `json:"regression_number"`
	Flows            string          `json:"flows"`
	Testmode         string          `json:"testmode"`
	IssueID          uint64          `json:"issue_id,string"`
	VerifyTime       *helper.CSTTime `json:"verify_time"`
	RejectTime       *helper.CSTTime `json:"reject_time"`
	ReopenTime       *helper.CSTTime `json:"reopen_time"`
	AuditTime        *helper.CSTTime `json:"audit_time"`
	SuspendTime      *helper.CSTTime `json:"suspend_time"`
	Deadline         *helper.CSTTime `json:"deadline"`
	InProgressTime   *helper.CSTTime `json:"in_progress_time"`
	AssignedTime     *helper.CSTTime `json:"assigned_time"`
	TemplateID       uint64          `json:"template_id,string"`
	StoryID          uint64          `json:"story_id,string"`
	StdStatus        string
	StdType          string
	Type             string
	Url              string

	SupportID       uint64  `json:"support_id,string"`
	SupportForumID  uint64  `json:"support_forum_id,string"`
	TicketID        uint64  `json:"ticket_id,string"`
	Follower        string  `json:"follower"`
	SyncType        string  `json:"sync_type"`
	Label           string  `json:"label"`
	Effort          float32 `json:"effort,string"`
	EffortCompleted float32 `json:"effort_completed,string"`
	Exceed          float32 `json:"exceed,string"`
	Remain          float32 `json:"remain,string"`
	Progress        string  `json:"progress"`
	Estimate        float32 `json:"estimate,string"`
	Bugtype         string  `json:"bugtype"`

	Milestone string `json:"milestone" gorm:"type:varchar(255)"`
}

func (TapdBug) TableName() string {
	return "_tool_tapd_bugs"
}
