from typing import Iterable

import iso8601 as iso8601

from azure.api import AzureDevOpsAPI
from azure.models import GitRepository, GitCommit
from azure.streams.repositories import GitRepositories
from pydevlake import Substream, Stream, DomainType, Context
from pydevlake.domain_layer.code import Commit as DomainCommit
from pydevlake.domain_layer.code import RepoCommit as DomainRepoCommit


class GitCommits(Substream):
    tool_model = GitCommit
    domain_types = [DomainType.CODE]
    parent_stream = GitRepositories

    def collect(self, state, context, parent: GitRepository) -> Iterable[tuple[object, dict]]:
        connection = context.connection
        options = context.options
        azure_api = AzureDevOpsAPI(connection.base_url, connection.pat)
        # grab this info off the parent results
        response = azure_api.commits(options["org"], options["project"], parent.id)
        for raw_commit in response:
            raw_commit["repo_id"] = parent.id
            yield raw_commit, state

    def extract(self, raw_data: dict) -> GitCommit:
        return extract_raw_commit(raw_data)

    def convert(self, commit: GitCommit, ctx: Context) -> Iterable[DomainCommit]:
        yield DomainCommit(
            sha=commit.commit_sha,
            additions=commit.additions,
            deletions=commit.deletions,
            message=commit.comment,
            author_name=commit.author_name,
            author_email=commit.author_email,
            authored_date=commit.authored_date,
            author_id=commit.author_name,
            committer_name=commit.committer_name,
            committer_email=commit.committer_email,
            committed_date=commit.commit_date,
            committer_id=commit.committer_name,
        )

        yield DomainRepoCommit(
                repo_id=commit.repo_id,
                commit_sha=commit.commit_sha,
        )


def extract_raw_commit(stream: Stream, raw_data: dict, ctx: Context) -> GitCommit:
    commit: GitCommit = stream.tool_model(**raw_data)
    commit.project_id = ctx.options["project"]
    commit.repo_id = raw_data["repo_id"]
    commit.commit_sha = raw_data["commitId"]
    commit.author_name = raw_data["author"]["name"]
    commit.author_email = raw_data["author"]["email"]
    commit.authored_date = iso8601.parse_date(raw_data["author"]["date"])
    commit.committer_name = raw_data["committer"]["name"]
    commit.committer_email = raw_data["committer"]["email"]
    commit.commit_date = iso8601.parse_date(raw_data["committer"]["date"])
    if "changeCounts" in raw_data:
        commit.additions = raw_data["changeCounts"]["Add"]
        commit.deletions = raw_data["changeCounts"]["Delete"]
    return commit
