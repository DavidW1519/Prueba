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

import "github.com/apache/incubator-devlake/core/models/common"

type ApiBambooProject struct {
	Key         string            `json:"key"`
	Expand      string            `json:"expand"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Link        ApiBambooLink     `json:"link"`
	Plans       ApiBambooSizeData `json:"plans"`
}

type ApiBambooProjects struct {
	ApiBambooSizeData
	Expand   string             `json:"expand"`
	Link     ApiBambooLink      `json:"link"`
	Projects []ApiBambooProject `json:"project"`
}

type ApiBambooProjectResponse struct {
	Expand   string            `json:"expand"`
	Link     ApiBambooLink     `json:"link"`
	Projects ApiBambooProjects `json:"projects"`
}

type BambooProject struct {
	ConnectionId         uint64 `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey"`
	Key                  string `json:"key" gorm:"primaryKey;type:varchar(256)"`
	TransformationRuleId uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId"`
	Name                 string `json:"name" gorm:"index;type:varchar(256)"`
	Description          string `json:"description"`
	Href                 string `json:"link"`
	Rel                  string `json:"rel" gorm:"type:varchar(100)"`
	common.NoPKModel     `json:"-" mapstructure:"-"`
}

func (b *BambooProject) Convert(apiProject *ApiBambooProject) {
	b.Key = apiProject.Key
	b.Name = apiProject.Name
	b.Description = apiProject.Description
	b.Href = apiProject.Link.Href
}

func (b *BambooProject) TableName() string {
	return "_tool_bamboo_projects"
}
