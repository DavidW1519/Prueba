package app

import (
	"context"

	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"go.temporal.io/sdk/activity"
)

func DevLakeTaskActivity(ctx context.Context, configJson []byte, taskId uint64) error {
	cfg, log, db, err := loadResources(configJson)
	if err != nil {
		return err
	}
	log.Info("received task #%d", taskId)
	progressDetail := &models.TaskProgressDetail{}
	progChan := make(chan core.RunningProgress)
	defer close(progChan)
	go func() {
		for p := range progChan {
			runner.UpdateProgressDetail(db, taskId, progressDetail, &p)
			activity.RecordHeartbeat(ctx, progressDetail)
		}
	}()
	err = runner.RunTask(cfg, log, db, ctx, progChan, taskId)
	if err != nil {
		log.Error("failed to execute task #%d: %w", taskId, err)
	}
	log.Info("finished task #%d", taskId)
	return err
}
