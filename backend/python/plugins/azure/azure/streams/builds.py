import typing
from typing import Iterable

import iso8601 as iso8601

from azure.api import AzureDevOpsAPI
from azure.helper import db
from azure.models import AzureDevOpsConnection, GitRepository
from azure.models import Build
from pydevlake import Context, DomainType, Stream, logger
from pydevlake.domain_layer.devops import *


class Builds(Stream):
    tool_model = Build
    domain_types = [DomainType.CICD]

    def collect(self, state, context) -> Iterable[tuple[object, dict]]:
        connection: AzureDevOpsConnection = context.connection
        options = context.options
        azure_api = AzureDevOpsAPI(connection.base_url, connection.pat)
        # grab this info off the parent results
        cached_repos = dict()
        response = azure_api.builds(options["org"], options["project"])
        for raw_build in response:
            if self.validate_repo(context, raw_build, cached_repos):
                yield raw_build, state

    def extract(self, raw_data: dict) -> Build:
        build: Build = self.tool_model(**raw_data)
        build.id = raw_data["id"]
        build.project_id = raw_data["project"]["id"]
        build.repo_id = raw_data["repository"]["id"]
        build.repo_url = raw_data["repository"]["url"]
        build.source_branch = raw_data["sourceBranch"]
        build.source_version = raw_data["sourceVersion"]
        build.build_number = raw_data["buildNumber"]
        if "buildNumberRevision" in raw_data:
            build.build_number_revision = raw_data["buildNumberRevision"]
        build.start_time = iso8601.parse_date(raw_data["startTime"])
        build.finish_time = iso8601.parse_date(raw_data["finishTime"])
        build.status = raw_data["status"]
        build.tags = ",".join(raw_data["tags"])
        build.priority = raw_data["priority"]
        build.build_result = raw_data["result"]
        trigger_info: dict = raw_data["triggerInfo"]
        if "ci.sourceSha" in trigger_info: # this key is not guaranteed to be in here per docs
            assert build.source_version == trigger_info["ci.sourceSha"]
        return build

    def convert(self, b: Build, ctx: Context) -> Iterable[DomainModel]:
        yield CICDPipeline(
                name=b.id,
                status=b.status,
                created_date=b.start_time,
                finished_date=b.finish_time,
                result=b.build_result.value,
                duration_sec=abs(b.finish_time.second-b.start_time.second),
                environment=CICDEnvironment.PRODUCTION.value,
                type=CICDType.DEPLOYMENT.value,
                cicd_scope_id=b.repo_id,
        )
        yield CiCDPipelineCommit(
                pipeline_id=b.id,
                commit_sha=b.source_version,
                branch=b.source_branch,
                repo_id=b.repo_id,
                repo=b.repo_url,
        )

    # workaround because azure also returns builds for unmanaged repos (we don't want them)
    @classmethod
    def validate_repo(cls, context: Context, raw_build: dict, cached_repos: typing.Dict[str, GitRepository]) -> bool:
        repo_id = raw_build["repository"]["id"]
        if repo_id not in cached_repos:
            repo: GitRepository = db.get(context, GitRepository, GitRepository.id == repo_id)
            if repo is None:
                logger.warn(f"no Azure repo associated with {repo_id}")
            cached_repos[repo_id] = repo
        if cached_repos[repo_id] is None:
            return False
        raw_build["repository"]["url"] = cached_repos[repo_id].url
        return True
