package sde

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/atotto/clipboard"
)

const (
	ProtofitsBaseURL = "http://www.protofits.com/fittings/getCLF/"
)

// GetFitClipboard gets a CLF fit from the clipboard.
func GetFitClipboard() (*Fit, error) {
	defer Debug(time.Now())

	data, err := clipboard.ReadAll()
	fit := &Fit{}
	if err != nil {
		return fit, err
	}
	merr := json.Unmarshal([]byte(data), fit)
	if merr != nil {
		return fit, merr
	}
	fit.FillFields()
	return fit, nil
}

// GetFitProtoFits gets a CLF fit from Protofits.com
// You must provide the id of the fit and it must be shared
func GetFitProtofits(id string) (*Fit, error) {
	defer Debug(time.Now())

	resp, err := http.Get(fmt.Sprintf("%v%v", ProtofitsBaseURL, id))
	fit := &Fit{}
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
	fit.FillFields()
	return fit, nil
}

// GetFitFromFile loads a CLF fit from file.
func GetFitFromFile(filename string) (*Fit, error) {
	defer Debug(time.Now())

	data, err := ioutil.ReadFile(filename)
	fit := &Fit{}
	if err != nil {
		return fit, err
	}
	merr := json.Unmarshal(data, fit)
	if merr != nil {
		return fit, merr
	}
	fit.FillFields()
	return fit, nil
}
