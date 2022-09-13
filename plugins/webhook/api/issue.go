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

package api

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
	"net/http"
	"time"
)

type WebhookIssueRequest struct {
	BoardKey                string     `mapstructure:"board_key" validate:"required"`
	Url                     string     `mapstructure:"url"`
	IssueKey                string     `mapstructure:"issue_key" validate:"required"`
	Title                   string     `mapstructure:"title" validate:"required"`
	Description             string     `mapstructure:"description"`
	EpicKey                 string     `mapstructure:"epic_key"`
	Type                    string     `mapstructure:"type"`
	Status                  string     `mapstructure:"status" validate:"required" validate:"oneof=TODO DONE IN_PROGRESS"`
	OriginalStatus          string     `mapstructure:"original_status" validate:"required"`
	StoryPoint              int64      `mapstructure:"story_point"`
	ResolutionDate          *time.Time `mapstructure:"resolution_date"`
	CreatedDate             *time.Time `mapstructure:"created_date" validate:"required"`
	UpdatedDate             *time.Time `mapstructure:"updated_date"`
	LeadTimeMinutes         uint       `mapstructure:"lead_time_minutes"`
	ParentIssueKey          string     `mapstructure:"parent_issue_key"`
	Priority                string     `mapstructure:"priority"`
	OriginalEstimateMinutes int64      `mapstructure:"original_estimate_minutes"`
	TimeSpentMinutes        int64      `mapstructure:"time_spent_minutes"`
	TimeRemainingMinutes    int64      `mapstructure:"time_remaining_minutes"`
	CreatorId               string     `mapstructure:"creator_id"`
	CreatorName             string     `mapstructure:"creator_name"`
	AssigneeId              string     `mapstructure:"assignee_id"`
	AssigneeName            string     `mapstructure:"assignee_name"`
	Severity                string     `mapstructure:"severity"`
	Component               string     `mapstructure:"component"`
	//IconURL               string
	//DeploymentId          string
}

/*
PostIssue receive a record as below and save it
POST /plugins/webhook/:connectionId/issue
{
  "board_key": "DLK",
  "url": "",
  "issue_key": "DLK-1234",
  "title": "a feature from DLK",
  "description": "",
  "epic_key": "",
  "type": "BUG",
  "status": "TODO",
  "original_status": "created",
  "story_point": 0,
  "resolution_date": null,
  "created_date": "2020-01-01T12:00:00+00:00",
  "updated_date": null,
  "lead_time_minutes": 0,
  "parent_issue_key": "DLK-1200",
  "priority": "",
  "original_estimate_minutes": 0,
  "time_spent_minutes": 0,
  "time_remaining_minutes": 0,
  "creator_id": "user1131",
  "creator_name": "Nick name 1",
  "assignee_id": "user1132",
  "assignee_name": "Nick name 2",
  "severity": "",
  "component": ""
}
*/
func PostIssue(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	// TODO save issue
	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

/*
POST /plugins/webhook/:connectionId/issue/:issueId/close
{
	"name": "Webhook data connection name",
	"endpoint": "Webhook api endpoint, i.e. https://example.com",
	"username": "username, usually should be email address",
	"password": "Webhook api access token"
}
*/
func CloseIssue(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	// TODO close issue
	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}
