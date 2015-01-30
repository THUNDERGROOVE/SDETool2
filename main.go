/*
	SDETool is a command line tool for searching information on items in the
	DUST514 Static Data Export
*/
package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/args"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"github.com/THUNDERGROOVE/SDETool2/sde"
	"github.com/THUNDERGROOVE/SDETool2/web"
	"os"
	"time"
)

func init() {
	os.Mkdir(os.Getenv("HOME")+"/.SDETool/", 0777)
	if err := os.Chdir(os.Getenv("HOME") + "/.SDETool/"); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	defer sde.Debug(time.Now())
	SDE, err := sde.Open(*args.Version)

	sde.PrimarySDE = &SDE

	if err != nil {
		log.LogError(err.Error())
		return
	}

	var t *sde.SDEType

	if *args.UseCache {
		log.Info("Warning, Cache should only be used for long running setups like the server flags.")
		log.Info("loading cache")
		err := sde.LoadCache(fmt.Sprintf("%v.sde", *args.Version))
		if err != nil {
			log.LogError(err.Error())
		}
		defer sde.SaveCache(fmt.Sprintf("%v.sde", *args.Version))
		SDE.Cache = true
	}

	if *args.Server {
		web.StartServer()
	}

	if *args.TypeID != -1 {
		log.Info("Using TypeID as selector")

		tt, err := SDE.GetType(*args.TypeID)
		if err != nil {
			log.LogError(err.Error())
		}
		t = &tt
	}

	if *args.TypeName != "" {
		log.Info("Using TypeName as selector")

		tt, err := SDE.GetTypeByName(*args.TypeName)
		if err != nil {
			log.LogError(err.Error())
		}
		t = &tt
	}

	if *args.TypeDisplayName != "" {
		log.Info("Using mDisplayName as selector")

		tt, err := SDE.GetTypeWhereNameContains(*args.TypeDisplayName)
		if err != nil {
			log.LogError(err.Error())
		}
		if len(tt) == 0 {
			fmt.Println("No such type.")
			return
		}
		if len(tt) > 1 {
			log.Info("We found more than one type, using the first.")
		}

		t = tt[0]
	}

	if t == nil {
		fmt.Println("No such type or no selectors used.")
		return
	}

	t.GetAttributes()

	fmt.Printf("%v | %v | %v\n", t.GetName(), t.TypeID, t.TypeName)

	if *args.ToJSON {
		fmt.Println(t.ToJSON())
		return
	}

	if *args.Attributes {
		for k, v := range t.Attributes {
			fmt.Printf("'%v' => '%v'\n", k, v)
		}
		return
	}
	if *args.Stats {
		sde.PrintWorthyStats(*t)
	}
}
