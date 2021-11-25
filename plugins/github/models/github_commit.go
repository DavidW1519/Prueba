package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type GithubCommit struct {
	Sha            string `gorm:"primaryKey"`
	RepositoryId   int    `gorm:"index"`
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   time.Time
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
	Message        string
	Url            string
	Additions      int `gorm:"comment:Added lines of code"`
	Deletions      int `gorm:"comment:Deleted lines of code"`

	models.NoPKModel
}
