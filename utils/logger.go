package utils

import (
	"bytes"
	"io"
	"log"
	"os"
)

func Logger(isVerbose bool) *log.Logger {
	out := io.Discard
	if isVerbose {
		out = io.Writer(os.Stderr)
	}

	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "", log.Ltime)
	)

	logger.SetOutput(out)

	return logger
}
