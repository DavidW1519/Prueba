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

package securitytesting

import "github.com/apache/incubator-devlake/core/models/domainlayer"

type StIssueCodeBlock struct {
	domainlayer.DomainEntity
	Id          string `gorm:"primaryKey"`
	IssueKey    string `json:"key" gorm:"index"`
	Component   string `json:"component" gorm:"index"`
	Project     string `json:"project" gorm:"index"`
	Msg         string `json:"msg" `
	StartLine   int    `json:"startLine" `
	EndLine     int    `json:"endLine" `
	StartOffset int    `json:"startOffset" `
	EndOffset   int    `json:"endOffset" `
}

func (StIssueCodeBlock) TableName() string {
	return "st_issue_code_blocks"
}
