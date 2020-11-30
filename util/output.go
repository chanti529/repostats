package util

import "github.com/jfrog/jfrog-client-go/utils/log"

type LogIoWriter struct {
	buffer []byte
}

func (w *LogIoWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if b != '\n' {
			w.buffer = append(w.buffer, b)
		}
		if b == '\n' {
			log.Output(string(w.buffer))
			w.buffer = make([]byte, 128)
		}
	}
	return len(p), nil
}
