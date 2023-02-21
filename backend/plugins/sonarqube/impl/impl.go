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

package impl

import (
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/dal"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"

	"github.com/apache/incubator-devlake/plugins/sonarqube/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/sonarqube/tasks"
)

// make sure interface is implemented
var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginTask
	plugin.PluginMigration
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
	plugin.PluginSource
} = (*Sonarqube)(nil)

type Sonarqube struct{}

func (p Sonarqube) Description() string {
	return "collect some Sonarqube data"
}

func (p Sonarqube) Init(br context.BasicRes) errors.Error {
	api.Init(br)
	return nil
}

func (p Sonarqube) Connection() interface{} {
	return &models.SonarqubeConnection{}
}

func (p Sonarqube) Scope() interface{} {
	return &models.SonarqubeProject{}
}

func (p Sonarqube) TransformationRule() interface{} {
	return nil
}

func (p Sonarqube) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.SonarqubeConnection{},
		&models.SonarqubeProject{},
		&models.SonarqubeIssue{},
		&models.SonarqubeIssueCodeBlock{},
		&models.SonarqubeHotspot{},
		&models.SonarqubeFileMetrics{},
		&models.SonarqubeAccount{},
	}
}

func (p Sonarqube) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectIssuesMeta,
		tasks.ExtractIssuesMeta,
		tasks.CollectHotspotsMeta,
		tasks.ExtractHotspotsMeta,
		tasks.CollectFilemetricsMeta,
		tasks.ExtractFilemetricsMeta,
		tasks.CollectAccountsMeta,
		tasks.ExtractAccountsMeta,
		tasks.ConvertProjectsMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueCodeBlocksMeta,
		tasks.ConvertHotspotsMeta,
		tasks.ConvertFileMetricsMeta,
		tasks.ConvertAccountsMeta,
	}
}

func (p Sonarqube) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.SonarqubeConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Sonarqube connection by the given connection ID")
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Sonarqube API client instance")
	}
	taskData := &tasks.SonarqubeTaskData{
		Options:   op,
		ApiClient: apiClient,
	}
	if op.ProjectKey != "" {
		var scope *models.SonarqubeProject
		var apiProject *tasks.SonarqubeApiProject
		// support v100 & advance mode
		// If we still cannot find the record in db, we have to request from remote server and save it to db
		db := taskCtx.GetDal()
		err = db.First(&scope, dal.Where("connection_id = ? AND project_key = ?", op.ConnectionId, op.ProjectKey))
		if err != nil && db.IsErrorNotFound(err) {
			apiProject, err = api.GetApiProject(op.ProjectKey, apiClient)
			if err != nil {
				return nil, err
			}
			logger.Debug(fmt.Sprintf("Current project: %s", apiProject.ProjectKey))
			scope = tasks.ConvertProject(apiProject)
			scope.ConnectionId = op.ConnectionId
			err = taskCtx.GetDal().CreateIfNotExist(&scope)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find project: %s", op.ProjectKey))
		}
	}
	var createdDateAfter time.Time
	if op.CreatedDateAfter != "" {
		createdDateAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.CreatedDateAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `createdDateAfter`")
		}
	}
	if !createdDateAfter.IsZero() {
		taskData.CreatedDateAfter = &createdDateAfter
		logger.Debug("collect data updated createdDateAfter %s", createdDateAfter)
	}
	return taskData, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p Sonarqube) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/sonarqube"
}

func (p Sonarqube) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Sonarqube) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": api.SearchRemoteScopes,
		},
		"connections/:connectionId/scopes/:projectKey": {
			"GET":   api.GetScope,
			"PATCH": api.UpdateScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
	}
}

func (p Sonarqube) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy) (pp plugin.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (p Sonarqube) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.SonarqubeTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
