package sde

const (
	SDE1_7  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.7_674383.zip"
	SDE1_8  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_752135.zip" // Old: http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_739147.zip
	SDE1_8D = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_851720.zip"
	SDE1_9  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.9_853193.zip"
)

var (
	Versions map[string]string
)

func init() {
	Versions = make(map[string]string, 0)
	Versions["1.7"] = SDE1_7
	Versions["1.8"] = SDE1_8
	Versions["1.8-delta"] = SDE1_8D
	Versions["1.9"] = SDE1_9
}

func DownloadAllVersions() {
	for k, _ := range Versions {
		download(k)
	}
}
