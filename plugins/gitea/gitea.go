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

package main

import (
	"github.com/apache/incubator-devlake/plugins/gitea/impl"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

var PluginEntry impl.Gitea //nolint

// Gitea API Docs: https://docs.gitea.io/zh-cn/api-usage/
func main() {
	cmd := &cobra.Command{Use: "gitea"}
	connectionId := cmd.Flags().Uint64P("connectionId", "c", 0, "gitea connection id")
	owner := cmd.Flags().StringP("owner", "o", "", "gitea owner")
	repo := cmd.Flags().StringP("repo", "r", "", "gitea repo")
	_ = cmd.MarkFlagRequired("connectionId")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"connectionId": *connectionId,
			"owner":        *owner,
			"repo":         *repo,
		})
	}
	runner.RunCmd(cmd)
}
