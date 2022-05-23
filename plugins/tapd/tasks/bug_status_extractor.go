package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractBugStatus

var ExtractBugStatusMeta = core.SubTaskMeta{
	Name:             "extractBugStatus",
	EntryPoint:       ExtractBugStatus,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_bugStatus",
}

type TapdBugStatusRes struct {
	Data map[string]string
}

func ExtractBugStatus(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_BUG_STATUS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var statusRes TapdBugStatusRes
			err := json.Unmarshal(row.Data, &statusRes)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0)
			for k, v := range statusRes.Data {
				toolL := &models.TapdBugStatus{
					ConnectionId: data.Connection.ID,
					WorkspaceID:  data.Options.WorkspaceID,
					EnglishName:  k,
					ChineseName:  v,
					IsLastStep:   false,
				}
				results = append(results, toolL)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
