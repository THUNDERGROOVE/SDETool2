package sde

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"net/http"
)

const (
	ProtofitsBaseURL = "http://www.protofits.com/fittings/getCLF/"
)

// GetFitClipboard gets a CLF fit from the clipboard.
func GetFitClipboard() (*SDEFit, error) {
	data, err := clipboard.ReadAll()
	fit := &SDEFit{}
	if err != nil {
		return fit, err
	}
	merr := json.Unmarshal([]byte(data), fit)
	if merr != nil {
		return fit, merr
	}
	fit.fillFields()
	return fit, nil
}

// GetFitProtoFits gets a CLF fit from Protofits.com
// You must provide the id of the fit and it must be shared
func GetFitProtofits(id string) (*SDEFit, error) {
	resp, err := http.Get(fmt.Sprintf("%v%v", ProtofitsBaseURL, id))
	fit := &SDEFit{}
	if err != nil {
		return fit, err
	}
	data, rerr := ioutil.ReadAll(resp.Body)
	if rerr != nil {
		return fit, rerr
	}

	merr := json.Unmarshal(data, fit)
	if merr != nil {
		return fit, merr
	}
	fit.fillFields()
	return fit, nil
}

// GetFitFromFile loads a CLF fit from file.
func GetFitFromFile(filename string) (*SDEFit, error) {
	data, err := ioutil.ReadFile(filename)
	fit := &SDEFit{}
	if err != nil {
		return fit, err
	}
	merr := json.Unmarshal(data, fit)
	if merr != nil {
		return fit, merr
	}
	fit.fillFields()
	return fit, nil
}
