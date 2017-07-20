package slog

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	cfg := Config{
		File:   "/tmp/a",
		Debug:  false,
		Prefix: "TST",
	}
	Init(cfg)
	defer os.Remove(cfg.File)

	P("this p-test should be in `%s'", cfg.File)
	D("AIYO d-test should not appear in `%s'", cfg.File)

	txt, err := ioutil.ReadFile(cfg.File)
	if err != nil {
		t.Errorf("failed to read `%s': %v", cfg.File, err)
	}
	s := string(txt)

	if strings.Contains(s, "WARN") == false {
		t.Errorf("logfile `%s' does not WARN prefix", cfg.File)
	}
	if strings.Contains(s, "AIYO") == true {
		t.Errorf("logfile `%s' contains debug when debug not set", cfg.File)
	}
}

func TestStack(t *testing.T) {
	cfg := Config{
		File:   "/tmp/a",
		Debug:  true,
		Prefix: "TST",
	}
	Init(cfg)
	defer os.Remove(cfg.File)

	f := func() { P("say hi to stack") }
	f()

	txt, err := ioutil.ReadFile(cfg.File)
	if err != nil {
		t.Errorf("failed to read `%s': %v", cfg.File, err)
	}
	s := string(txt)

	if strings.Contains(s, "hi") == false {
		t.Errorf("logfile `%s' does not contain canary ", cfg.File)
	}
}

func TestDebug(t *testing.T) {
	cfg := Config{
		File:   "/tmp/a",
		Debug:  true,
		Prefix: "TST",
	}
	Init(cfg)
	defer os.Remove(cfg.File)

	P("this p-test should be in `%s'", cfg.File)
	D("AIYO d-test should appear in `%s'", cfg.File)

	txt, err := ioutil.ReadFile(cfg.File)
	if err != nil {
		t.Errorf("failed to read `%s': %v", cfg.File, err)
	}
	s := string(txt)

	if strings.Contains(s, "WARN") == false {
		t.Errorf("logfile `%s' does not WARN prefix", cfg.File)
	}
	if strings.Contains(s, "AIYO") == false {
		t.Errorf("logfile `%s' does not contain debug when debug set", cfg.File)
	}
}

func TestAuditLogging(t *testing.T) {
	cfg := Config{
		File:      "/tmp/a",
		Debug:     false,
		AuditFile: "/tmp/b",
		Prefix:    "TST",
	}
	Init(cfg)
	defer os.Remove(cfg.File)
	defer os.Remove(cfg.AuditFile)

	A("a-test did some testing for `%s'", cfg.AuditFile)
	txt, err := ioutil.ReadFile(cfg.AuditFile)
	if err != nil {
		t.Errorf("failed to read `%s': %v", cfg.AuditFile, err)
	}

	s := string(txt)
	if strings.Contains(s, "testing") == false {
		t.Errorf("logfile `%s' does not testing line", cfg.AuditFile)
	}
}
