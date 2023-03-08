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

from typing import Optional
import base64

from pydevlake.api import API, request_hook, Paginator, Request


class AzurePaginator(Paginator):
    def get_items(self, response) -> Optional[list[object]]:
        return response.json['value']

    def get_next_page_id(self, response) -> Optional[int | str]:
        return response.headers.get('x-ms-continuation')

    def set_next_page_param(self, request, next_page_id):
        request.query_args['continuationToken'] = next_page_id


class AzureDevOpsAPI(API):
    paginator = AzurePaginator()

    def __init__(self, base_url: str, pat: str):
        self._base_url = base_url or "https://dev.azure.com/"
        self.pat = pat

    @property
    def base_url(self):
        return self._base_url

    @request_hook
    def authenticate(self, request: Request):
        if self.pat:
            pat_b64 = base64.b64encode((':' + self.pat).encode()).decode()
            request.headers['Authorization'] = 'Basic ' + pat_b64

    @request_hook
    def set_api_version(self, request: Request):
        request.query_args['api-version'] = "7.0"

    def my_profile(self):
        req = Request('https://app.vssps.visualstudio.com/_apis/profile/profiles/me')
        return self.send(req)

    def accounts(self, member_id: str):
        req = Request('https://app.vssps.visualstudio.com/_apis/accounts', query_args={"memberId": member_id})
        return self.send(req)

    def orgs(self) -> list[str]:
        response = self.accounts()
        return [acct["AccountName"] for acct in response.json]

    def projects(self, org: str):
        return self.get(org, '_apis/projects')

    # Get a project
    def project(self, org: str, project: str):
        return self.get(org, '_apis/projects', project)

    # List repos under an org
    def git_repos(self, org: str, project: str):
        return self.get(org, project, '_apis/git/repositories')

    def git_repo_pull_requests(self, org: str, project: str, repo_id: str):
        # see https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests?view=azure-devops-rest-7.1&tabs=HTTP
        return self.get(org, project, '_apis/git/repositories', repo_id, 'pullrequests')

    def git_repo_pull_request_commits(self, org: str, project: str, repo_id: str, pull_request_id: int):
        return self.get(org, project, '_apis/git/repositories', repo_id, 'pullRequests', pull_request_id, 'commits')

    def git_repo_pull_request_comments(self, org: str, project: str, repo_id: str, pull_request_id: int):
        return self.get(org, project, '_apis/git/repositories', repo_id, 'pullRequests', pull_request_id, 'threads')

    # not needed
    def commits(self, org: str, project: str, repo_id: str):
        return self.get(org, project, '_apis/git/repositories', repo_id, 'commits')

    def builds(self, org: str, project: str):
        return self.get(org, project, '_apis/build/builds')

    def jobs(self, org: str, project: str, build_id: int):
        return self.get(org, project, '_apis/build/builds', build_id, 'timeline')

    # unused
    def deployments(self, org: str, project: str):
        return self.get(org, project, '_apis/release/deployments')

    # unused
    def releases(self, org: str, project: str):
        return self.get(org, project, '_apis/release/releases')
