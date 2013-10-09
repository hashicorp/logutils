package logutils

import (
	"bytes"
	"io"
)

type LogLevel string

// LevelFilter is an io.Writer that can be used with a logger that
// will filter out log messages that aren't at least a certain level.
type LevelFilter struct {
	// Levels is the list of log levels, in increasing order of
	// severity. Example might be: {"DEBUG", "WARN", "ERROR"}.
	Levels []LogLevel

	// MinLevel is the minimum level allowed through
	MinLevel LogLevel

	// The underlying io.Writer where log messages that pass the filter
	// will be set.
	Writer io.Writer
}

func (f *LevelFilter) Write(p []byte) (n int, err error) {
	// Check for a log level
	var level LogLevel
	x := bytes.IndexByte(p, '[')
	if x >= 0 {
		y := bytes.IndexByte(p[x:], ']')
		if y >= 0 {
			level = LogLevel(p[x+1:y])
		}
	}

	if level != "" {
		for _, l := range f.Levels {
			// If we reached a level we care about, skip it
			if l == f.MinLevel {
				break
			}

			if l == level {
				return len(p), nil
			}
		}
	}

	return f.Writer.Write(p)
}
