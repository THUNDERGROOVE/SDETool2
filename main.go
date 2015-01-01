/*
	SDETool is a command line tool for searching information on items in the
	DUST514 Static Data Export
*/
package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"github.com/THUNDERGROOVE/SDETool2/market"
	"github.com/THUNDERGROOVE/SDETool2/web"
	"os"

	"github.com/THUNDERGROOVE/SDETool2/args"
	"github.com/THUNDERGROOVE/SDETool2/sde"
)

func init() {
	// Move to a data directory
	os.Mkdir(os.Getenv("HOME")+"/.SDETool/", 0777)
	if err := os.Chdir(os.Getenv("HOME") + "/.SDETool/"); err != nil {
		fmt.Println(err.Error())
	}
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

	if *args.Dump {
		err := s.Dump()
		if err != nil {
			fmt.Println("Error in SDE.Dump()", err.Error())
			return
		}
	}

	/* Begin type parsing  */
	if *args.TypeName != "" {
		types, _ := s.GetTypeWhereNameContains(*args.TypeName)
		if len(types) < 1 {
			fmt.Printf("No types returned\n")
			return
		}
		t := types[0]
		t.GetAttributes()
		HandleType(t)
	}

	if *args.MultiType != "" {
		types, _ := s.GetTypeWhereNameContains(*args.MultiType)
		for _, v := range types {
			v.GetAttributes()
			if v.IsAurum() || v.IsFaction() {
				continue
			}
			HandleType(v)
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
	fmt.Printf("%v\n", t.GetName())
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
		d := market.GetUnitsSold(t)
		fmt.Println("-> Total sold in the last 6 months", d)
	}
}
