# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


from pydantic import BaseModel


class Message(BaseModel):
    pass


class SubtaskMeta(BaseModel):
    name: str
    entry_point_name: str
    required: bool
    enabled_by_default: bool
    description: str
    domain_types: list[str]
    arguments: list[str] = None


class PluginInfo(Message):
    name: str
    description: str
    connection_schema: dict
    transformation_rule_schema: dict
    plugin_path: str
    subtask_metas: list[SubtaskMeta]
    extension: str = "datasource"
    type: str = "python-poetry"


class SwaggerDoc(Message):
    name: str
    resource: str
    spec: dict


class PluginDetails(Message):
    plugin_info: PluginInfo
    swagger: SwaggerDoc


class RemoteProgress(Message):
    increment: int = 0
    current: int = 0
    total: int = 0


class Connection(Message):
    pass


class TransformationRule(Message):
    pass


class PipelineTask(Message):
    plugin: str
    # Do not snake_case this attribute,
    # it must match the json tag name in PipelineTask go struct
    skipOnFail: bool
    subtasks: list[str]
    options: dict[str, object]


class PipelineStage(Message):
    tasks: list[PipelineTask]


class PipelinePlan(Message):
    stages: list[PipelineStage]


class PipelineScope(Message):
    id: str
    name: str
    table_name: str


class BlueprintScope(Message):
    id: str
    name: str
