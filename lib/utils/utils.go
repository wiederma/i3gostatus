package utils

import "encoding/json"

func Json(data interface{}) string {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(json)
}
