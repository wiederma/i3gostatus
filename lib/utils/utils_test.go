package utils

import (
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"strings"
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

func TestHumanReadableByteCount(t *testing.T) {
	expected := "408 GiB"
	res := HumanReadableByteCount(uint64(437875942755))
	if strings.Compare(res, expected) != 0 {
		t.Errorf("Fail; expected: %s", expected)
	}

	expected = "500 B"
	res = HumanReadableByteCount(uint64(500))
	if strings.Compare(res, expected) != 0 {
		t.Logf("Wrong result: %s", res)
		t.Errorf("Expected: %s", expected)
	}
}

func TestWhich(t *testing.T) {
	if cmd, err := Which("ls"); err == nil {
		t.Logf("cmd %s found", cmd)
	} else {
		t.Log("ls is not available?")
		t.Errorf("failed with error: %s", err)
	}

	if _, err := Which("kalsdfjlsajf"); err == nil {
		t.Errorf("Found non existing command!")
	} else {
		t.Logf("correctly reported error: %s", err)
	}
}
