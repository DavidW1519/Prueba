package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type GitlabApiParams struct {
	ProjectId int
}

type GitlabInput struct {
	GitlabId int
	Iid      int
}

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	total := res.Header.Get("X-Total-Pages")
	if total == "" {
		return 0, nil
	}
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		return 0, err
	}
	return totalInt, nil
}

func GetRawMessageFromResponse(res *http.Response) ([]json.RawMessage, error) {
	rawMessages := []json.RawMessage{}

	if res == nil {
		return nil, fmt.Errorf("res is nil")
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, res.Request.URL.String())
	}

	err = json.Unmarshal(resBody, &rawMessages)
	if err != nil {
		return nil, fmt.Errorf("%w %s %s", err, res.Request.URL.String(), string(resBody))
	}

	return rawMessages, nil
}

func GetQuery(reqData *helper.RequestData) (url.Values, error) {
	query := url.Values{}
	query.Set("with_stats", "true")
	query.Set("sort", "asc")
	query.Set("page", strconv.Itoa(reqData.Pager.Page))
	query.Set("per_page", strconv.Itoa(reqData.Pager.Size))
	return query, nil
}

func CreateRawDataSubTaskArgs(taskCtx core.SubTaskContext, Table string) (*helper.RawDataSubTaskArgs, *GitlabTaskData) {
	data := taskCtx.GetData().(*GitlabTaskData)
	RawDataSubTaskArgs := &helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GitlabApiParams{
			ProjectId: data.Options.ProjectId,
		},
		Table: Table,
	}
	return RawDataSubTaskArgs, data
}

func GetMergeRequestsIterator(taskCtx core.SubTaskContext) (*helper.CursorIterator, error) {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GitlabTaskData)
	cursor, err := db.Model(&models.GitlabMergeRequest{}).Where("project_id = ?", data.Options.ProjectId).Select("gitlab_id,iid").Rows()
	if err != nil {
		return nil, err
	}

	return helper.NewCursorIterator(db, cursor, reflect.TypeOf(GitlabInput{}))
}

func GetPipelinesIterator(taskCtx core.SubTaskContext) (*helper.CursorIterator, error) {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GitlabTaskData)
	cursor, err := db.Model(&models.GitlabPipeline{}).Where("project_id = ?", data.Options.ProjectId).Select("gitlab_id").Rows()
	if err != nil {
		return nil, err
	}

	return helper.NewCursorIterator(db, cursor, reflect.TypeOf(GitlabInput{}))
}
