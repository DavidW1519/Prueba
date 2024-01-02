/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	gocontext "context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/bitbucket_server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input,
		func(basicRes context.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.BitbucketServerConnection) ([]models.ProjectItem, errors.Error) {
			if gid != "" {
				return nil, nil
			}
			query := initialQuery(queryData)

			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			var res *http.Response
			res, err = apiClient.Get("rest/api/1.0/projects", query, nil)
			if err != nil {
				return nil, err
			}

			resBody := &models.ProjectsResponse{}
			err = api.UnmarshalResponse(res, resBody)
			if err != nil {
				return nil, err
			}

			return resBody.Values, err
		},
		func(basicRes context.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.BitbucketServerConnection) ([]models.BitbucketApiRepo, errors.Error) {
			if gid == "" {
				return nil, nil
			}
			query := initialQuery(queryData)

			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			var res *http.Response
			// list projects part
			res, err = apiClient.Get(fmt.Sprintf("rest/api/1.0/projects/%s/repos", gid), query, nil)
			if err != nil {
				return nil, err
			}
			var resBody models.ReposResponse
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			return resBody.Values, err
		},
	)
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/bitbucket_server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} api.SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.SearchRemoteScopes(input,
		func(basicRes context.BasicRes, queryData *api.RemoteQueryData, connection models.BitbucketServerConnection) ([]models.BitbucketApiRepo, errors.Error) {
			// create api client
			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}

			// request search
			query := initialQuery(queryData)
			if len(queryData.Search) == 0 {
				return nil, errors.BadInput.New("empty search query")
			}
			s := queryData.Search[0]
			gid, searchName := getSearch(s)
			// query.Set("sort", "name")
			// query.Set("fields", "values.name,values.full_name,values.language,values.description,values.owner.display_name,values.created_on,values.updated_on,values.links.clone,values.links.html,pagelen,page,size")
			queryString := fmt.Sprintf("name=%s", searchName)
			if len(gid) > 0 {
				queryString = fmt.Sprintf("&projectkey=%s", gid)
			}
			query.Set("q", queryString)

			// list repos part
			res, err := apiClient.Get("rest/api/1.0/repos", query, nil)
			if err != nil {
				return nil, err
			}

			var resBody models.ReposResponse
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}

			return resBody.Values, err
		},
	)
}

func getSearch(s string) (string, string) {
	gid := ""
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) == 2 {
			// KEY/repo
			gid = parts[0]
			s = strings.Join(parts[1:], "/")
		} else if len(parts) >= 3 {
			// KEY/repos/repo
			gid = parts[0]
			s = strings.Join(parts[2:], "/")
		}
	}
	return gid, s
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	start := (queryData.Page - 1) * queryData.PerPage
	query.Set("start", fmt.Sprintf("%v", start))
	query.Set("limit", fmt.Sprintf("%v", queryData.PerPage))
	return query
}
