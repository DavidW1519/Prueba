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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
)

type OpsgenieParams struct {
	ConnectionId uint64
	ScopeId      string
}

type Service struct {
	common.NoPKModel
	ConnectionId uint64 `json:"connection_id" mapstructure:"connectionId,omitempty" gorm:"primaryKey" `
	Id           string `json:"id" mapstructure:"id" gorm:"primaryKey;autoIncrement:false" `
	Url          string `json:"url" mapstructure:"url"`
	Name         string `json:"name" mapstructure:"name"`
	TeamId       string `json:"team_id" mapstructure:"team_id"`
}

func (s Service) ScopeId() string {
	return s.Name
}

func (s Service) ScopeConnectionId() uint64 {
	return s.ConnectionId
}

func (s Service) ScopeScopeConfigId() uint64 {
	return 0
}

func (s Service) ScopeName() string {
	return s.Name
}

func (s Service) ScopeFullName() string {
	return s.Name
}

func (s Service) ScopeParams() interface{} {
	return &OpsgenieParams{
		ConnectionId: s.ConnectionId,
		ScopeId:      s.Id,
	}
}

func (s Service) TableName() string {
	return "_tool_opsgenie_services"
}

var _ plugin.ToolLayerScope = (*Service)(nil)
