package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractIterations

var ExtractIterationMeta = core.SubTaskMeta{
	Name:             "extractIterations",
	EntryPoint:       ExtractIterations,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdIterationRes struct {
	Iteration models.TapdIteration
}

func ExtractIterations(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_ITERATION_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var iterBody TapdIterationRes
			err := json.Unmarshal(row.Data, &iterBody)
			if err != nil {
				return nil, err
			}
			iter := iterBody.Iteration

			iter.ConnectionId = data.Connection.ID
			workspaceIter := &models.TapdWorkspaceIteration{
				ConnectionId: data.Connection.ID,
				WorkspaceID:  iter.WorkspaceID,
				IterationId:  iter.ID,
			}
			return []interface{}{
				&iter, workspaceIter,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
