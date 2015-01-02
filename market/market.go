/*
	market is a package to lookup market information for DUST514
*/
package market

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"github.com/THUNDERGROOVE/SDETool2/sde"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

const BaseURL = "http://public-crest.eveonline.com/market/"

// MarketData is a set to group a slice of MarketDataEntry
type MarketData struct {
	Items []MarketDataEntry `json:"items"`
}

// MarketDataEntry is a struct for Unmarhaling Market data
type MarketDataEntry struct {
	AveragePrice float64 `json:"avgPrice"`
	Date         string  `json:"date"`
	HighPrice    float64 `json:"highPrice"`
	LowPrice     float64 `json:"lowPrice"`
	OrderCount   float64 `json:"orderCount"`
	Volume       float64 `json:"volume"`
	//OrderCountString string `json:"orderCount_str"`
	//VolumeString     string `json:"volume_str"`
}

func GetMarketData(s *sde.SDEType) map[string]MarketData {
	defer sde.Debug(time.Now())
	out := make(map[string]MarketData)
	TypeID := s.TypeID

	if CheckCache(TypeID) {
		c, err := GetCache(TypeID)
		if err != nil {
			log.LogError("Error getting cache", err.Error())
			return out
		}
		return c.Data
	}

	for _, v := range Regions.Regions {
		log.Trace("Getting data for", s.GetName(), "in region", v.Name)
		r, err := http.Get(fmt.Sprintf("%v%v/types/%v/history/", BaseURL, v.TypeID, TypeID))
		if err != nil {
			log.LogError("Error getting market data for Type", s.GetName(), "in region", v.Name, "error:", err.Error())
			continue
		}
		data, rerr := ioutil.ReadAll(r.Body)
		if rerr != nil {
			log.LogError("Error reading http response", rerr.Error())
			continue
		}
		var Data MarketData
		merr := json.Unmarshal(data, &Data)
		if merr != nil {
			log.LogError("Error unmarshaling json", merr.Error())
			continue
		}
		out[v.Name] = Data
	}
	if CacheData(TypeID, out) != nil {
		log.LogError("Error creating cache for type", s.GetName())
	}
	return out
}

func GetUnitsSold(i interface{}) (int, error) {
	defer sde.Debug(time.Now())

	var data map[string]MarketData
	var out int

	switch reflect.TypeOf(i) {
	case reflect.TypeOf(&sde.SDEType{}):
		data = GetMarketData(i.(*sde.SDEType))

	case reflect.TypeOf(make(map[string]MarketData, 0)):
		data = i.(map[string]MarketData)
	default:
		return 0, errors.New("GetUnitsSold was given a type that was not an SDEType or map[string]MarketData")
	}
	for _, v := range data {
		for _, vv := range v.Items {
			out += int(vv.Volume)
		}
	}
	return out, nil
}
