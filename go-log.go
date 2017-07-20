package slog

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

// Config holds our configuration details
//    File = filename for development and production logs, default is production
//    Debug = flag to switch to development logging
//    AuditFile = filename to store audit / accounting logs
//    Prefix = short prefix to the filename+function hash, e.g. ORA or NGINX
type Config struct {
	File      string
	Debug     bool
	AuditFile string
	Prefix    string
}

// D for Developers' use
var D func(format string, args ...interface{})

// P is for Production use
var P func(format string, args ...interface{})

// A is for Audit / Accounting
var A func(format string, args ...interface{})

// Init -- must call this first to set up logging
func Init(cfg Config) {
	logFile, err := os.OpenFile(cfg.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Printf("Open logfile `%s' failed: %v", cfg.File, err)
	}
	log.SetOutput(logFile)
	if cfg.Debug == true {
		log.SetFlags(log.Lmicroseconds | log.LstdFlags | log.Lshortfile)
	}

	D = func(format string, args ...interface{}) {
		if cfg.Debug == true {
			log.Printf(format, args...)
		}
	}

	P = func(format string, args ...interface{}) {
		// make the hash of filename + function name
		pc := make([]uintptr, 10) // at least 1 entry needed
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		file, _ := f.FileLine(pc[0])
		h := md5.New()
		io.WriteString(h, path.Base(file))
		io.WriteString(h, path.Base(f.Name()))
		id := fmt.Sprintf("%X", h.Sum(nil))[0:8]

		var b bytes.Buffer
		b.WriteString("WARN ")
		if cfg.Prefix != "" {
			b.WriteString(cfg.Prefix)
			b.WriteString("-")
		}
		b.WriteString(id)
		b.WriteString(" ")
		b.WriteString(format)
		log.Printf(b.String(), args...)
	}

	if cfg.AuditFile == "" {
		return
	}

	alogFile, err := os.OpenFile(cfg.AuditFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Printf("Open logfile `%s' failed: %v", cfg.AuditFile, err)
	}
	auditLog := log.New(alogFile, "", log.LstdFlags)

	A = func(format string, args ...interface{}) {
		auditLog.Printf(format, args...)
	}
}
