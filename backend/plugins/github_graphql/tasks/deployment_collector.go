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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.SubTaskEntryPoint = CollectDeployment

const (
	RAW_DEPLOYMENT = "github_deployment"
)

type SingleEntityTaskGroup struct {
	Collector plugin.SubTaskMeta
	Extractor plugin.SubTaskMeta
	Convertor plugin.SubTaskMeta
}

var Deployment = SingleEntityTaskGroup{
	Collector: CollectDeploymentMeta,
	Extractor: ExtractDeploymentMeta,
	Convertor: ConvertDeploymentMeta,
}

var CollectDeploymentMeta = plugin.SubTaskMeta{
	Name:             "CollectDeployment",
	EntryPoint:       CollectDeployment,
	EnabledByDefault: true,
	Description:      "",
	DomainTypes:      []string{},
}

func CollectDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	return nil
}
