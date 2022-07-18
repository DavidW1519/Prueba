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

package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"time"
)

var ExtractApiIssueCommentsMeta = core.SubTaskMeta{
	Name:             "extractApiIssueComments",
	EntryPoint:       ExtractApiIssueComments,
	EnabledByDefault: true,
	Required:         true,
	Description:      "Extract raw issue comments data into tool layer table BitbucketIssueComments",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

type BitbucketIssueCommentsResponse struct {
	Type        string    `json:"type"`
	BitbucketId int       `json:"id"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
	Content     struct {
		Type string
		Raw  string
	} `json:"content"`
	User  *BitbucketAccountResponse
	Issue struct {
		Type       string
		Id         int
		Repository *BitbucketApiRepo
		Links      struct {
			Self struct {
				Href string
			}
		}
		Title string
	}
	Links struct {
		Self struct {
			Href string
		} `json:"self"`
		Html struct {
			Href string
		} `json:"html"`
	} `json:"links"`
}

func ExtractApiIssueComments(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_COMMENTS_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			issueComment := &BitbucketIssueCommentsResponse{}
			err := json.Unmarshal(row.Data, issueComment)
			if err != nil {
				return nil, err
			}

			toolIssueComment, err := convertIssueComment(issueComment)
			toolIssueComment.ConnectionId = data.Options.ConnectionId
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 2)

			results = append(results, toolIssueComment)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertIssueComment(issueComment *BitbucketIssueCommentsResponse) (*models.BitbucketIssueComment, error) {
	bitbucketIssueComment := &models.BitbucketIssueComment{
		BitbucketId:    issueComment.BitbucketId,
		AuthorUserId:   issueComment.User.AccountId,
		IssueId:        issueComment.Issue.Id,
		AuthorUsername: issueComment.User.DisplayName,
		CreatedAt:      issueComment.CreatedOn,
		Type:           issueComment.Type,
	}
	return bitbucketIssueComment, nil
}
