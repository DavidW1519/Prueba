package plugins

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/merico-dev/lake/config"

	"github.com/merico-dev/lake/logger"

	. "github.com/merico-dev/lake/plugins/core"
)

// LoadPlugins load plugins from local directory
func LoadPlugins(pluginsDir string) error {
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		println(fileName, path)
		if strings.HasSuffix(fileName, ".so") {
			pluginName := fileName[0 : len(d.Name())-3]
			plug, loadErr := plugin.Open(path)
			if loadErr != nil {
				return loadErr
			}
			symPluginEntry, pluginEntryError := plug.Lookup("PluginEntry")
			if pluginEntryError != nil {
				return pluginEntryError
			}
			plugEntry, ok := symPluginEntry.(Plugin)
			if !ok {
				return fmt.Errorf("%v PluginEntry must implement Plugin interface", pluginName)
			}
			plugEntry.Init()
			logger.Info(`[plugins] init a plugin success`, pluginName)
			err = RegisterPlugin(pluginName, plugEntry)
			if err != nil {
				return nil
			}
			logger.Info("[plugins] plugin loaded", pluginName)
		}
		return nil
	})
	return walkErr
}

func RunPlugin(taskId uint64, name string, options map[string]interface{}, progress chan<- float32) error {
	plugin, err := GetPlugin(name)
	if err != nil {
		return err
	}
	plugin.Execute(options, taskId, progress)
	return nil
}

func PluginDir() string {
	pluginDir := config.V.GetString("PLUGIN_DIR")
	if !path.IsAbs(pluginDir) {
		wd := config.V.GetString("WORKING_DIRECTORY")
		pluginDir = filepath.Join(wd, pluginDir)
	}
	return pluginDir
}
