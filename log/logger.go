package log

import (
	"fmt"
	"github.com/joshlf13/term"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	Log          *log.Logger
	DebugLog     bool
	Color        bool
	TerminalMode bool
)

func callstr() string {
	_, path, line, _ := runtime.Caller(2)
	_, file := filepath.Split(path)
	return fmt.Sprintf("%v:%v", file, line)
}

// TypeString was initially made for our logging functions however it's can be
// used all over the codebase
func TypeString(i []interface{}) string {
	s := ""
	for _, v := range i {
		switch k := v.(type) {
		case int:
			s += fmt.Sprintf("%v ", k)
		case string:
			s += fmt.Sprintf("%v ", k)
		case float64:
			s += fmt.Sprintf("%v ", k)
		default:
			s += fmt.Sprintf("%v ", k)
		}
	}
	return s
}

// LogInit is called to init the logging portion of the util package.  If you
// try using any of the logging functions before calling this you will get a
// nil pointer exception.
func init() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("error opening log file")
	}
	Color = true
	DebugLog = false
	Log = log.New(f, "", log.Ltime)
	TerminalMode = true
	Info("Log started!")
}

// LErr is a function for non-fatal errors.  It will always log to file and
// optionally with -debug print to stdout and also will have nice colored
// output as long as -nocolor is not supplied.
func LogError(i ...interface{}) {
	s := TypeString(i)
	if DebugLog {
		if Color {
			term.Red(os.Stdout, fmt.Sprintf("\r%v Error: %v\n", callstr(), s))
		} else {
			fmt.Printf("\r%v Error: %v\n", callstr(), s)
			if TerminalMode {
				fmt.Printf("\n\r> ")
			}
		}
	}
	Log.SetPrefix("WARN ")
	Log.Println("Error: ", s)
	Log.SetPrefix("")
}

// Info is a function for informing the user what is going on.  It will
// always log to file and optionally with -debug print to stdout and
// also will have nice coloredoutput as long as -nocolor is not supplied.
func Info(i ...interface{}) {
	s := TypeString(i)
	if DebugLog {
		if Color {
			term.Cyan(os.Stdout, fmt.Sprintf("\r%v Info: %v\n", callstr(), s))

		} else {
			fmt.Printf("\r%v Info: %v\n", callstr(), s)
			if TerminalMode {
				fmt.Printf("\n\r> ")
			}
		}
	}
	Log.SetPrefix("INFO ")
	Log.Println(s)
	log.SetPrefix("")
}

// Trace is a function for non-helper functions to call on call.  It will
// always log to file and optionally with -debug print to stdout and
// also will have nice coloredoutput as long as -nocolor is not supplied.
// Don't:
//   Use on primitive, short or otherwise uneeded functions.  An example
//   of ones would be the logging functions
// Do:
//   Use on complicated functions, an example would be most of the sde.go
//   file and any method that uses util.TimeFunction on a defer.
func Trace(i ...interface{}) {
	s := TypeString(i)
	if DebugLog {
		if Color {
			term.Green(os.Stdout, "\rTrace: "+s+"\n")
		} else {
			fmt.Print("\rTrace: " + s)
			if TerminalMode {
				fmt.Printf("\n\r> ")
			}
		}
	}
	Log.SetPrefix("TRACE ")
	Log.Println(s)
	log.SetPrefix("")
}
