package tasks

import (
	"reflect"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

func ConvertApiMergeRequests(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)

	domainMrIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	domainRepoIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabProject{})

	//Find all piplines associated with the current projectid
	cursor, err := lakeModels.Db.Model(&models.GitlabMergeRequest{}).Where("project_id=?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMergeRequest{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabMr := inputRow.(*models.GitlabMergeRequest)
			err = lakeModels.Db.ScanRows(cursor, gitlabMr)

			domainPr := &code.PullRequest{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainMrIdGenerator.Generate(gitlabMr.GitlabId),
				},
				RepoId:      domainRepoIdGenerator.Generate(gitlabMr.ProjectId),
				Status:      gitlabMr.State,
				Title:       gitlabMr.Title,
				Url:         gitlabMr.WebUrl,
				CreatedDate: gitlabMr.GitlabCreatedAt,
				MergedDate:  gitlabMr.MergedAt,
				ClosedAt:    gitlabMr.ClosedAt,
			}

			return []interface{}{
				domainPr,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
