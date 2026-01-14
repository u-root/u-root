package defaultlog

import (
	"log"
	"os"
)

func Default() *log.Logger {
	return log.New(os.Stderr, "", 0)
}
