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
	"github.com/apache/incubator-devlake/plugins/helper"
)

// This object conforms to what the frontend currently sends.
type GitlabConnection struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
}

type TestConnectionRequest struct {
	Endpoint           string `json:"endpoint"`
	Proxy              string `json:"proxy"`
	helper.AccessToken `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type GitlabResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	GitlabConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

type Config struct {
	MrType               string `mapstructure:"MrType" env:"GITLAB_PR_TYPE" json:"MrType"`
	MrComponent          string `mapstructure:"MrComponent" env:"GITLAB_PR_COMPONENT" json:"MrComponent"`
	IssueSeverity        string `mapstructure:"issueSeverity" env:"GITLAB_ISSUE_SEVERITY" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority" env:"GITLAB_ISSUE_PRIORITY" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent" env:"GITLAB_ISSUE_COMPONENT" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" env:"GITLAB_ISSUE_TYPE_BUG" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"issueTypeIncident" env:"GITLAB_ISSUE_TYPE_INCIDENT" json:"issueTypeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" env:"GITLAB_ISSUE_TYPE_REQUIREMENT" json:"issueTypeRequirement"`
}

func (GitlabConnection) TableName() string {
	return "_tool_gitlab_connections"
}
