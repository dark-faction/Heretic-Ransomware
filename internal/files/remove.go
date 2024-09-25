package files

import (
	"go.uber.org/zap"
	"os"
)

func Remove(path string) {
	err := os.Remove(path) // remove a single file
	if err != nil {
		zap.L().Info(err.Error())
	}
}
