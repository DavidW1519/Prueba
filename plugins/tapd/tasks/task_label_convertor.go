package tasks

import (
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertTaskLabelsMeta = core.SubTaskMeta{
	Name:             "convertTaskLabels",
	EntryPoint:       ConvertTaskLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table tapd_issue_labels into  domain layer table issue_labels",
}

func ConvertTaskLabels(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*TapdTaskData)

	cursor, err := db.Model(&models.TapdTaskLabel{}).
		Joins(`left join _tool_tapd_workspace_tasks on _tool_tapd_workspace_tasks.task_id = _tool_tapd_task_labels.task_id`).
		Where("_tool_tapd_workspace_tasks.workspace_id = ?", data.Options.WorkspaceID).
		Order("task_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,

				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_BUG_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdTaskLabel{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issueLabel := inputRow.(*models.TapdTaskLabel)
			domainTaskLabel := &ticket.IssueLabel{
				IssueId:   IssueIdGen.Generate(issueLabel.TaskId),
				LabelName: issueLabel.LabelName,
			}
			return []interface{}{
				domainTaskLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
