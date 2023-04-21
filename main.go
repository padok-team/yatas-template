package main

import (
	"encoding/gob"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/padok-team/yatas-template/internal"
	"github.com/padok-team/yatas-template/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

type YatasPlugin struct {
	logger hclog.Logger
}

// Don't remove this function.
// This function is called by YATAS through the RPC.
// This is the entrypoint of the plugin.
//
// It receives as argument the YATAS config of the user and returns the results
// of all the checks executed.
func (g *YatasPlugin) Run(c *commons.Config) []commons.Tests {
	// Set the global logger to the one used by the plugin
	logger.Logger = g.logger

	// Read the configuration sent by YATAS to retrieve a common account
	// configuration that you can use in your checks to make API calls to a
	// cloud provider for example
	var accounts []internal.FakeAccount
	var err error
	accounts, err = internal.UnmarshalConfig(c)
	if err != nil {
		logger.Logger.Error("Error unmarshaling accounts", "error", err)
		return nil
	}

	// TODO: Sent the `accounts` variable to the checks
	logger.Logger.Info("Accounts", "accounts", accounts)

	var checksAll []commons.Tests

	checks, err := runPlugin(c)
	if err != nil {
		logger.Logger.Error("Error running plugins", "error", err)
	}

	// TODO: Comment
	checksAll = append(checksAll, checks...)
	return checksAll
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// You do not need to change this function.
// This is the main entrypoint for the program, which launches the plugin RPC
// server.
func main() {
	// Register the types that will be used serialized and deserialized
	// and used in the RPC communication
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})

	// Here we setup the logger that will be used by the plugin and whose
	// output will be transmitted via RPC to YATAS
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	// Create an instance of YatasPlugin, which implements the Plugin interface
	// from hashicorp/go-plugin.
	yatasPlugin := &YatasPlugin{
		logger: logger,
	}

	// `pluginMap` is the map of plugins we can dispense.
	// Just this plugin in our case.
	var pluginMap = map[string]plugin.Plugin{
		internal.PluginName: &commons.YatasPlugin{Impl: yatasPlugin},
	}

	// Launch the plugin RPC server
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

// TODO: Refacto?
// Function that runs the checks or things to do.
func runPlugin(c *commons.Config) ([]commons.Tests, error) {
	var checksAll []commons.Tests

	// Run the checks here
	// TODO: placeholder?

	return checksAll, nil
}
