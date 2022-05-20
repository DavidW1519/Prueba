package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/helper"
)

type TapdStoryCommit struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ID           uint64 `gorm:"primaryKey;type:BIGINT" json:"id,string"`

	UserID          string `json:"user_id" gorm:"type:varchar(255)"`
	HookUserName    string `json:"hook_user_name" gorm:"type:varchar(255)"`
	CommitID        string `json:"commit_id" gorm:"type:varchar(255)"`
	WorkspaceID     uint64 `json:"workspace_id,string" gorm:"type:varchar(255)"`
	Message         string `json:"message"`
	Path            string `json:"path" gorm:"type:varchar(255)"`
	WebURL          string `json:"web_url" gorm:"type:varchar(255)"`
	HookProjectName string `json:"hook_project_name" gorm:"type:varchar(255)"`

	Ref        string          `json:"ref" gorm:"type:varchar(255)"`
	RefStatus  string          `json:"ref_status" gorm:"type:varchar(255)"`
	GitEnv     string          `json:"git_env" gorm:"type:varchar(255)"`
	FileCommit string          `json:"file_commit"`
	CommitTime *helper.CSTTime `json:"commit_time"`
	Created    *helper.CSTTime `json:"created"`

	StoryId uint64
	common.NoPKModel
}

func (TapdStoryCommit) TableName() string {
	return "_tool_tapd_story_commits"
}
