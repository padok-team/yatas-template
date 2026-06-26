package example

import "github.com/padok-team/yatas/plugins/commons"

// ExampleCompliantID is the check's unique, stable ID (<PLUGIN>_<CATEGORY>_<NNN>).
// Users reference it in `.yatas.yml`, so never reuse or renumber it.
const ExampleCompliantID = "TEMPLATE_EXAMPLE_001"

// checkIfResourceCompliant turns fetched resources into OK/FAIL results. A check
// never calls an API itself — that is the getter's job.
func checkIfResourceCompliant(checkConfig commons.CheckConfig, resources []Resource, testName string) {
	var check commons.Check
	check.InitCheck("Resources are compliant", "Check that each resource is compliant", testName, []string{"Security", "Good Practice"})

	for _, resource := range resources {
		if resource.Compliant {
			check.AddResult(commons.Result{Status: "OK", Message: "Resource " + resource.ID + " is compliant", ResourceID: resource.ID})
		} else {
			check.AddResult(commons.Result{Status: "FAIL", Message: "Resource " + resource.ID + " is not compliant", ResourceID: resource.ID})
		}
	}

	checkConfig.Queue <- check
}
