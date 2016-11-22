package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"math"
	"time"
)

func Json(data interface{}) string {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(json)
}

// TODO: Function name?
func FindFastestModule(configTree *toml.TomlTree) time.Duration {
	res := configTree.Get("modules").([]interface{})
	var smallest int64
	var current int64

	smallest = model.DefaultPeriod

	for _, module := range res {
		moduleStr := module.(string)
		// TODO
		current = configTree.GetDefault(moduleStr+".period", int64(model.DefaultPeriod)).(int64)
		if current < smallest {
			smallest = current
		}
	}

	return time.Duration(smallest) * time.Millisecond
}

func HumanReadableByteCount(x uint64) string {
	base := float64(1024)
	prefixes := []rune("kMGTPE")
	a := float64(x)

	if a < base {
		return fmt.Sprintf("%f B", a)
	}

	// https://en.wikipedia.org/wiki/Binary_prefix
	exp := math.Floor(math.Log2(a) / math.Log2(base))
	unit := string(prefixes[int(exp)-1])

	return fmt.Sprintf("%.0f %siB", a/math.Pow(base, exp), unit)
}
