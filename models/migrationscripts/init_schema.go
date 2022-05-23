package migrationscripts

import (
	"context"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type initSchemas struct{}

func (*initSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.Task{},
		&archived.Notification{},
		&archived.Pipeline{},
		&archived.Blueprint{},
		&archived.User{},
		&archived.Repo{},
		&archived.Commit{},
		&archived.CommitParent{},
		&archived.PullRequest{},
		&archived.PullRequestCommit{},
		&archived.PullRequestComment{},
		&archived.PullRequestLabel{},
		&archived.Note{},
		&archived.RepoCommit{},
		&archived.Ref{},
		&archived.RefsCommitsDiff{},
		&archived.CommitFile{},
		&archived.Board{},
		&archived.Issue{},
		&archived.IssueLabel{},
		&archived.IssueComment{},
		&archived.BoardIssue{},
		&archived.BoardSprint{},
		&archived.Changelog{},
		&archived.Sprint{},
		&archived.SprintIssue{},
		&archived.IssueStatusHistory{},
		&archived.IssueSprintsHistory{},
		&archived.IssueAssigneeHistory{},
		&archived.Job{},
		&archived.Build{},
		&archived.IssueWorklog{},
		&archived.BoardRepo{},
		&archived.PullRequestIssue{},
		&archived.IssueCommit{},
		&archived.IssueRepoCommit{},
		&archived.RefsIssuesDiffs{},
		&archived.RefsPrCherrypick{},
	)
}

func (*initSchemas) Version() uint64 {
	return 20220406212344
}

func (*initSchemas) Owner() string {
	return "Framework"
}

func (*initSchemas) Name() string {
	return "create init schemas"
}
