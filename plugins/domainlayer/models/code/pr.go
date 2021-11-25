package code

import (
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Pr struct {
	base.DomainEntity
	RepoId      uint64 `gorm:"index"`
	State       string `gorm:"comment:open/closed or other"`
	Title       string
	Url         string
	CreatedDate time.Time
	MergedDate  *time.Time
	ClosedAt    *time.Time
}
