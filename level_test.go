package logutils

import (
	"bytes"
	"log"
	"io"
	"testing"
)

func TestLevelFilter_impl(t *testing.T) {
	var _ io.Writer = new(LevelFilter)
}

func TestLevelFilter(t *testing.T) {
	buf := new(bytes.Buffer)
	filter := &LevelFilter{
		Levels: []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer: buf,
	}

	logger := log.New(filter, "", 0)
	logger.Print("[WARN] foo")
	logger.Println("[ERROR] bar")
	logger.Println("[DEBUG] baz")
	logger.Println("[WARN] buzz")

	result := buf.String()
	expected := "[WARN] foo\n[ERROR] bar\n[WARN] buzz\n"
	if result != expected {
		t.Fatalf("bad: %#v", result)
	}
}
