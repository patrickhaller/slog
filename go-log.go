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
/*  File      = filename for development and production logs, default is production
    Debug     = flag to switch to development logging
    AuditFile = filename to store audit / accounting logs
    Prefix    = short prefix to the filename+function hash, e.g. ORA or NGINX */
type Config struct {
	File      string
	Debug     bool
	AuditFile string
	Prefix    string
}

// F for use before slog is initialized
var F = log.Fatalf

// D for Developers' use
var D = F

// P is for Production use
var P = F

// A is for Audit / Accounting
var A = F

/* Now add drop-in support so people can just
import (
	log "github.com/patrickhaller/slog"
)
*/

// Fatal is for drop-in support
var Fatal = log.Fatal

// Fatalf is for drop-in support
var Fatalf = log.Fatalf

// Fatalln is for drop-in support
var Fatalln = log.Fatalln

// Panic is for drop-in support
var Panic = log.Panic

// Panicf is for drop-in support
var Panicf = log.Panicf

// Panicln is for drop-in support
var Panicln = log.Panicln

// Print is for drop-in support
var Print = log.Print

// Printf is for drop-in support
var Printf = log.Printf

// Println is for drop-in support
var Println = log.Println

func parseLogFile(filename string) *os.File {
	if filename == "STDERR" {
		return os.Stderr
	}
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Printf("Open logfile `%s' failed: %v", filename, err)
		return os.Stderr
	}
	return logFile
}

// Init -- must call this first to set up logging
func Init(cfg Config) {
	log.SetOutput(parseLogFile(cfg.File))
	if cfg.Debug == true {
		log.SetFlags(log.Lmicroseconds | log.LstdFlags | log.Lshortfile)
	}

	D = func(format string, args ...interface{}) {
		if cfg.Debug == false {
			return
		}
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		id := path.Base(f.Name())
		var b bytes.Buffer
		b.WriteString(id)
		b.WriteString(" ")
		b.WriteString(format)
		log.Printf(b.String(), args...)
	}

	P = func(format string, args ...interface{}) {
		// make the hash of filename + function name
		// https://stackoverflow.com/questions/25927660/golang-get-function-name

		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		file, _ := f.FileLine(pc[0])
		h := md5.New()
		io.WriteString(h, path.Base(file))
		io.WriteString(h, path.Base(f.Name()))
		id := fmt.Sprintf("%X", h.Sum(nil))[0:8]
		if cfg.Debug {
			id = fmt.Sprintf("%s:%s", id, path.Base(f.Name()))
		}

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

	auditLog := log.New(parseLogFile(cfg.AuditFile), "", log.LstdFlags)

	A = func(format string, args ...interface{}) {
		auditLog.Printf(format, args...)
	}
}
