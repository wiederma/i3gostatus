package utils

import (
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	expected := `{"full_text":"Hans","short_text":"Hans Short","color":"#ff0000","background":"#f00000","name":"HansBlock","instance":"eolfj09209ijkmn,2qrw","urgent":true,"separator":true}`
	i3BarBlock := model.I3BarBlock{
		FullText:   "Hans",
		ShortText:  "Hans Short",
		Color:      "#ff0000",
		Background: "#f00000",
		Name:       "HansBlock",
		Instance:   "eolfj09209ijkmn,2qrw",
		Urgent:     true,
		Separator:  true,
	}

	if jsonStr := Json(i3BarBlock); jsonStr != expected {
		t.Errorf("Wrong json string: %s", jsonStr)
	}
}

func TestFindFastestModule(t *testing.T) {
	expected := 1000 * time.Millisecond
	config := `modules = ["datetime", "temperature"]`
	configTree, _ := toml.Load(config)

	if res := FindFastestModule(configTree); res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}
}
