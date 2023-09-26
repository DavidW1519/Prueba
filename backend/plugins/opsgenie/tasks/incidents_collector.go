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

package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_INCIDENTS_TABLE = "opsgenie_incidents"

var _ plugin.SubTaskEntryPoint = CollectIncidents

type (
	collectedIncidents struct {
		TotalCount int               `json:"totalCount"`
		Data       []json.RawMessage `json:"data"`
	}
)

var CollectIncidentsMeta = plugin.SubTaskMeta{
	Name:             "collectIncidents",
	EntryPoint:       CollectIncidents,
	EnabledByDefault: true,
	Description:      "Collect Opsgenie incidents",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*OpsgenieTaskData)

	collectorWithState, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_INCIDENTS_TABLE,
	})
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.Client,
		PageSize:    100,
		UrlTemplate: "v1/incidents",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}

			query.Set("query", fmt.Sprintf("impactedServices:%s", data.Options.ServiceId))
			query.Set("sort", "createdAt")
			query.Set("order", "desc")
			query.Set("limit", fmt.Sprintf("%d", reqData.Pager.Size))
			query.Set("offset", fmt.Sprintf("%d", reqData.Pager.Skip))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			rawResult := collectedIncidents{}
			err := api.UnmarshalResponse(res, &rawResult)

			return rawResult.Data, err
		},
	})

	if err != nil {
		return errors.Default.Wrap(err, "error collecting opsgenie incidents")
	}

	return collectorWithState.Execute()
}
