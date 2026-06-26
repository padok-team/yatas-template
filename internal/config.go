package internal

import (
	"github.com/padok-team/yatas-template/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

// UnmarshalConfig decodes the plugin configuration YATAS reads from `.yatas.yml`
// into the list of targets to audit.
//
// YATAS hands every plugin's config to every plugin, so we first pick the block
// whose `pluginName` matches ours, then turn each entry under `accounts` into a
// Target. Extend the inner switch to read each field you add to Target.
func UnmarshalConfig(c *commons.Config) ([]Target, error) {
	var targets []Target

	// 1. Find this plugin's configuration block.
	var pluginConfig map[string]interface{}
	for _, config := range c.PluginConfig {
		if config["pluginName"] == PluginName {
			pluginConfig = config
		}
	}
	if pluginConfig == nil {
		logger.Logger.Error("No configuration found for plugin", "plugin", PluginName)
		return targets, nil
	}

	// 2. Read the list of targets to audit.
	accounts, ok := pluginConfig["accounts"].([]interface{})
	if !ok {
		logger.Logger.Error("No `accounts` list found in plugin config")
		return targets, nil
	}

	// 3. Decode each entry into a Target.
	for _, entry := range accounts {
		fields, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		var target Target
		for key, value := range fields {
			// TODO: add a case for every field you add to the Target struct.
			switch key {
			case "name":
				target.Name, _ = value.(string)
			// case "server":
			// 	target.Server, _ = value.(string)
			// case "token":
			// 	target.Token, _ = value.(string)
			default:
				logger.Logger.Warn("Unknown field in target config", "field", key)
			}
		}
		targets = append(targets, target)
	}

	logger.Logger.Debug("Targets to audit", "targets", targets)
	return targets, nil
}
