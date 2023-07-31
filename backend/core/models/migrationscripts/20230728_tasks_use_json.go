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

package migrationscripts

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*tasksUsesJSON)(nil)

type tasksUsesJSON struct{}

type srcTaskSubtaskJSON20230731 struct {
	archived.Model
	Subtasks json.RawMessage
}

type dstTaskSubtaskJSON20230731 struct {
	archived.Model
	Subtasks []string `gorm:"type:json;serializer:json"`
}

func (script *tasksUsesJSON) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.TransformColumns(
		basicRes,
		script,
		"_devlake_tasks",
		[]string{"subtasks"},
		func(src *srcTaskSubtaskJSON20230731) (*dstTaskSubtaskJSON20230731, errors.Error) {
			dst := &dstTaskSubtaskJSON20230731{
				Model: src.Model,
			}
			if len(src.Subtasks) == 0 {
				return nil, nil
			}
			println("src.Subtask", string(src.Subtasks))
			errors.Must(json.Unmarshal(src.Subtasks, &dst.Subtasks))
			return dst, nil
		},
	)
}

func (*tasksUsesJSON) Version() uint64 {
	return 20230728162121
}

func (*tasksUsesJSON) Name() string {
	return "tasks uses json"
}
