package models

import (
	"github.com/merico-dev/lake/models"
)

type JiraProject struct {
	models.NoPKModel

	// collected fields
	SourceId uint64 `gorm:"primarykey"`
	Id       string `gorm:"primaryKey"`
	Key      string
	Name     string
}
