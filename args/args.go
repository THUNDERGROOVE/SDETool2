package args

import (
	"flag"
	"github.com/THUNDERGROOVE/SDETool2/log"
)

var (
	// Stuff that effects general behavior
	Version  *string
	UseCache *bool
	Server   *bool
	Port     *int

	// Selectors.  Things that chose the type to work with.
	TypeName        *string
	TypeID          *int
	TypeDisplayName *string

	// Selector Modifiers.  Adjusts what's printed.
	ToJSON     *bool
	Attributes *bool
	Stats      *bool

	//Debug
	Trace  *bool
	Errors *bool
	Info   *bool
)

func init() {
	// Behavior
	Version = flag.String("v", "1.9", "SDE version to use.")
	UseCache = flag.Bool("c", false, "Uses a cache file instead of a database.")
	Server = flag.Bool("http", false, "Starts an SDETool server.")
	Port = flag.Int("port", 80, "Port to be used")

	// Selectors
	TypeName = flag.String("tn", "", "A TypeName selector")
	TypeID = flag.Int("tid", -1, "A TypeID selector")
	TypeDisplayName = flag.String("td", "", "Searches mDisplayName")

	// Selector modifiers
	ToJSON = flag.Bool("json", false, "Prints type in json")
	Attributes = flag.Bool("a", false, "Prints type attributes")
	Stats = flag.Bool("s", false, "Prints useful stats about a type")

	//Debug
	Trace = flag.Bool("trace", false, "Print function timings")
	Errors = flag.Bool("err", false, "Print errors")
	Info = flag.Bool("info", false, "Print info")

	flag.Parse()

	log.TraceLog = *Trace
	log.InfoLog = *Info
	log.ErrLog = *Errors
}
