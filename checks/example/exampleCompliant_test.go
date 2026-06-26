package example

import (
	"sync"
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
)

// Checks are pure, so they test with hand-built data and no mocks. Cover an
// "all pass" and a "has a failure" case.
func Test_checkIfResourceCompliant(t *testing.T) {
	want := map[bool]string{true: "OK", false: "FAIL"}
	for compliant, status := range want {
		checkConfig := commons.CheckConfig{Queue: make(chan commons.Check, 1), Wg: &sync.WaitGroup{}}
		resources := []Resource{{ID: "test", Compliant: compliant}}

		checkIfResourceCompliant(checkConfig, resources, ExampleCompliantID)

		checkConfig.Wg.Add(1)
		go func() {
			for check := range checkConfig.Queue {
				if check.Status != status {
					t.Errorf("checkIfResourceCompliant(compliant=%v) = %v, want %s", compliant, check.Status, status)
				}
				checkConfig.Wg.Done()
			}
		}()
		checkConfig.Wg.Wait()
	}
}
