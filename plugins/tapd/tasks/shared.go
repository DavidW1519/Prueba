package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Page struct {
	Data Data `json:"data"`
}
type Data struct {
	Count int `json:"count"`
}

var UserIdGen *didgen.DomainIdGenerator
var WorkspaceIdGen *didgen.DomainIdGenerator
var IssueIdGen *didgen.DomainIdGenerator
var IterIdGen *didgen.DomainIdGenerator

// res will not be used
func GetTotalPagesFromResponse(r *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	data := args.Ctx.GetData().(*TapdTaskData)
	apiClient, err := NewTapdApiClient(args.Ctx.TaskContext(), data.Connection)
	if err != nil {
		return 0, err
	}
	query := url.Values{}
	query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceID))
	res, err := apiClient.Get(fmt.Sprintf("%s/count", r.Request.URL.Path), query, nil)
	if err != nil {
		return 0, err
	}
	var page Page
	err = helper.UnmarshalResponse(res, &page)

	count := page.Data.Count
	totalPage := count/args.PageSize + 1

	return totalPage, err
}

func parseIterationChangelog(taskCtx core.SubTaskContext, old string, new string) (uint64, uint64, error) {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDb()
	iterationFrom := &models.TapdIteration{}
	err := db.Model(&models.TapdIteration{}).
		Where("connection_id = ? and workspace_id = ? and name = ?",
			data.Connection.ID, data.Options.WorkspaceID, old).Limit(1).Find(iterationFrom).Error
	if err != nil {
		return 0, 0, err
	}
	iterationTo := &models.TapdIteration{}
	err = db.Model(&models.TapdIteration{}).
		Where("connection_id = ? and workspace_id = ? and name = ?",
			data.Connection.ID, data.Options.WorkspaceID, new).Limit(1).Find(iterationTo).Error
	if err != nil {
		return 0, 0, err
	}
	return iterationFrom.ID, iterationTo.ID, nil
}
func GetRawMessageDirectFromResponse(res *http.Response) ([]json.RawMessage, error) {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return []json.RawMessage{body}, nil
}
