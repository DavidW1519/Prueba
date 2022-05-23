package app

import (
	"bytes"

	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func loadResources(configJson []byte) (*viper.Viper, core.Logger, *gorm.DB, error) {
	// prepare
	cfg := viper.New()
	cfg.SetConfigType("json")
	err := cfg.ReadConfig(bytes.NewBuffer(configJson))
	if err != nil {
		return nil, nil, nil, err
	}
	// TODO: should be redirected to server
	logger := logger.Global.Nested("worker")
	db, err := runner.NewGormDb(cfg, logger)
	if err != nil {
		return nil, nil, nil, err
	}
	return cfg, logger, db, err
}
