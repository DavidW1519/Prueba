package archived

import (
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type GitlabMergeRequestComment struct {
	GitlabId        int `gorm:"primaryKey"`
	MergeRequestId  int `gorm:"index"`
	MergeRequestIid int `gorm:"comment:Used in API requests ex. /api/merge_requests/<THIS_IID>"`
	Body            string
	AuthorUsername  string `gorm:"type:varchar(255)"`
	AuthorUserId    int
	GitlabCreatedAt time.Time
	Resolvable      bool `gorm:"comment:Is or is not review comment"`
	archived.NoPKModel
}

func (GitlabMergeRequestComment) TableName() string {
	return "_tool_gitlab_merge_request_comments"
}
