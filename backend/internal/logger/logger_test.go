package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/harshbarnawa/mintok/backend/internal/config"
)

func TestNewWithWriterUsesJSONInProduction(t *testing.T) {
	var buffer bytes.Buffer
	log := NewWithWriter(config.Config{
		Environment: "production",
		ServiceName: "mintok-api",
		LogLevel:    "info",
	}, &buffer)

	log.Info("startup")

	output := buffer.String()
	if !strings.Contains(output, `"msg":"startup"`) {
		t.Fatalf("expected json log output, got %q", output)
	}
}

func TestNewWithWriterHonorsDebugLevel(t *testing.T) {
	var buffer bytes.Buffer
	log := NewWithWriter(config.Config{
		Environment: "development",
		ServiceName: "mintok-api",
		LogLevel:    "debug",
	}, &buffer)

	log.Debug("debug-enabled")

	output := buffer.String()
	if !strings.Contains(output, "debug-enabled") {
		t.Fatalf("expected debug log output, got %q", output)
	}
}
