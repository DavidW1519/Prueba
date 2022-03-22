package services

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/worker"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	cfg := config.GetConfig()
	db, err = worker.NewGormDb(cfg, logger.Global.Nested("db"))

	if err != nil {
		panic(err)
	}

	// load plugins
	err = worker.LoadPlugins(config.GetConfig().GetString("PLUGIN_DIR"), logger.Global.Nested("plugin"))
	if err != nil {
		panic(err)
	}

	// migrate framework tables
	err = db.AutoMigrate(
		&models.Task{},
		&models.Notification{},
		&models.Pipeline{},
	)
	if err != nil {
		panic(err)
	}

	// migrate data tables if run in standalone mode
	if cfg.GetBool("STAND_ALONE") {
		err = worker.MigrateDb(db)
		if err != nil {
			panic(err)
		}
	}

	// set all unfinished tasks to failed
	db.Model(&models.Task{}).Where("status = ?", models.TASK_RUNNING).Update("status", models.TASK_FAILED)

}
