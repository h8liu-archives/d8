package printer

import (
	"io"
)

type noopWriter struct{}

var noop io.Writer = new(noopWriter)

func (w *noopWriter) Write(buf []byte) (int, error) {
	return len(buf), nil
}
