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

package models

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetRepoIdMap(t *testing.T) {
	scopeConfig := &BambooScopeConfig{
		RepoMap: map[string][]int{
			"repo1": {1, 2},
			"repo2": {3},
			"repo3": {4, 5, 6},
		},
	}
	expectedRepoMap := map[int]string{
		1: "repo1",
		2: "repo1",
		3: "repo2",
		4: "repo3",
		5: "repo3",
		6: "repo3",
	}

	assert.Equal(t, scopeConfig.GetRepoIdMap(), expectedRepoMap, "Unexpected result")
}
