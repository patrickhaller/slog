package slog

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	cfg := Config{
		File:      "/tmp/a",
		Debug:     false,
		AuditFile: "/tmp/b",
		Prefix:    "TST",
	}
	Init(cfg)
	P("this p-test should be in `%s'", cfg.File)
	D("XXX d-test should not appear in `%s'", cfg.File)

	if txt, err := ioutil.ReadFile(cfg.File); err == nil {
		s := string(txt)

		if strings.Contains(s, "WARN") == false {
			t.Errorf("logfile `%s' does not WARN prefix", cfg.File)
		}
		if strings.Contains(s, "XXX") == true {
			t.Errorf("logfile `%s' contains debug when debug not set", cfg.File)
		}
		println(s)
	} else {
		t.Errorf("failed to read `%s': %v", cfg.File, err)
	}
	os.Remove(cfg.File)

	if cfg.AuditFile == "" {
		return
	}

	A("a-test did some testing for `%s'", cfg.AuditFile)
	if txt, err := ioutil.ReadFile(cfg.AuditFile); err == nil {
		s := string(txt)

		if strings.Contains(s, "testing") == false {
			t.Errorf("logfile `%s' does not testing line", cfg.AuditFile)
		}
	} else {
		t.Errorf("failed to read `%s': %v", cfg.AuditFile, err)
	}
	os.Remove(cfg.AuditFile)
}
