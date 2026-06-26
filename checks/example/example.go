package example

import (
	"sync"

	"github.com/padok-team/yatas-template/internal"
	"github.com/padok-team/yatas-template/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

// RunChecks is the entrypoint of the category (≈ yatas-aws's s3.RunChecks). It
// fetches the data once, fans out every check, collects the results and sends
// them back through queue. main.go registers it in initTest.
//
// To add a check: write its `checkIf...` function in its own file, give it a new
// incremented ID, and add one `go commons.CheckTest(...)` line below.
func RunChecks(wg *sync.WaitGroup, s internal.Session, c *commons.Config, queue chan []commons.Check) {
	logger.Logger.Debug("Example - Checks started")

	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	resources := GetResources(s)

	// CheckTest honours the user's include/exclude config (a disabled check is a
	// no-op) and registers the check on the wait group.
	go commons.CheckTest(checkConfig.Wg, c, ExampleCompliantID, checkIfResourceCompliant)(checkConfig, resources, ExampleCompliantID)

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
	logger.Logger.Debug("Example - Checks done")
}
