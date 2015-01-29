package args

import (
	"flag"
)

var (
	// Stuff that effects general behavior
	Version *string

	// Selectors.  Things that chose the type to work with.
	TypeName        *string
	TypeID          *int
	TypeDisplayName *string

	// Selector Modifiers
	ToJSON     *bool
	Attributes *bool
	Stats      *bool
)

func init() {
	// Behavior
	Version = flag.String("v", "1.9", "SDE version to use.")

	// Selectors
	TypeName = flag.String("tn", "", "A TypeName selector")
	TypeID = flag.Int("tid", -1, "A TypeID selector")
	TypeDisplayName = flag.String("td", "", "Searches mDisplayName")

	// Selector modifiers
	ToJSON = flag.Bool("json", false, "Prints type in json")
	Attributes = flag.Bool("a", false, "Prints type attributes")
	Stats = flag.Bool("s", false, "Prints useful stats about a type")

	flag.Parse()
}
