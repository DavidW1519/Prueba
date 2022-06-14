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

import "github.com/apache/incubator-devlake/plugins/helper"

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required,url"`
	Auth     string `json:"auth" validate:"required"`
	Proxy    string `json:"proxy"`
}

type GithubConnection struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
}

type Config struct {
	PrType               string `mapstructure:"prType" json:"prType"`
	PrComponent          string `mapstructure:"prComponent" json:"prComponent"`
	PrBodyClosePattern   string `mapstructure:"prBodyClosePattern" json:"prBodyClosePattern"`
	IssueSeverity        string `mapstructure:"issueSeverity" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"issueTypeIncident" json:"issueTypeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" json:"issueTypeRequirement"`
}

func (GithubConnection) TableName() string {
	return "_tool_github_connections"
}

// Using Public Email because it requires authentication, and it is public information anyway.
// We're not using email information for anything here.
type PublicEmail struct {
	Email      string
	Primary    bool
	Verified   bool
	Visibility string
}
