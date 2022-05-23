package code

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

type Note struct {
	domainlayer.DomainEntity
	PrId        string `gorm:"index;comment:References the pull request for this note;type:varchar(100)"`
	Type        string `gorm:"type:varchar(100)"`
	Author      string `gorm:"type:varchar(255)"`
	Body        string
	Resolvable  bool `gorm:"comment:Is or is not a review comment"`
	IsSystem    bool `gorm:"comment:Is or is not auto-generated vs. human generated"`
	CreatedDate time.Time
}
