package runner

import (
	"context"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/spf13/cobra"
)

func RunCmd(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("tasks", "t", nil, "specify what tasks to run, --tasks=collectIssues,extractIssues")
	cmd.Execute()
}

func DirectRun(cmd *cobra.Command, args []string, pluginTask core.PluginTask, options map[string]interface{}) {
	tasks, err := cmd.Flags().GetStringSlice("tasks")
	if err != nil {
		panic(err)
	}
	options["tasks"] = tasks
	cfg := config.GetConfig()
	log := logger.Global.Nested(cmd.Use)
	db, err := NewGormDb(cfg, log)
	if err != nil {
		panic(err)
	}
	if pluginInit, ok := pluginTask.(core.PluginInit); ok {
		pluginInit.Init(cfg, log, db)
	}
	err = RunPluginSubTasks(
		cfg,
		log,
		db,
		context.Background(),
		cmd.Use,
		options,
		pluginTask,
		nil,
	)
	if err != nil {
		panic(err)
	}
}
