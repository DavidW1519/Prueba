package models

import "github.com/apache/incubator-devlake/models/common"

type JiraIssueCommit struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	IssueId      uint64 `gorm:"primaryKey"`
	CommitSha    string `gorm:"primaryKey;type:varchar(40)"`
	CommitUrl    string `gorm:"type:varchar(255)"`
}

func (JiraIssueCommit) TableName() string {
	return "_tool_jira_issue_commits"
}
