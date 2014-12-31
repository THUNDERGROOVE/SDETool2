package args

import (
	"flag"
)

var (
	Version        *string
	VersionCompare *string
	TypeName       *string

	DPS          *bool
	ListVersions *bool
	DownloadAll  *bool
	Compare      *bool
	ToJSON       *bool
	Server       *bool
	Port         *int
	Dump         *bool
	Debug        *bool

	ProtoFits *string
	Clipboard *bool
)

func init() {
	Version = flag.String("v", "", "The SDE version to use.")
	VersionCompare = flag.String("vc", "", "The SDE version used for comparisons.  Currently does nothing")
	TypeName = flag.String("t", "", "A type name to look up")

	DPS = flag.Bool("dps", false, "Print dps of a weapon")
	ListVersions = flag.Bool("versions", false, "Display a list of all versions.")
	DownloadAll = flag.Bool("dl", false, "Download all available versions.")
	Compare = flag.Bool("c", false, "Compare type to other types that share the same main tag.  May not work for all types.")
	ToJSON = flag.Bool("json", false, "Print type to JSON")
	Server = flag.Bool("http", false, "Run a web server to return types as JSON")
	Port = flag.Int("port", 80, "Port used by the http server")
	Dump = flag.Bool("dump", false, "Dump relevant typeids to file")
	Debug = flag.Bool("debug", false, "Print debug information about function timings.")

	ProtoFits = flag.String("pf", "", "Gets a fit from protofits")
	Clipboard = flag.Bool("clip", false, "Get a fit from your clipboard in CLF format")

	flag.Parse()
}
