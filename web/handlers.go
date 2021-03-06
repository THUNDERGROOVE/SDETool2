package web

import (
	"encoding/json"
	"errors"
	"github.com/THUNDERGROOVE/SDETool2/args"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"github.com/THUNDERGROOVE/SDETool2/sde"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// HandleSerach is used to handle /search/[search]
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
		log.LogError("Error encountered while handling a response", err.Error())
		return
	}
	if t, err := SDE.GetType(typeID); err != nil {
		procErr(err, res)
	} else {
		if _, ok := t.Attributes["mDisplayName"]; !ok {
			t.GetAttributes()
		}
		//t.Lookup(2)
		j, _ := t.ToJSON()
		response = []byte(j)
	}
	res.Write(response)
}

// HandleVersion is used to handle /version/
func HandleVersion(res http.ResponseWriter, req *http.Request) {
	v := &struct {
		Version string `json:"version"`
	}{
		Version: *args.Version,
	}
	data, _ := json.Marshal(v)
	res.Write(data)
}

//HandleFit accepts a fit via POST
func HandleFit(res http.ResponseWriter, req *http.Request) {
	d := json.NewDecoder(req.Body)
	var f sde.Fit
	err := d.Decode(&f)
	if procErr(err, res) {
		return
	}
	f.FillFields()
	s := f.Stats()
	res.Write([]byte(s.ToJSON()))
}

// FourOhFour is a 404 handler.
func FourOhFour(res http.ResponseWriter, req *http.Request) {
	err := errors.New("404: Not found")
	m, _ := json.MarshalIndent(Error{err.Error()}, "", "    ")
	res.WriteHeader(404)
	res.Write(m)
}
