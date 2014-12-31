package sde

import (
	"time"
)

const (
	sde17  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.7_674383.zip"
	sde18  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_752135.zip" // Old: http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_739147.zip
	sde18D = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_851720.zip"
	sde19  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.9_853193.zip"
)

var (
	// Versions is a map of all of the available versions.
	Versions map[string]string
)

func init() {
	Versions = make(map[string]string, 0)
	Versions["1.7"] = sde17
	Versions["1.8"] = sde18
	Versions["1.8-delta"] = sde18D
	Versions["1.9"] = sde19
}

// DownloadAllVersions is a function to download every version of the SDE
// that is available.
func DownloadAllVersions() {
	defer Debug(time.Now())

	for k := range Versions {
		download(k)
	}
}
