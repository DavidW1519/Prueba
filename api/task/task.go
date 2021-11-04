package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/services"
)

func Post(ctx *gin.Context) {
	// We use a 2D array because the request body must be an array of a set of tasks
	// to be executed concurrently, while each set is to be executed sequentially.
	var data [][]services.NewTask

	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		logger.Error("", err)
		ctx.JSON(http.StatusBadRequest, "You must send down an array of objects")
		return
	}

	tasks := services.CreateTasksInDBFromJSON(data)
	// Return all created tasks to the User
	ctx.JSON(http.StatusCreated, tasks)

	go func() {
		err := services.RunAllTasks(data, tasks)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}()
}

func Get(ctx *gin.Context) {
	var query services.TaskQuery
	err := ctx.BindQuery(&query)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tasks, count, err := services.GetTasks(&query)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks, "count":count})
}

func Delete(ctx *gin.Context) {
	taskId := ctx.Param("taskId")
	id, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid task id")
		return
	}
	err = services.CancelTask(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	} else {
		err = models.Db.Model(&models.Task{}).Where("id = ?", id).Update("status", "CANCELLED").Error
		if err != nil {
			logger.Error("Could not upsert: ", err)
		}
	}
}
