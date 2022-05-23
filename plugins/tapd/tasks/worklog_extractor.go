package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractWorklogs

var ExtractWorklogMeta = core.SubTaskMeta{
	Name:             "extractWorklogs",
	EntryPoint:       ExtractWorklogs,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdWorklogRes struct {
	Timesheet models.TapdWorklog
}

func ExtractWorklogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_WORKLOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var worklogBody TapdWorklogRes
			err := json.Unmarshal(row.Data, &worklogBody)
			if err != nil {
				return nil, err
			}
			toolL := worklogBody.Timesheet

			toolL.ConnectionId = data.Connection.ID
			results := make([]interface{}, 0, 1)
			results = append(results, &toolL)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
