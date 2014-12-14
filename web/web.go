package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/args"
	"github.com/THUNDERGROOVE/SDETool2/sde"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
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
	r.HandleFunc("/version", HandleVersion)
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, r))
	r.NotFoundHandler = http.HandlerFunc(FourOhFour)
	http.ListenAndServe(fmt.Sprintf(":%v", *args.Port), nil)
}

func HandleGetType(res http.ResponseWriter, req *http.Request) {
	response := make([]byte, 0)
	v := mux.Vars(req)
	vs := v["typeID"]
	typeID, err := strconv.Atoi(vs)
	if procErr(err, res) {
		return
	}
	if t, err := SDE.GetType(typeID); err != nil {
		procErr(err, res)
	} else {
		t.GetAtrributes()
		j, _ := t.ToJSON()
		response = []byte(j)
	}
	res.Write(response)
}
func HandleVersion(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(*args.Version))
}

func FourOhFour(res http.ResponseWriter, req *http.Request) {
	err := errors.New("404: Not found")
	m, _ := json.MarshalIndent(Error{err.Error()}, "", "    ")
	res.WriteHeader(404)
	res.Write(m)
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
