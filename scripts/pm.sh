#!/bin/sh
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

LAKE_ENDPOINT=${LAKE_ENDPOINT-'http://localhost:8080'}
LAKE_PIPELINE_URL=$LAKE_ENDPOINT/pipelines

debug() {
    $SCRIPT_DIR/compile-plugins.sh -gcflags=all="-N -l"
    dlv debug
}

run() {
    $SCRIPT_DIR/compile-plugins.sh
    go run $SCRIPT_DIR/../main.go
}

jira_source_post() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/connections" --data '
    {
        "name": "test-jira-connection",
        "endpoint": "'"$JIRA_ENDPOINT"'",
        "basicAuthEncoded": "'"$JIRA_BASIC_AUTH_ENCODED"'",
        "epicKeyField": "'"$JIRA_ENDPOINT"'",
        "storyPointField": "'"$JIRA_ISSUE_STORYPOINT_FIELD"'",
    }
    ' | jq
}

jira_source_post_full() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/connections" --data '
    {
        "name": "test-jira-connection",
        "endpoint": "'"$JIRA_ENDPOINT"'",
        "basicAuthEncoded": "'"$JIRA_BASIC_AUTH_ENCODED"'",
        "epicKeyField": "'"$JIRA_ENDPOINT"'",
        "storyPointField": "'"$JIRA_ISSUE_STORYPOINT_FIELD"'",
        "typeMappings": {
            "Story": {
                "standardType": "Requirement",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    },
                    "已解决": {
                        "standardStatus": "Resolved"
                    }
                }
            },
            "Incident": {
                "standardType": "Incident",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    }
                }
            },
            "Bug": {
                "standardType": "Bug",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    }
                }
            }
        }
    }' | jq
}

jira_source_post_fail() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/connections" --data @- <<'    JSON' | jq
    {
        "name": "test-jira-connection-fail",
        "endpoint": "https://merico.atlassian.net/rest",
        "basicAuthEncoded": "basicAuth",
        "epicKeyField": "epicKeyField",
        "storyPointField": "storyPointField",
        "typeMappings": "ehhlow"
    }
    JSON
}

jira_source_put() {
    curl -v -XPUT "$LAKE_ENDPOINT/plugins/jira/connections/$1" --data @- <<'    JSON' | jq
    {
        "name": "test-jira-connection-updated",
        "endpoint": "https://merico.atlassian.net/rest",
        "basicAuthEncoded": "basicAuth",
        "epicKeyField": "epicKeyField",
        "storyPointField": "storyPointField",
    }
    JSON
}

jira_source_put_full() {
    curl -v -XPUT "$LAKE_ENDPOINT/plugins/jira/connections/$1" --data '
    {
        "name": "test-jira-connection-updated",
        "endpoint": "'"$JIRA_ENDPOINT"'",
        "basicAuthEncoded": "'"$JIRA_BASIC_AUTH_ENCODED"'",
        "epicKeyField": "'"$JIRA_ENDPOINT"'",
        "storyPointField": "'"$JIRA_ISSUE_STORYPOINT_FIELD"'",
        "typeMappings": {
            "Story": {
                "standardType": "Requirement",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    },
                    "已解决": {
                        "standardStatus": "Resolved"
                    }
                }
            },
            "Incident": {
                "standardType": "Incident",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    }
                }
            },
            "Bug": {
                "standardType": "Bug",
                "statusMappings": {
                    "已完成": {
                        "standardStatus": "Resolved"
                    }
                }
            }
        }
    }' | jq
}

jira_source_list() {
    curl -v "$LAKE_ENDPOINT/plugins/jira/connections" | jq
}

jira_source_get() {
    curl -v "$LAKE_ENDPOINT/plugins/jira/connections/$1" | jq
}

jira_source_delete() {
    curl -v -XDELETE "$LAKE_ENDPOINT/plugins/jira/connections/$1"
}

jira_typemapping_post() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/connections/$1/type-mappings" --data @- <<'    JSON' | jq
    {
        "userType": "userType",
        "standardType": "standardType"
    }
    JSON
}

jira_typemapping_put() {
    curl -v -XPUT "$LAKE_ENDPOINT/plugins/jira/connections/$1/type-mappings/$2" --data @- <<'    JSON' | jq
    {
        "standardType": "standardTypeUpdated"
    }
    JSON
}

jira_typemapping_delete() {
    curl -v -XDELETE "$LAKE_ENDPOINT/plugins/jira/connections/$1/type-mappings/$2"
}

jira_typemapping_list() {
    curl -v "$LAKE_ENDPOINT/plugins/jira/connections/$1/type-mappings" | jq
}

jira_statusmapping_post() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/connections/$1/type-mappings/$2/status-mappings" --data @- <<'    JSON' | jq
    {
        "userStatus": "userStatus",
        "standardStatus": "standardStatus"
    }
    JSON
}

jira_statusmapping_put() {
    curl -v -XPUT "$LAKE_ENDPOINT/plugins/jira/connections/$1/type-mappings/$2/status-mappings/$3" --data @- <<'    JSON' | jq
    {
        "standardStatus": "standardStatusUpdated"
    }
    JSON
}

jira_statusmapping_delete() {
    curl -v -XDELETE "$LAKE_ENDPOINT/plugins/jira/connections/$1/type-mappings/$2/status-mappings/$3"
}

jira_statusmapping_list() {
    curl -v "$LAKE_ENDPOINT/plugins/jira/connections/$1/type-mappings/$2/status-mappings" | jq
}

jira_echo() {
    curl -v -XPOST "$LAKE_ENDPOINT/plugins/jira/echo" --data @- <<'    JSON' | jq
    {
        "plugin": "jira",
        "options": {
            "boardId": 8
        }
    }
    JSON
}

pipeline_new() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON' | jq
    {
        "name": "test-all",
        "tasks": [
            [
                {
                    "plugin": "jira",
                    "options": {
                        "connectionId": 1,
                        "boardId": 8
                    }
                },
                {
                    "plugin": "jenkins",
                    "options": {}
                }
            ]
        ]
    }
    JSON
}

pipelines() {
    curl -v $LAKE_PIPELINE_URL'?'$1 | jq
}

pipeline() {
    curl -v $LAKE_PIPELINE_URL/$1 | jq
}

pipeline_cancel() {
    curl -v -XDELETE $LAKE_PIPELINE_URL/$1
}

pipeline_tasks() {
    curl -v $LAKE_PIPELINE_URL/$1/tasks'?'$2 | jq
}

jira() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data '
    {
        "name": "test-jira",
        "tasks": [
            [
                {
                    "plugin": "jira",
                    "options": {
                        "connectionId": '$1',
                        "boardId": '$2',
                        "tasks": ['"$3"']
                    }
                }
            ]
        ]
    }
    ' | jq
}

gitlab() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON'
    {
        "name": "test-gitlab",
        "tasks": [
            [
                {
                    "plugin": "gitlab",
                    "options": {
                        "projectId": 8967944,
                        "tasks": ["collectMrs"]
                    }
                }
            ]
        ]
    }
    JSON
}

github() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON'
    {
        "name": "test-github",
        "tasks": [
            [
                {
                    "plugin": "github",
                    "options": {
                        "repo": "lake",
                        "owner": "merico-dev",
                        "tasks": ["collectCommits"]
                    }
                }
            ]
        ]
    }
    JSON
}

jenkins() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON'
    {
        "name": "test-jenkins",
        "tasks": [
            [
                {
                    "plugin": "jenkins",
                    "options": {}
                }
            ]
        ]
    }
    JSON
}

ae() {
    curl -v -XPOST $LAKE_PIPELINE_URL --data @- <<'    JSON'
    {
        "name": "test-ae",
        "tasks": [
            [
                {
                    "plugin": "ae",
                    "options": {
                        "projectId": 13
                    }
                }
            ]
        ]
    }
    JSON
}

truncate() {
    SQL=$()
    echo "SET FOREIGN_KEY_CHECKS=0;"
    echo 'show tables' | mycli local-lake | tail -n +2 | xargs -I{} -n 1 echo "truncate table {};"
    echo "SET FOREIGN_KEY_CHECKS=1;"
}

lint() {
    golangci-lint run -v
}

"$@"
