package syncthing

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "syncthing ", log.LstdFlags)
}
