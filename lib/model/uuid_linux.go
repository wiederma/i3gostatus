package model

import (
	"io/ioutil"
	"strings"
)

func uuidV4() string {
	uuid, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(uuid))
}
