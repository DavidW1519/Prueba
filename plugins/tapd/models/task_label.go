package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type TapdTaskLabel struct {
	TaskId    uint64 `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (TapdTaskLabel) TableName() string {
	return "_tool_tapd_task_labels"
}
