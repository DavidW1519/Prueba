package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/helper"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/plugins/core"
)

const RAW_REPOSITORIES_TABLE = "github_api_repositories"

// this struct should be moved to `gitub_api_common.go`

var _ core.SubTaskEntryPoint = CollectApiRepositories

func CollectApiRepositories(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_REPOSITORIES_TABLE,
		},
		ApiClient: data.ApiClient,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{body}, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
