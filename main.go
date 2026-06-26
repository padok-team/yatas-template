package main

import (
	"encoding/gob"
	"os"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/padok-team/yatas-template/checks/example"
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
	targets, err := internal.UnmarshalConfig(c)
	if err != nil {
		logger.Logger.Error("Error unmarshaling config", "error", err)
		return nil
	}

	return run(c, targets)
}

// run audits every target concurrently. A queue + wait group collect the results
// so targets never block one another.
func run(c *commons.Config, targets []internal.Target) []commons.Tests {
	var wg sync.WaitGroup
	queue := make(chan commons.Tests, 10)
	var checks []commons.Tests

	wg.Add(len(targets))
	for _, target := range targets {
		go runTestsForTarget(target, c, queue)
	}

	go func() {
		for t := range queue {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()
	return checks
}

func runTestsForTarget(target internal.Target, c *commons.Config, queue chan commons.Tests) {
	session := initSession(target)
	queue <- initTest(session, c, target)
}

// initSession authenticates against the target and returns a ready-to-use client
// (≈ yatas-aws's initAuth). Build your real client here.
func initSession(target internal.Target) internal.Session {
	logger.Logger.Debug("Init session", "target", target.Name)
	// TODO: build and return your client, e.g. kubernetes.NewForConfig(...).
	return internal.Session{Target: target}
}

// initTest runs every check category for one target and aggregates their checks.
// Add one `go commons.CheckMacroTest(...)` line per category you create under
// checks/ (this is the only place to register a category).
func initTest(s internal.Session, c *commons.Config, target internal.Target) commons.Tests {
	var checks commons.Tests
	checks.Account = target.Name

	var wg sync.WaitGroup
	queue := make(chan []commons.Check, 100)

	go commons.CheckMacroTest(&wg, c, example.RunChecks)(&wg, s, c, queue)

	go func() {
		for t := range queue {
			checks.Checks = append(checks.Checks, t...)
			wg.Done()
		}
	}()

	wg.Wait()
	return checks
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
