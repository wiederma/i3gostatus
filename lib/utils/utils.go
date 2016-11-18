package utils

import (
	"encoding/json"
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
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
		current = configTree.GetDefault(moduleStr+".period", int64(model.DefaultPeriod)).(int64)
		if current < smallest {
			smallest = current
		}
	}

	return time.Duration(smallest) * time.Millisecond
}
