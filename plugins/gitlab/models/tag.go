package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type GitlabTag struct {
	Name               string `gorm:"primaryKey;type:varchar(60)"`
	Message            string
	Target             string `gorm:"type:varchar(255)"`
	Protected          bool
	ReleaseDescription string
	common.NoPKModel
}

func (GitlabTag) TableName() string {
	return "_tool_gitlab_tags"
}
