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

package archived

import (
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type BaseConnection struct {
	Name string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	archived.Model
}

type RestConnection struct {
	BaseConnection   `mapstructure:",squash"`
	Endpoint         string `mapstructure:"endpoint" validate:"required" json:"endpoint"`
	Proxy            string `mapstructure:"proxy" json:"proxy"`
	RateLimitPerHour int    `comment:"api request rate limit per hour" json:"rateLimit"`
}

type SonarqubeConnection struct {
	RestConnection `mapstructure:",squash"`
	AccessToken    `mapstructure:",squash"`
}

type TestConnectionRequest struct {
	Endpoint    string `json:"endpoint"`
	Proxy       string `json:"proxy"`
	AccessToken `mapstructure:",squash"`
}

// For sonarqube, we can `use user_token:`
type AccessToken struct {
	Token string `mapstructure:"token" validate:"required" json:"token" encrypt:"yes"`
}

// This object conforms to what the frontend currently expects.
type SonarqubeResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	SonarqubeConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

func (SonarqubeConnection) TableName() string {
	return "_tool_sonarqube_connections"
}
