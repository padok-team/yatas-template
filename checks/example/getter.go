package example

import (
	"github.com/padok-team/yatas-template/internal"
	"github.com/padok-team/yatas-template/logger"
)

// Resource is a fake resource standing in for whatever this category audits (a
// Pod, an Application, a bucket, ...). Replace it with the type your API returns.
type Resource struct {
	ID        string
	Compliant bool
}

// GetResources fetches the resources to audit. All API/IO lives in getters so
// the checks stay pure and unit-testable. In a real plugin this uses the
// session's client; here it returns static data so the template runs as-is.
func GetResources(s internal.Session) []Resource {
	logger.Logger.Debug("Example - Fetching resources", "target", s.Target.Name)
	return []Resource{
		{ID: "resource-ok", Compliant: true},
		{ID: "resource-bad", Compliant: false},
	}
}
