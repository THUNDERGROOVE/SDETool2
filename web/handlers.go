package web

import (
	"encoding/json"
	"errors"
	"github.com/THUNDERGROOVE/SDETool2/args"
	"github.com/THUNDERGROOVE/SDETool2/sde"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func HandleSearch(res http.ResponseWriter, req *http.Request) {
	defer sde.Debug(time.Now())
	response := make([]byte, 0)
	v := mux.Vars(req)
	vs := v["search"]

	if types, err := SDE.GetTypeWhereNameContains(vs); err != nil {
		procErr(err, res)
	} else {
		for _, v := range types {
			v.GetAttributes()
		}
		var err1 error
		response, err1 = json.MarshalIndent(types, "", "    ")
		if err1 != nil {
			procErr(err1, res)
		}
	}
	res.Write(response)
}

// HandleGetType is used to handle /type/[TypeID]
func HandleGetType(res http.ResponseWriter, req *http.Request) {
	defer sde.Debug(time.Now())
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
		t.GetAttributes()
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
