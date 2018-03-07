// andrzej lichnerowicz, unlicensed (~public domain)

package sources

// Provider is an interface that must be implemented
// by source handler
type Provider interface {
	GetName() string
	CanHandle(uri string) bool
	Download(path string) error
}

// Source is a generic structure for holding source meta-information
// Needs to implement Provider interface.
type Source struct {
	uri      string
	provider string
}

// RegisterProvider adds provider to global map
func RegisterProvider(provider Provider) {
	providers[provider.GetName()] = provider
}

var providers map[string]Provider

func init() {
	providers = make(map[string]Provider)
}
