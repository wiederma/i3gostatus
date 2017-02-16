package model

import "io/ioutil"

func uuidV4() string {
	uuid, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		panic(err)
	}

	return string(uuid)
}
