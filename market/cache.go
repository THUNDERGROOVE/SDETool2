package market

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/THUNDERGROOVE/SDETool2/log"
)

const (
	// TimeFormat is used internally for market cache
	TimeFormat = "01-02 2006 03:04PM"
)

// Cache is our struct for each cache file
type Cache struct {
	Data   map[string]MarketData `json:"data"`
	Time   string                `json:"time"`
	TypeID int                   `json:"typeid"`
}

// CacheData stores our market cache
func CacheData(t int, m map[string]MarketData) error {
	c := Cache{}
	c.Data = m
	c.Time = time.Now().Format(TimeFormat)
	log.Info("Saving time as", time.Now().Format(time.ANSIC))
	c.TypeID = t
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("cache/%v.json", t), data, 0777)
	if err != nil {
		return err
	}
	return nil
}

// CheckCache checks to see if we have valid cache for the type.
func CheckCache(t int) bool {
	checkCache()
	if _, err := os.Stat(fmt.Sprintf("cache/%v.json", t)); err == nil {
		// c, _ := GetCache(t)
		// t, errt := time.Parse(c.Time, TimeFormat)
		// if errt != nil {
		// 	log.LogError("Error", errt.Error(), c.Time)
		// }
		// if time.Now().After(t.Add(time.Hour * 2)) {
		// 	log.Info("-> Cache was old.  Returning false instead")
		// 	return false
		// }
		log.Info("Cache file found.", fmt.Sprintf("cache/%v.json", t), "Using it")
		return true
	}

	log.Info("No cache file found.")

	return false
}

// GetCache returns a cache for a typeid
func GetCache(t int) (Cache, error) {
	out := Cache{}
	data, err := ioutil.ReadFile(fmt.Sprintf("cache/%v.json", t))
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

func checkCache() {
	os.Mkdir("cache", 0777)
}
