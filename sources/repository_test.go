// andrzej lichnerowicz, unlicensed (~public domain)
package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestProvider struct{}

func (t *TestProvider) GetName() string            { return "test" }
func (t *TestProvider) CanHandle(string) bool  { return false }
func (t *TestProvider) Download(path string) error { return nil }

func TestRegisterProvider(t *testing.T) {
	RegisterProvider(&TestProvider{})
	_, ok := providers["test"]
	assert.True(t, ok)
}
