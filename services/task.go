package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/utils"
	"github.com/merico-dev/lake/worker"
)

type RunningTask struct {
	mu    sync.Mutex
	tasks map[uint64]context.CancelFunc
}

func (rt *RunningTask) Add(taskId uint64, cancel context.CancelFunc) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if _, ok := rt.tasks[taskId]; ok {
		return fmt.Errorf("task with id %v already running", taskId)
	}
	rt.tasks[taskId] = cancel
	return nil
}

func (rt *RunningTask) Remove(taskId uint64) (context.CancelFunc, error) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if cancel, ok := rt.tasks[taskId]; ok {
		delete(rt.tasks, taskId)
		return cancel, nil
	}
	return nil, fmt.Errorf("task with id %v not found", taskId)
}

var runningTasks RunningTask

type TaskQuery struct {
	Status     string `form:"status"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	Plugin     string `form:"plugin"`
	PipelineId uint64 `form:"pipelineId" uri:"pipelineId"`
	Pending    int    `form:"pending"`
}

func init() {
	// set all previous unfinished tasks to status failed
	runningTasks.tasks = make(map[uint64]context.CancelFunc)
}

func CreateTask(newTask *models.NewTask) (*models.Task, error) {
	b, err := json.Marshal(newTask.Options)
	if err != nil {
		return nil, err
	}
	task := models.Task{
		Plugin:      newTask.Plugin,
		Options:     b,
		Status:      models.TASK_CREATED,
		Message:     "",
		PipelineId:  newTask.PipelineId,
		PipelineRow: newTask.PipelineRow,
		PipelineCol: newTask.PipelineCol,
	}
	err = db.Save(&task).Error
	if err != nil {
		logger.Error("save task failed", err)
		return nil, errors.InternalError
	}
	return &task, nil
}

func GetTasks(query *TaskQuery) ([]models.Task, int64, error) {
	db := db.Model(&models.Task{}).Order("id DESC")
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Plugin != "" {
		db = db.Where("plugin = ?", query.Plugin)
	}
	if query.PipelineId > 0 {
		db = db.Where("pipeline_id = ?", query.PipelineId)
	}
	if query.Pending > 0 {
		db = db.Where("finished_at is null")
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	if query.Page > 0 && query.PageSize > 0 {
		offset := query.PageSize * (query.Page - 1)
		db = db.Limit(query.PageSize).Offset(offset)
	}
	tasks := make([]models.Task, 0)
	err = db.Find(&tasks).Error
	if err != nil {
		return nil, count, err
	}
	return tasks, count, nil
}

func GetTask(taskId uint64) (*models.Task, error) {
	task := &models.Task{}
	err := db.Find(task, taskId).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func RunTask(taskId uint64) error {
	// load task information from database
	task, err := GetTask(taskId)
	if err != nil {
		return err
	}
	if task.Status != models.TASK_CREATED {
		return fmt.Errorf("invalid task status")
	}
	beganAt := time.Now()
	// make sure task status always correct even if it panicked
	defer func() {
		_, _ = runningTasks.Remove(task.ID)
		if r := recover(); r != nil {
			err = fmt.Errorf("run task failed with panic (%s): %v", utils.GatherCallFrames(), r)
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
				logger.Error("eror is not nil", err)
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

	// for task cancelling
	ctx, cancel := context.WithCancel(context.Background())
	err = runningTasks.Add(taskId, cancel)
	if err != nil {
		return err
	}
	// start execution
	logger.Info("start executing task ", task.ID)
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

	return worker.RunPluginTask(
		config.GetConfig(),
		logger.Global.Nested(task.Plugin),
		db,
		ctx,
		task.Plugin,
		options,
	)
}

func CancelTask(taskId uint64) error {
	cancel, err := runningTasks.Remove(taskId)
	if err != nil {
		return err
	}
	cancel()
	return nil
}
