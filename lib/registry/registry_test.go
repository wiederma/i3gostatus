package registry

import (
	"github.com/pelletier/go-toml"
	"testing"
)

func TestInitialize(t *testing.T) {
	config := `modules = ["datetime", "temperature", "backlight"]`
	configTree, _ := toml.Load(config)
	enabledModules := Initialize(configTree)

	if n := len(enabledModules); n != 3 {
		t.Error("len(enabledModules) does not match!")
		t.Errorf("Expected: %d", n)
	}
}
