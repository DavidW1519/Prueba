/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

// RunTask FIXME ...
func RunTask(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	ctx context.Context,
	progress chan core.RunningProgress,
	taskId uint64,
) error {
	task := &models.Task{}
	err := db.Find(task, taskId).Error
	if err != nil {
		return err
	}
	if task.Status == models.TASK_COMPLETED {
		return fmt.Errorf("invalid task status")
	}
	beganAt := time.Now()
	// make sure task status always correct even if it panicked
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("run task failed with panic (%s): %v", utils.GatherCallFrames(0), r)
		}
		finishedAt := time.Now()
		spentSeconds := finishedAt.Unix() - beganAt.Unix()
		if err != nil {
			subTaskName := ""
			if pluginErr, ok := err.(*errors.SubTaskError); ok {
				subTaskName = pluginErr.GetSubTaskName()
			}
			dbe := db.Model(task).Updates(map[string]interface{}{
				"status":          models.TASK_FAILED,
				"message":         err.Error(),
				"finished_at":     finishedAt,
				"spent_seconds":   spentSeconds,
				"failed_sub_task": subTaskName,
			}).Error
			if dbe != nil {
				logger.Error("failed to finalize task status into db: %w", err)
			}
		} else {
			err = db.Model(task).Updates(map[string]interface{}{
				"status":        models.TASK_COMPLETED,
				"message":       "",
				"finished_at":   finishedAt,
				"spent_seconds": spentSeconds,
			}).Error
		}
	}()

	// start execution
	logger.Info("start executing task: %d", task.ID)
	err = db.Model(task).Updates(map[string]interface{}{
		"status":   models.TASK_RUNNING,
		"message":  "",
		"began_at": beganAt,
	}).Error
	if err != nil {
		return err
	}

	var options map[string]interface{}
	err = json.Unmarshal(task.Options, &options)
	if err != nil {
		return err
	}
	var subtasks []string
	err = json.Unmarshal(task.Subtasks, &subtasks)
	if err != nil {
		return err
	}

	err = RunPluginTask(
		config.GetConfig(),
		logger.Nested(task.Plugin),
		db,
		ctx,
		task.Plugin,
		subtasks,
		options,
		progress,
	)
	return err
}

// RunPluginTask FIXME ...
func RunPluginTask(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	ctx context.Context,
	name string,
	subtasks []string,
	options map[string]interface{},
	progress chan core.RunningProgress,
) error {
	pluginMeta, err := core.GetPlugin(name)
	if err != nil {
		return err
	}
	pluginTask, ok := pluginMeta.(core.PluginTask)
	if !ok {
		return fmt.Errorf("plugin %s doesn't support PluginTask interface", name)
	}

	return RunPluginSubTasks(
		cfg,
		logger,
		db,
		ctx,
		name,
		subtasks,
		options,
		pluginTask,
		progress,
	)
}

// RunPluginSubTasks FIXME ...
func RunPluginSubTasks(
	cfg *viper.Viper,
	logger core.Logger,
	db *gorm.DB,
	ctx context.Context,
	name string,
	subtasks []string,
	options map[string]interface{},
	pluginTask core.PluginTask,
	progress chan core.RunningProgress,
) error {
	logger.Info("start plugin")
	// find out all possible subtasks this plugin can offer
	subtaskMetas := pluginTask.SubTaskMetas()
	subtasksFlag := make(map[string]bool)
	for _, subtaskMeta := range subtaskMetas {
		subtasksFlag[subtaskMeta.Name] = subtaskMeta.EnabledByDefault
	}
	/* subtasksFlag example
	subtasksFlag := map[string]bool{
		"collectProject": true,
		"convertCommits": true,
		...
	}
	*/

	// user specifies what subtasks to run
	if len(subtasks) != 0 {
		// decode user specified subtasks
		var specifiedTasks []string
		err := mapstructure.Decode(subtasks, &specifiedTasks)
		if err != nil {
			return err
		}
		if len(specifiedTasks) > 0 {
			// first, disable all subtasks
			for task := range subtasksFlag {
				subtasksFlag[task] = false
			}
			// second, check specified subtasks is valid and enable them if so
			for _, task := range specifiedTasks {
				if _, ok := subtasksFlag[task]; ok {
					subtasksFlag[task] = true
				} else {
					return fmt.Errorf("subtask %s does not exist", task)
				}
			}
		}
	}

	// make sure `Required` subtasks are always enabled
	for _, subtaskMeta := range subtaskMetas {
		if subtaskMeta.Required {
			subtasksFlag[subtaskMeta.Name] = true
		}
	}

	// calculate total step(number of task to run)
	steps := 0
	for _, enabled := range subtasksFlag {
		if enabled {
			steps++
		}
	}

	taskCtx := helper.NewDefaultTaskContext(cfg, logger, db, ctx, name, subtasksFlag, progress)
	if closeablePlugin, ok := pluginTask.(core.CloseablePluginTask); ok {
		defer closeablePlugin.Close(taskCtx)
	}
	taskData, err := pluginTask.PrepareTaskData(taskCtx, options)
	if err != nil {
		return err
	}
	taskCtx.SetData(taskData)

	// execute subtasks in order
	taskCtx.SetProgress(0, steps)
	i := 0
	for _, subtaskMeta := range subtaskMetas {
		subtaskCtx, err := taskCtx.SubTaskContext(subtaskMeta.Name)
		if err != nil {
			// sth went wrong
			return err
		}
		if subtaskCtx == nil {
			// subtask was disabled
			continue
		}

		// run subtask
		logger.Info("executing subtask %s", subtaskMeta.Name)
		i++
		if progress != nil {
			progress <- core.RunningProgress{
				Type:          core.SetCurrentSubTask,
				SubTaskName:   subtaskMeta.Name,
				SubTaskNumber: i,
			}
		}
		err = subtaskMeta.EntryPoint(subtaskCtx)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: subtaskMeta.Name,
				Message:     err.Error(),
			}
		}
		taskCtx.IncProgress(1)
	}

	return nil
}

// UpdateProgressDetail FIXME ...
func UpdateProgressDetail(db *gorm.DB, taskId uint64, progressDetail *models.TaskProgressDetail, p *core.RunningProgress) {
	task := &models.Task{}
	task.ID = taskId
	switch p.Type {
	case core.TaskSetProgress:
		progressDetail.TotalSubTasks = p.Total
		progressDetail.FinishedSubTasks = p.Current
	case core.TaskIncProgress:
		progressDetail.FinishedSubTasks = p.Current
		// TODO: get rid of db update
		pct := float32(p.Current) / float32(p.Total)
		err := db.Model(task).Update("progress", pct).Error
		if err != nil {
			logger.Global.Error("failed to update progress: %w", err)
		}
	case core.SubTaskSetProgress:
		progressDetail.TotalRecords = p.Total
		progressDetail.FinishedRecords = p.Current
	case core.SubTaskIncProgress:
		progressDetail.FinishedRecords = p.Current
	case core.SetCurrentSubTask:
		progressDetail.SubTaskName = p.SubTaskName
		progressDetail.SubTaskNumber = p.SubTaskNumber
	}
}
