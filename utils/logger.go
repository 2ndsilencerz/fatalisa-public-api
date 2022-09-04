package utils

import (
	"bytes"
	"fmt"
	"github.com/subchen/go-log"
	"github.com/subchen/go-log/formatters"
	"strconv"
	"sync"
	"time"
)

const (
	FileLogLocation = "logs/"
)

var (
	FileLogName = FileLogLocation + "gin-" + GetPodName() + ".log"
	fmtBuffer   = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

type LogFormat struct {
	TimeFormat string
	init       sync.Once
	isterm     bool
}

func (f *LogFormat) Format(level log.Level, msg string, logger *log.Logger) []byte {
	f.init.Do(func() {

		if f.TimeFormat == "" {
			f.TimeFormat = "2006/01/02 15:04:05 -0700"
		}

		f.isterm = formatters.IsTerminal(logger.Out)
	})

	buf := fmtBuffer.Get().(*bytes.Buffer)
	buf.Reset()
	defer fmtBuffer.Put(buf)

	// timestamp
	timeStr := time.Now().Format(f.TimeFormat)
	buf.WriteString(timeStr)

	// level
	buf.WriteByte(' ')
	if f.isterm {
		buf.WriteString(fmt.Sprintf("%-7s", "["+level.ColorString()+"]"))
	} else {
		buf.WriteString(fmt.Sprintf("%-7s", "["+level.String()+"]"))
	}

	// separator for message
	buf.WriteByte(' ')
	buf.WriteString(":")

	// file, line
	if level == log.ERROR {
		file, line := formatters.FilelineCaller(5)
		buf.WriteByte(' ')
		buf.WriteString(file)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
	}

	// msg
	buf.WriteByte(' ')
	buf.WriteString(msg)

	// newline
	buf.WriteByte('\n')

	return buf.Bytes()
}
