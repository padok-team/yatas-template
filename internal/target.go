package internal

// Target is one thing to audit: a Kubernetes cluster, an ArgoCD instance, a
// cloud account, ... A plugin built from this template audits a list of them.
//
// Rename it to fit your plugin and add the connection fields you need (this is
// the equivalent of yatas-aws's AWS_Account). They are read from `.yatas.yml`
// by UnmarshalConfig.
type Target struct {
	Name string `yaml:"name"` // Name of the target, shown in the reports
	// TODO: add the fields your plugin connects with, e.g.:
	// Server  string `yaml:"server"`  // API server URL
	// Token   string `yaml:"token"`   // API token
	// Context string `yaml:"context"` // kubeconfig context
}

// Session is the authenticated client used to query the API (a
// *kubernetes.Clientset, an ArgoCD API client, an aws.Config, ...). It is built
// once per target by initSession in main.go and passed to the getters.
//
// Use a Session when connecting is non-trivial or shared: it authenticates once
// (auth handshake, token exchange, SDK client) and every getter reuses it.
//
// You don't need a Session when there is no client to authenticate or reuse —
// e.g. checking local files, parsing manifests, or hitting a public endpoint. In
// that case delete this type and initSession, and have getters take the Target
// directly (GetResources(t Target)).
//
// Replace the placeholder with your real client.
type Session struct {
	Target Target
	// TODO: add your client, e.g. Client *kubernetes.Clientset
}
