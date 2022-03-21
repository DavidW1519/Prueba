package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jenkins/models"
)

// this struct should be moved to `gitub_api_common.go`

var _ core.SubTaskEntryPoint = ExtractApiJobs

func ExtractApiJobs(taskCtx core.SubTaskContext) error {
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			/*
				Table store raw data
			*/
			Table: RAW_JOB_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &models.Job{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)

			job := &models.JenkinsJob{
				JenkinsJobProps: models.JenkinsJobProps{
					Name:  body.Name,
					Class: body.Class,
					Color: body.Color,
				},
			}
			results = append(results, job)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
