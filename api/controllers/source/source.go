package source

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/api/services"
	"github.com/merico-dev/lake/api/types"
	"github.com/merico-dev/lake/logger"
)

func Post(ctx *gin.Context) {
	var data types.CreateSource
	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	logger.Debug("Created Source", data)
	logger.Info("Created Source", data)
	logger.Error("Created Source", data)
	logger.Warn("Created Source", data)
	err = services.NewSource(data)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, source)
}
