package logger

import (
	"fmt"
	"io"
	"os"
	"storj.io/ditto/pkg/models"
)

type Logger interface {
	Log(string)
	LogE(error)
}

type lg struct {
	buffer 			 io.Writer
	errBuffer 		 io.Writer
	bucketDiffBuffer []models.DiffModel

	format, errFormat string
}

func (l *lg) Log(msg string) {
	if msg == "" {
		return
	}

	l.buffer.Write([]byte(fmt.Sprintf(l.format, msg)))
}

func (l *lg) LogE(err error) {
	if err == nil {
		return
	}

	l.errBuffer.Write([]byte(fmt.Sprintf(l.errFormat, err)))
}

func (l *lg) Write(b []byte) (n int, err error) {
	return l.buffer.Write(b)
}

var StdOutLogger = lg{
	buffer: os.Stdout,
	errBuffer: os.Stdout,
	format: "Log: %s\n",
	errFormat: "Err: %s\n",
}