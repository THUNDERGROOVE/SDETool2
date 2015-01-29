/*
	SDETool is a command line tool for searching information on items in the
	DUST514 Static Data Export
*/
package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"github.com/THUNDERGROOVE/SDETool2/market"
	"github.com/THUNDERGROOVE/SDETool2/market/graph"
	"github.com/THUNDERGROOVE/SDETool2/web"
	"os"

	"github.com/THUNDERGROOVE/SDETool2/args"
	"github.com/THUNDERGROOVE/SDETool2/sde"
)

var (
	MarketData map[string]map[string]market.MarketData
)

func init() {
	// Move to a data directory
	os.Mkdir(os.Getenv("HOME")+"/.SDETool/", 0777)
	if err := os.Chdir(os.Getenv("HOME") + "/.SDETool/"); err != nil {
		fmt.Println(err.Error())
	}
	MarketData = make(map[string]map[string]market.MarketData, 0)
}

func main() {

	sde.PrintDebug = *args.Debug
	log.DebugLog = *args.Debug

	/* Set default version to 1.9, the latest if version isn't set. */
	if *args.Version == "" {
		*args.Version = "1.9"
	}

	if *args.Server {
		web.StartServer()
		return
	}

	/* Download all SDE versions with -dl */
	if *args.DownloadAll {
		fmt.Printf("Downloading all SDE versions\n")
		sde.DownloadAllVersions()
	}

	/* List all SDE versions */
	if *args.ListVersions {
		fmt.Println("Versions:")
		for k, _ := range sde.Versions {
			fmt.Println(k)
		}
	}

	/* Open chosen SDE version */
	s, err := sde.Open(*args.Version)
	if err != nil {
		fmt.Printf("Unable to open SDE file: %v\n", err.Error())
		return
	}
	fmt.Printf("Opened DB with version %v\n", s.Version)

	sde.GiveSDE(&s)

	//log.Info("Attempting to read cache file")
	//sde.LoadCache(fmt.Sprintf("%v.sde", *args.Version))
	//defer sde.SaveCache(fmt.Sprintf("%v.sde", *args.Version))

	if *args.DoCache {
		s.GobDump()
		return
	}

	if *args.Dump {
		err := s.Dump()
		if err != nil {
			fmt.Println("Error in SDE.Dump()", err.Error())
			return
		}
	}

	if *args.ClassSearch != "" {
		d, class := s.GetTypesByClassName(*args.ClassSearch)
		fmt.Println("Class found: ", class)
		for _, v := range d {
			//v.GetAttributes()
			fmt.Println(v.GetName())
		}
	}

	/* Begin type parsing  */
	if *args.TypeName != "" {
		types, err := s.Search(*args.TypeName)
		log.Info("Found", len(types), "types")
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
		if len(types) == 0 {
			fmt.Println("No such type")
			return
		}
		t := types[0]
		t.GetAttributes()

		if t.TypeID == 367765 || t.TypeName == "container_meta1" {
			t.ESBA()
		}

		HandleType(t)
	}

	if *args.TID != 0 {
		t, _ := s.GetType(*args.TID)

		t.GetAttributes()
		t.Lookup(2)

		HandleType(&t)
	}

	if *args.MultiType != "" {
		types, err := s.Search(*args.MultiType)
		if err != nil {
			log.LogError("Error:", err.Error())
		}
		for _, v := range types {
			v.GetAttributes()
			if v.IsAurum() || v.IsFaction() {
				continue
			}
			HandleType(v)
		}
		if *args.Plot {
			graph.BarSuitData(MarketData)
		}
	}
	if *args.ProtoFits != "" {
		fit, err := sde.GetFitProtofits(*args.ProtoFits)
		if err != nil {
			fmt.Printf("Error getting fit from ProtoFits.com: %v\n", err.Error())
			return
		}
		fmt.Println(fit)
	}
	if *args.Clipboard {
		fit, err := sde.GetFitClipboard()
		if err != nil {
			fmt.Printf("Error getting fit from clipboard %v\n", err.Error())
			return
		}
		fmt.Println(fit)
	}
}

func HandleType(t *sde.SDEType) {
	if t.TypeID == 367765 {
		t.ESBA()
	}
	n := t.GetName()
	if n == "" {
		n = t.TypeName
	}
	fmt.Printf("Name: '%v' | %v\n", n, t.TypeID)
	if *args.DPS {
		fmt.Printf("DPS: %v\n", t.GetDPS())
	}
	if *args.Tags {
		t.PrintTags()
	}
	if *args.Compare {
		tt, err := t.GetSharedTagTypes()
		if err != nil {
			fmt.Printf("Error while getting shared tag types %v\n", err.Error())
			return
		}
		for _, v := range tt {
			v.GetAttributes()
			if !v.IsAurum() && v.IsWeapon() {
				fmt.Printf("DPS: %v %v\n", v.GetDPS(), v.GetName())
			}
		}
	}
	if *args.ToJSON {
		v, err := t.ToJSON()
		if err != nil {
			fmt.Printf("Error marshaling JSON data %v\n", err.Error())
		}
		fmt.Println(v)
	}
	if *args.Market {
		d, _ := market.GetUnitsSold(t)
		fmt.Println("-> Total sold in the last 6 months", d)
		// Do if we have multiple items to work with
		if *args.MultiType != "" {
			MarketData[t.GetName()] = market.GetMarketData(t)
		}
	}

}
