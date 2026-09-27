// Package logutils augments the standard log package with levels.
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
//
// Level tags are matched against the configured Levels list (e.g. "[WARN]"),
// not the first pair of square brackets on the line. That way a log.SetPrefix
// value that itself contains brackets does not disable filtering (#8).
func (f *LevelFilter) Check(line []byte) bool {
	f.once.Do(f.init)

	level, ok := f.levelFromLine(line)
	if !ok {
		// No recognized level tag: keep previous behavior and allow the line.
		return true
	}
	_, bad := f.badLevels[level]
	return !bad
}

// levelFromLine finds a configured level tag like "[WARN]" in the line.
func (f *LevelFilter) levelFromLine(line []byte) (LogLevel, bool) {
	for _, lvl := range f.Levels {
		needle := make([]byte, 0, len(lvl)+2)
		needle = append(needle, '[')
		needle = append(needle, lvl...)
		needle = append(needle, ']')
		if bytes.Contains(line, needle) {
			return lvl, true
		}
	}
	return "", false
}

func (f *LevelFilter) Write(p []byte) (n int, err error) {
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

// SetMinLevel is used to update the minimum log level
func (f *LevelFilter) SetMinLevel(min LogLevel) {
	f.MinLevel = min
	f.init()
}

func (f *LevelFilter) init() {
	badLevels := make(map[LogLevel]struct{})
	for _, level := range f.Levels {
		if level == f.MinLevel {
			break
		}
		badLevels[level] = struct{}{}
	}
	f.badLevels = badLevels
}
