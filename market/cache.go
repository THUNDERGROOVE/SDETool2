package market

import (
	"encoding/json"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"io/ioutil"
	"os"
	"time"
)

const (
	TimeFormat = "01-02 2006 03:04PM"
)

type Cache struct {
	Data   map[string]MarketData `json:"data"`
	Time   string                `json:"time"`
	TypeID int                   `json:"typeid"`
}

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
	} else {
		log.Info("No cache file found.")
	}
	return false
}

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
