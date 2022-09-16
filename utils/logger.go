package utils

import (
	"bytes"
	"io"
	"log"
	"os"
)

func Logger(isDiscard bool) *log.Logger {
	out := io.Writer(os.Stderr)
	if isDiscard {
		out = io.Discard
	}

	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "", log.Ltime)
	)

	logger.SetOutput(out)

	return logger
}
