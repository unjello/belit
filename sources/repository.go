// andrzej lichnerowicz, unlicensed (~public domain)

package sources

// Provider is an interface that must be implemented
// by source handler
type Provider interface {
	GetName() string
	CanHandle(uri string) bool
	Download(path string, debug bool, handlers ...interface{}) error
}

// RegisterProvider adds provider to global map
func RegisterProvider(provider Provider) {
	providers[provider.GetName()] = provider
}

var providers = initMap()

func initMap() map[string]Provider {
	return make(map[string]Provider)
}

type remoteHandler interface {
	Download(url string, dir string, debug bool) error
	Open(path string) error
}