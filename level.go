package logutils

import (
	"bytes"
	"io"
	"sync"
)

type LogLevel string

// LevelFilter is an io.Writer that can be used with a logger that
// will filter out log messages that aren't at least a certain level.
//
// Once the filter is in use somewhere, it is not safe to modify
// the structure.
type LevelFilter struct {
	// Levels is the list of log levels, in increasing order of
	// severity. Example might be: {"DEBUG", "WARN", "ERROR"}.
	Levels []LogLevel

	// MinLevel is the minimum level allowed through
	MinLevel LogLevel

	// The underlying io.Writer where log messages that pass the filter
	// will be set.
	Writer io.Writer

	badLevels map[LogLevel]struct{}
	once      sync.Once
}

// Check will check a given line if it would be included in the level
// filter.
func (f *LevelFilter) Check(line []byte) bool {
	// Check for a log level
	var level LogLevel
	x := bytes.IndexByte(line, '[')
	if x >= 0 {
		y := bytes.IndexByte(line[x:], ']')
		if y >= 0 {
			level = LogLevel(line[x+1 : y])
		}
	}

	_, ok := f.badLevels[level]
	return !ok
}

func (f *LevelFilter) Write(p []byte) (n int, err error) {
	f.once.Do(f.init)

	// Note in general that io.Writer can receive any byte sequence
	// to write, but the "log" package always guarantees that we only
	// get a single line. We use that as a slight optimization within
	// this method, assuming we're dealing with a single, complete line
	// of log data.

	if !f.Check(p) {
		return len(p), nil
	}

	return f.Writer.Write(p)
}

func (f *LevelFilter) init() {
	f.badLevels = make(map[LogLevel]struct{})
	for _, level := range f.Levels {
		if level == f.MinLevel {
			break
		}

		f.badLevels[level] = struct{}{}
	}
}
