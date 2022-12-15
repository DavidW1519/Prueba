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
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var _ interface {
	core.PluginMeta
	core.PluginInit
	core.PluginTask
	core.PluginModel
	core.PluginMigration
	core.PluginBlueprintV100
	core.DataSourcePluginBlueprintV200
	core.CloseablePluginTask
	core.PluginSource
} = (*Gitlab)(nil)

type Gitlab string

func (plugin Gitlab) Init(basicRes core.BasicRes) errors.Error {
	api.Init(basicRes)
	return nil
}

func (plugin Gitlab) Connection() interface{} {
	return &models.GitlabConnection{}
}

func (plugin Gitlab) Scope() interface{} {
	return &models.GitlabProject{}
}

func (plugin Gitlab) TransformationRule() interface{} {
	return &models.GitlabTransformationRule{}
}

func (plugin Gitlab) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*core.BlueprintScopeV200, syncPolicy core.BlueprintSyncPolicy) (core.PipelinePlan, []core.Scope, errors.Error) {
	return api.MakePipelinePlanV200(plugin.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}

func (plugin Gitlab) GetTablesInfo() []core.Tabler {
	return []core.Tabler{
		&models.GitlabConnection{},
		&models.GitlabAccount{},
		&models.GitlabCommit{},
		&models.GitlabIssue{},
		&models.GitlabIssueLabel{},
		&models.GitlabJob{},
		&models.GitlabMergeRequest{},
		&models.GitlabMrComment{},
		&models.GitlabMrCommit{},
		&models.GitlabMrLabel{},
		&models.GitlabMrNote{},
		&models.GitlabPipeline{},
		&models.GitlabPipelineProject{},
		&models.GitlabProject{},
		&models.GitlabProjectCommit{},
		&models.GitlabReviewer{},
		&models.GitlabTag{},
	}
}

func (plugin Gitlab) Description() string {
	return "To collect and enrich data from Gitlab"
}

func (plugin Gitlab) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,
		tasks.CollectApiMergeRequestsMeta,
		tasks.ExtractApiMergeRequestsMeta,
		tasks.CollectApiMrNotesMeta,
		tasks.ExtractApiMrNotesMeta,
		tasks.CollectApiMrCommitsMeta,
		tasks.ExtractApiMrCommitsMeta,
		tasks.CollectApiPipelinesMeta,
		tasks.ExtractApiPipelinesMeta,
		tasks.CollectApiJobsMeta,
		tasks.ExtractApiJobsMeta,
		tasks.EnrichMergeRequestsMeta,
		tasks.CollectAccountsMeta,
		tasks.ExtractAccountsMeta,
		tasks.ConvertAccountsMeta,
		tasks.ConvertProjectMeta,
		tasks.ConvertApiMergeRequestsMeta,
		tasks.ConvertMrCommentMeta,
		tasks.ConvertApiMrCommitsMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueLabelsMeta,
		tasks.ConvertMrLabelsMeta,
		tasks.ConvertCommitsMeta,
		tasks.ConvertPipelineMeta,
		tasks.ConvertPipelineCommitMeta,
		tasks.ConvertJobMeta,
	}
}

func (plugin Gitlab) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}
	connection := &models.GitlabConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "connection not found")
	}
	apiClient, err := tasks.NewGitlabApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	var createdDateAfter time.Time
	if op.CreatedDateAfter != "" {
		createdDateAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.CreatedDateAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `createdDateAfter`")
		}
	}

	if op.GitlabTransformationRule == nil && op.ProjectId != 0 {
		repo, err := api.GetRepoByConnectionIdAndscopeId(op.ConnectionId, strconv.Itoa(op.ProjectId))
		if err != nil {
			return nil, err
		}
		transformationRule, err := api.GetTransformationRuleByRepo(repo)
		if err != nil {
			return nil, err
		}
		op.GitlabTransformationRule = transformationRule
	}

	taskData := tasks.GitlabTaskData{
		Options:   op,
		ApiClient: apiClient,
	}

	if !createdDateAfter.IsZero() {
		taskData.CreatedDateAfter = &createdDateAfter
		logger.Debug("collect data updated createdDateAfter %s", createdDateAfter)
	}
	return &taskData, nil
}

func (plugin Gitlab) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitlab"
}

func (plugin Gitlab) MigrationScripts() []core.MigrationScript {
	return migrationscripts.All()
}

func (plugin Gitlab) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin Gitlab) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
			"GET":    api.GetConnection,
		},
		"connections/:connectionId/scopes/:projectId": {
			"GET":   api.GetScope,
			"PUT":   api.PutScope,
			"PATCH": api.UpdateScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
		},
		"transformation_rules": {
			"POST": api.CreateTransformationRule,
			"GET":  api.GetTransformationRuleList,
		},
		"transformation_rules/:id": {
			"PATCH": api.UpdateTransformationRule,
			"GET":   api.GetTransformationRule,
		},
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}

func (plugin Gitlab) Close(taskCtx core.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.GitlabTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
