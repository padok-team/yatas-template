package internal

import (
	"github.com/padok-team/yatas-template/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

// This function decodes the plugin configuration sent by YATAS to retrieve the
// fields this plugin needs and returns a slice of account configurations.
func UnmarshalConfig(c *commons.Config) ([]FakeAccount, error) {
	var accounts []FakeAccount
	var pluginConfig map[string]interface{}

	// YATAS sends all the plugin configs of `.yatas.yml`, we iterate over them
	// and find the one that matches the name of the plugin
	logger.Logger.Debug("Searching for plugin config")
	for _, config := range c.PluginConfig {
		if config["pluginName"] == PluginName {
			logger.Logger.Debug("Plugin config found ✅")
			pluginConfig = config
		}
	}

	accountsConfig := pluginConfig["accounts"]
	// Iterate over the accounts associated to the plugin
	for _, acc := range accountsConfig.([]interface{}) {
		var account FakeAccount
		logger.Logger.Debug("Inspecting account", "account", acc)
		for key, value := range acc.(map[string]interface{}) {
			// TODO: Add in this switch-case the fields of the plugin config you
			// want to read from
			switch key {
			case "region":
				account.Region = value.(string)
			}
		}
		accounts = append(accounts, account)
	}

	logger.Logger.Debug("Unmarshal Done ✅")
	logger.Logger.Debug("All accounts", "accounts", accounts)
	logger.Logger.Debug("Length of accounts", "len", len(accounts))
	return accounts, nil
}
