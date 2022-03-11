package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"regexp"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var _ core.SubTaskEntryPoint = ExtractApiIssues

type ApiIssuesResponse []IssuesResponse

type IssuesResponse struct {
	GithubId    int `json:"id"`
	Number      int
	State       string
	Title       string
	Body        string
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
	ClosedAt        *core.Iso8601Time `json:"closed_at"`
	GithubCreatedAt core.Iso8601Time  `json:"created_at"`
	GithubUpdatedAt core.Iso8601Time  `json:"updated_at"`
}

func ExtractApiIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	var issueSeverityRegex *regexp.Regexp
	var issueComponentRegex *regexp.Regexp
	var issuePriorityRegex *regexp.Regexp
	var issueTypeBugRegex *regexp.Regexp
	var issueTypeRequirementRegex *regexp.Regexp
	var issueTypeIncidentRegex *regexp.Regexp
	var issueSeverity = taskCtx.GetConfig("GITHUB_ISSUE_SEVERITY")
	var issueComponent = taskCtx.GetConfig("GITHUB_ISSUE_COMPONENT")
	var issuePriority = taskCtx.GetConfig("GITHUB_ISSUE_PRIORITY")
	var issueTypeBug = taskCtx.GetConfig("GITHUB_ISSUE_TYPE_BUG")
	var issueTypeRequirement = taskCtx.GetConfig("GITHUB_ISSUE_TYPE_REQUIREMENT")
	var issueTypeIncident = taskCtx.GetConfig("GITHUB_ISSUE_TYPE_INCIDENT")
	if len(issueSeverity) > 0 {
		issueSeverityRegex = regexp.MustCompile(issueSeverity)
	}
	if len(issueComponent) > 0 {
		issueComponentRegex = regexp.MustCompile(issueComponent)
	}
	if len(issuePriority) > 0 {
		issuePriorityRegex = regexp.MustCompile(issuePriority)
	}
	if len(issueTypeBug) > 0 {
		issueTypeBugRegex = regexp.MustCompile(issueTypeBug)
	}
	if len(issueTypeRequirement) > 0 {
		issueTypeRequirementRegex = regexp.MustCompile(issueTypeRequirement)
	}
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
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiIssuesResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, len(*body)*2)
			for _, apiIssue := range *body {
				if apiIssue.GithubId == 0 {
					return nil, nil
				}
				//If this is a pr, ignore
				if apiIssue.PullRequest.Url != "" {
					continue
				}
				githubIssue, err := convertGithubIssue(&apiIssue, data.Repo.GithubId)
				if err != nil {
					return nil, err
				}
				for _, label := range apiIssue.Labels {
					results = append(results, &models.GithubIssueLabel{
						IssueId:   githubIssue.GithubId,
						LabelName: label.Name,
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

			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
func convertGithubIssue(issue *IssuesResponse, repositoryId int) (*models.GithubIssue, error) {
	githubIssue := &models.GithubIssue{
		GithubId:        issue.GithubId,
		RepoId:          repositoryId,
		Number:          issue.Number,
		State:           issue.State,
		Title:           issue.Title,
		Body:            issue.Body,
		ClosedAt:        core.Iso8601TimeToTime(issue.ClosedAt),
		GithubCreatedAt: issue.GithubCreatedAt.ToTime(),
		GithubUpdatedAt: issue.GithubUpdatedAt.ToTime(),
	}

	if issue.Assignee != nil {
		githubIssue.AssigneeId = issue.Assignee.Id
		githubIssue.AssigneeName = issue.Assignee.Login
	}
	if issue.ClosedAt != nil {
		githubIssue.LeadTimeMinutes = uint(issue.ClosedAt.ToTime().Sub(issue.GithubCreatedAt.ToTime()).Minutes())
	}

	return githubIssue, nil
}
