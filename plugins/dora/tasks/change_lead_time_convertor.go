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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertChangeLeadTimeMeta = core.SubTaskMeta{
	Name:             "ConvertChangeLeadTime",
	EntryPoint:       ConvertChangeLeadTime,
	EnabledByDefault: true,
	Description:      "TODO",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

const RAW_ISSUES_TABLE = `dora_issues`

func ConvertChangeLeadTime(taskCtx core.SubTaskContext) error {
	//db := taskCtx.GetDal()
	//data := taskCtx.GetData().(*DoraTaskData)

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: DoraApiParams{
				// TODO
			},
			Table: RAW_ISSUES_TABLE,
		},
		//InputRowType: reflect.TypeOf(githubModels.GithubJob{}),
		//Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			// TODO

			return []interface{}{}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
