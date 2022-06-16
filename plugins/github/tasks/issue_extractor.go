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
	"regexp"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiIssuesMeta = core.SubTaskMeta{
	Name:             "extractApiIssues",
	EntryPoint:       ExtractApiIssues,
	EnabledByDefault: true,
	Description:      "Extract raw Issues data into tool layer table github_issues",
}

type IssuesResponse struct {
	GithubId    int `json:"id"`
	Number      int
	State       string
	Title       string
	Body        json.RawMessage
	HtmlUrl     string `json:"html_url"`
	PullRequest struct {
		Url     string `json:"url"`
		HtmlUrl string `json:"html_url"`
	} `json:"pull_request"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`

	Assignee *struct {
		Login string
		Id    int
	}
	User *struct {
		Login string
		Id    int
	}
	ClosedAt        *helper.Iso8601Time `json:"closed_at"`
	GithubCreatedAt helper.Iso8601Time  `json:"created_at"`
	GithubUpdatedAt helper.Iso8601Time  `json:"updated_at"`
}

func ExtractApiIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	config := data.Options.Config
	var issueSeverityRegex *regexp.Regexp
	var issueComponentRegex *regexp.Regexp
	var issuePriorityRegex *regexp.Regexp
	var issueTypeBugRegex *regexp.Regexp
	var issueTypeRequirementRegex *regexp.Regexp
	var issueTypeIncidentRegex *regexp.Regexp
	var issueSeverity = config.IssueSeverity
	if len(issueSeverity) > 0 {
		issueSeverityRegex = regexp.MustCompile(issueSeverity)
	}
	var issueComponent = config.IssueComponent
	if len(issueComponent) > 0 {
		issueComponentRegex = regexp.MustCompile(issueComponent)
	}
	var issuePriority = config.IssuePriority
	if len(issuePriority) > 0 {
		issuePriorityRegex = regexp.MustCompile(issuePriority)
	}
	var issueTypeBug = config.IssueTypeBug
	if len(issueTypeBug) > 0 {
		issueTypeBugRegex = regexp.MustCompile(issueTypeBug)
	}
	var issueTypeRequirement = config.IssueTypeRequirement
	if len(issueTypeRequirement) > 0 {
		issueTypeRequirementRegex = regexp.MustCompile(issueTypeRequirement)
	}
	var issueTypeIncident = config.IssueTypeIncident
	if len(issueTypeIncident) > 0 {
		issueTypeIncidentRegex = regexp.MustCompile(issueTypeIncident)
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &IssuesResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			if body.GithubId == 0 {
				return nil, nil
			}
			//If this is a pr, ignore
			if body.PullRequest.Url != "" {
				return nil, nil
			}
			results := make([]interface{}, 0, 2)
			githubIssue, err := convertGithubIssue(body, data.Options.ConnectionId, data.Repo.GithubId)
			if err != nil {
				return nil, err
			}
			for _, label := range body.Labels {
				results = append(results, &models.GithubIssueLabel{
					ConnectionId: data.Options.ConnectionId,
					IssueId:      githubIssue.GithubId,
					LabelName:    label.Name,
				})
				if issueSeverityRegex != nil {
					groups := issueSeverityRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubIssue.Severity = groups[1]
					}
				}

				if issueComponentRegex != nil {
					groups := issueComponentRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubIssue.Component = groups[1]
					}
				}

				if issuePriorityRegex != nil {
					groups := issuePriorityRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubIssue.Priority = groups[1]
					}
				}

				if issueTypeBugRegex != nil {
					if ok := issueTypeBugRegex.MatchString(label.Name); ok {
						githubIssue.Type = ticket.BUG
					}
				}

				if issueTypeRequirementRegex != nil {
					if ok := issueTypeRequirementRegex.MatchString(label.Name); ok {
						githubIssue.Type = ticket.REQUIREMENT
					}
				}

				if issueTypeIncidentRegex != nil {
					if ok := issueTypeIncidentRegex.MatchString(label.Name); ok {
						githubIssue.Type = ticket.INCIDENT
					}
				}
			}
			results = append(results, githubIssue)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
func convertGithubIssue(issue *IssuesResponse, connectionId uint64, repositoryId int) (*models.GithubIssue, error) {
	githubIssue := &models.GithubIssue{
		ConnectionId:    connectionId,
		GithubId:        issue.GithubId,
		RepoId:          repositoryId,
		Number:          issue.Number,
		State:           issue.State,
		Title:           issue.Title,
		Body:            string(issue.Body),
		Url:             issue.HtmlUrl,
		ClosedAt:        helper.Iso8601TimeToTime(issue.ClosedAt),
		GithubCreatedAt: issue.GithubCreatedAt.ToTime(),
		GithubUpdatedAt: issue.GithubUpdatedAt.ToTime(),
	}

	if issue.Assignee != nil {
		githubIssue.AssigneeId = issue.Assignee.Id
		githubIssue.AssigneeName = issue.Assignee.Login
	}
	if issue.User != nil {
		githubIssue.AuthorId = issue.User.Id
		githubIssue.AuthorName = issue.User.Login
	}
	if issue.ClosedAt != nil {
		githubIssue.LeadTimeMinutes = uint(issue.ClosedAt.ToTime().Sub(issue.GithubCreatedAt.ToTime()).Minutes())
	}

	return githubIssue, nil
}
