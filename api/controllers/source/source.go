package source

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/api/services"
	"github.com/merico-dev/lake/api/types"
	"github.com/merico-dev/lake/logger"
)

// PostSource godoc
// @Summary create a source
// @Description create a source for plugin
// @ID create-post
// @Accept  json
// @Produce  json
// @Param param body types.CreateSource true "task info"
// @Success 200 {object} models.Source
// @Header 200 {string} Token "qwerty"
// @Router /source [post]
func Post(ctx *gin.Context) {
	var data types.CreateSource
	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	logger.Debug("display data", data)
	source, err := services.NewSource(data)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, source)
}
