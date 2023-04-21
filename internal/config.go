package internal

import (
	"github.com/padok-team/yatas-template/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

// This function decodes the plugin configuration sent by YATAS to retrieve the
// fields this plugin needs and returns a slice of account configurations.
func UnmarshalConfig(c *commons.Config) ([]FakeAccount, error) {
	var accounts []FakeAccount

	// YATAS sends all the plugin configs of `.yatas.yml`, we iterate over them
	for _, r := range c.PluginConfig {
		var tmpAccounts []FakeAccount
		// Boolean to keep track if we've found the config of the plugin we're
		// interested in
		pluginFound := false
		for key, value := range r {

			switch key {
			case "pluginName":
				if value == "template" {
					pluginFound = true
				}
			case "accounts":
				for _, v := range value.([]interface{}) {
					var account FakeAccount
					logger.Logger.Debug("Inspecting account", "account", v)
					for keyaccounts, valueaccounts := range v.(map[string]interface{}) {
						// Add in this switch-case the fields of the plugin config you
						// want to read from
						switch keyaccounts {
						case "region":
							account.Region = valueaccounts.(string)
						}
					}
					tmpAccounts = append(tmpAccounts, account)
				}
			}
		}
		if pluginFound {
			logger.Logger.Debug("template config found ✅")
			accounts = tmpAccounts
		}
	}
	logger.Logger.Debug("Unmarshal Done ✅")
	logger.Logger.Debug("All accounts", "accounts", accounts)
	logger.Logger.Debug("Length of accounts", "len", len(accounts))
	return accounts, nil
}
