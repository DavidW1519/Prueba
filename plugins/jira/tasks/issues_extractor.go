package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"strings"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractIssues

func ExtractIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	sourceId := data.Source.ID
	boardId := data.Options.BoardId
	db := taskCtx.GetDb()
	logger := taskCtx.GetLogger()
	logger.Info("extract Issues, source_id=%d, board_id=%d", sourceId, boardId)
	// prepare getStdType function
	var typeMappingRows []*models.JiraIssueTypeMapping
	err := db.Find(&typeMappingRows, "source_id = ?", sourceId).Error
	if err != nil {
		return err
	}
	typeMappings := make(map[string]string)
	for _, typeMappingRow := range typeMappingRows {
		typeMappings[typeMappingRow.UserType] = typeMappingRow.StandardType
	}
	getStdType := func(userType string) string {
		stdType := typeMappings[userType]
		if stdType == "" {
			return strings.ToUpper(userType)
		}
		return strings.ToUpper(stdType)
	}
	// prepare getStdStatus function
	var statusMappingRows []*models.JiraIssueStatusMapping
	err = db.Find(&statusMappingRows, "source_id = ?", sourceId).Error
	if err != nil {
		return err
	}
	statusMappings := make(map[string]string)
	makeStatusMappingKey := func(userType string, userStatus string) string {
		return fmt.Sprintf("%v:%v", userType, userStatus)
	}
	for _, statusMappingRow := range statusMappingRows {
		k := makeStatusMappingKey(statusMappingRow.UserType, statusMappingRow.UserStatus)
		statusMappings[k] = statusMappingRow.StandardStatus
	}
	getStdStatus := func(statusKey string) string {
		if statusKey == "done" {
			return ticket.DONE
		} else if statusKey == "new" {
			return ticket.TODO
		} else {
			return ticket.IN_PROGRESS
		}
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var apiIssue apiv2models.Issue
			err := json.Unmarshal(row.Data, &apiIssue)
			if err != nil {
				return nil, err
			}
			var results []interface{}
			sprints, issue, _, worklogs, changelogs, changelogItems := apiIssue.ExtractEntities(data.Source.ID, data.Source.EpicKeyField, data.Source.StoryPointField)
			for _, sprintId := range sprints {
				sprintIssue := &models.JiraSprintIssue{
					SourceId: data.Source.ID,
					SprintId: sprintId,
					IssueId:  issue.IssueId,
				}
				results = append(results, sprintIssue)
			}
			if issue.ResolutionDate != nil {
				issue.LeadTimeMinutes = uint(issue.ResolutionDate.Unix()-issue.Created.Unix()) / 60
			}
			issue.StdStoryPoint = uint(issue.StoryPoint)
			issue.StdType = getStdType(issue.Type)
			issue.StdStatus = getStdStatus(issue.StatusKey)
			issue.SpentMinutes = issue.AggregateEstimateMinutes - issue.RemainingEstimateMinutes
			results = append(results, issue)
			for _, worklog := range worklogs {
				results = append(results, worklog)
			}
			for _, changelog := range changelogs {
				results = append(results, changelog)
			}
			for _, changelogItem := range changelogItems {
				results = append(results, changelogItem)
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
