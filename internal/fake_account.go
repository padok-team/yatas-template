package internal

// FakeAccount is a fake struct for demo purposes that represents an account
// to connect to an API (e.g. AWS, GCP, Azure, etc.)
type FakeAccount struct {
	Region string `yaml:"region"` // A field from .yatas.yml plugin config
}
