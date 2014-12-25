package web

import (
	"encoding/json"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/args"
	"github.com/THUNDERGROOVE/SDETool2/sde"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var SDE sde.SDE

type Error struct {
	Text string `json:"error"`
}

func StartServer() {
	fmt.Printf("Starting http server using SDE version: %v on port: %v\n", *args.Version, *args.Port)
	var err error
	SDE, err = sde.Open(*args.Version)

	if err != nil {
		fmt.Printf("Unable to open SDE file: %v\n", err.Error())
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/type/{typeID:[0-9]+}", HandleGetType)
	r.HandleFunc("/search/{search:(.*)}", HandleSearch)
	r.HandleFunc("/version", HandleVersion)
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, r))
	r.NotFoundHandler = http.HandlerFunc(FourOhFour)
	http.ListenAndServe(fmt.Sprintf(":%v", *args.Port), nil)
}

func procErr(err error, res http.ResponseWriter) bool {
	if err != nil {
		response, _ := json.MarshalIndent(Error{err.Error()}, "", "    ")
		res.Write(response)
		res.WriteHeader(500)
		return true
	}
	return false
}
