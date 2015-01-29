package args

import (
	"flag"
)

var (
	Version        *string
	VersionCompare *string
	TypeName       *string
	MultiType      *string
	TID            *int
	Plot           *bool

	DPS          *bool
	Tags         *bool
	ListVersions *bool
	DownloadAll  *bool
	Compare      *bool
	ToJSON       *bool
	Server       *bool
	ClassSearch  *string
	Port         *int
	Dump         *bool
	Debug        *bool
	Market       *bool
	DoCache      *bool

	ProtoFits *string
	Clipboard *bool
)

func init() {
	Version = flag.String("v", "", "The SDE version to use.")
	VersionCompare = flag.String("vc", "", "The SDE version used for comparisons.  Currently does nothing")
	TypeName = flag.String("t", "", "A type name to look up")
	MultiType = flag.String("mt", "", "Multiple Types to lookup one after another.")
	TID = flag.Int("id", 0, "A type to lookup by ID")
	Plot = flag.Bool("p", false, "Draw a plot of data. Notice: Only some operations are supported.")

	DPS = flag.Bool("dps", false, "Print dps of a weapon")
	Tags = flag.Bool("tag", false, "Print tags with the type")
	ListVersions = flag.Bool("versions", false, "Display a list of all versions.")
	DownloadAll = flag.Bool("dl", false, "Download all available versions.")
	Compare = flag.Bool("c", false, "Compare type to other types that share the same main tag.  May not work for all types.")
	ToJSON = flag.Bool("json", false, "Print type to JSON")
	Server = flag.Bool("http", false, "Run a web server to return types as JSON")
	ClassSearch = flag.String("class", "", "Print all classes of a typeID")
	Port = flag.Int("port", 80, "Port used by the http server")
	Dump = flag.Bool("dump", false, "Dump relevant typeids to file")
	Debug = flag.Bool("debug", false, "Print debug information about function timings.")
	Market = flag.Bool("m", false, "Get market information")

	ProtoFits = flag.String("pf", "", "Gets a fit from protofits")
	Clipboard = flag.Bool("clip", false, "Get a fit from your clipboard in CLF format")

	DoCache = flag.Bool("cache", false, "Stores all available types into our cache file")

	flag.Parse()
}
