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

package e2e

import (
	"testing"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
)

func TestTapdStoryChangelogDataFlow(t *testing.T) {

	var tapd impl.Tapd
	dataflowTester := e2ehelper.NewDataFlowTester(t, "tapd", tapd)

	taskData := &tasks.TapdTaskData{
		Options: &tasks.TapdOptions{
			ConnectionId: 1,
			CompanyId:    99,
			WorkspaceId:  991,
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_story_changelogs.csv",
		"_raw_tapd_api_story_changelogs")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdStoryChangelog{})
	dataflowTester.FlushTabler(&models.TapdStoryChangelogItem{})
	dataflowTester.Subtask(tasks.ExtractStoryChangelogMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdStoryChangelog{},
		"./snapshot_tables/_tool_tapd_story_changelogs.csv",
		[]string{
			"connection_id",
			"id",
			"workspace_id",
			"workitem_type_id",
			"creator",
			"created",
			"change_summary",
			"comment",
			"entity_type",
			"change_type",
			"story_id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	dataflowTester.VerifyTable(
		models.TapdStoryChangelogItem{},
		"./snapshot_tables/_tool_tapd_story_changelog_items.csv",
		[]string{
			"connection_id",
			"changelog_id",
			"field",
			"value_before_parsed",
			"value_after_parsed",
			"iteration_id_from",
			"iteration_id_to",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	dataflowTester.FlushTabler(&ticket.IssueChangelogs{})
	dataflowTester.Subtask(tasks.ConvertStoryChangelogMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueChangelogs{},
		"./snapshot_tables/issue_changelogs_story.csv",
		[]string{
			"id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"issue_id",
			"author_id",
			"author_name",
			"field_id",
			"field_name",
			"from_value",
			"to_value",
			"created_date",
			"original_from_value",
			"original_to_value",
		},
	)
}
