package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

type JiraApiAuthor struct {
	Self        string `json:"self,omitempty"`
	AccountId   string `json:"accountId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Active      bool   `json:"active,omitempty"`
	TimeZone    string `json:"timeZone,omitempty"`
	AccountType string `json:"accountType,omitempty"`
}

type JiraApiChangelogItem struct {
	Field      string `json:"field,omitempty"`
	FieldType  string `json:"fieldType,omitempty"`
	FieldId    string `json:"fieldId,omitempty"`
	From       string `json:"from,omitempty"`
	FromString string `json:"fromString,omitempty"`
	To         string `json:"to,omitempty"`
	ToString   string `json:"toString,omitempty"`
}

type JiraApiChangeLog struct {
	Id      string                 `json:"id,omitempty"`
	Author  JiraApiAuthor          `json:"author,omitempty"`
	Created core.Iso8601Time       `json:"created,omitempty"`
	Items   []JiraApiChangelogItem `json:"items,omitempty"`
}

type JiraApiChangelogsResponse struct {
	JiraPagination
	Values []JiraApiChangeLog `json:"values,omitempty"`
}

func CollectChangelogs(boardId uint64) error {
	jiraIssue := &models.JiraIssue{}

	// select all issues belongs to the board
	cursor, err := lakeModels.Db.Model(jiraIssue).
		Select("jira_issues.id").
		Joins("left join jira_board_issues on jira_board_issues.issue_id = jira_issues.id").
		Where("jira_board_issues.board_id = ?", boardId).
		Rows()
	if err != nil {
		return err
	}

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		err = collectChangelogsByIssueId(jiraIssue.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func collectChangelogsByIssueId(issueId uint64) error {
	jiraApiClient := GetJiraApiClient()
	return jiraApiClient.FetchPages(fmt.Sprintf("/api/3/issue/%v/changelog", issueId), nil,
		func(res *http.Response) error {
			// parse response
			jiraApiChangelogResponse := &JiraApiChangelogsResponse{}
			err := core.UnmarshalResponse(res, jiraApiChangelogResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}

			// process changelogs
			for _, jiraApiChangelog := range jiraApiChangelogResponse.Values {

				jiraChangelog, err := convertChangelog(&jiraApiChangelog)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				// save changelog
				err = lakeModels.Db.Save(jiraChangelog).Error
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}

				// process changelog items
				lakeModels.Db.Delete(models.JiraChangelogItem{}, "changelog_id = ?", jiraChangelog.ID)
				for _, jiraApiChangelogItem := range jiraApiChangelog.Items {
					jiraChangelogItem, err := convertChangelogItem(jiraChangelog.ID, &jiraApiChangelogItem)
					if err != nil {
						logger.Error("Error: ", err)
						return err
					}
					// save changelog item
					err = lakeModels.Db.Create(jiraChangelogItem).Error
					if err != nil {
						logger.Error("Error: ", err)
						return err
					}
				}
			}
			return nil
		})
}

func convertChangelog(jiraApiChangelog *JiraApiChangeLog) (*models.JiraChangelog, error) {
	id, err := strconv.ParseUint(jiraApiChangelog.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	return &models.JiraChangelog{
		Model:             lakeModels.Model{ID: id},
		AuthorAccountId:   jiraApiChangelog.Author.AccountId,
		AuthorDisplayName: jiraApiChangelog.Author.DisplayName,
		AuthorActive:      jiraApiChangelog.Author.Active,
		Created:           jiraApiChangelog.Created.ToTime(),
	}, nil
}

func convertChangelogItem(changelogId uint64, jiraApiChangeItem *JiraApiChangelogItem) (*models.JiraChangelogItem, error) {
	return &models.JiraChangelogItem{
		ChangelogId: changelogId,
		Field:       jiraApiChangeItem.Field,
		FieldType:   jiraApiChangeItem.FieldType,
		FieldId:     jiraApiChangeItem.FieldId,
		From:        jiraApiChangeItem.From,
		FromString:  jiraApiChangeItem.FromString,
		To:          jiraApiChangeItem.To,
		ToString:    jiraApiChangeItem.ToString,
	}, nil
}
