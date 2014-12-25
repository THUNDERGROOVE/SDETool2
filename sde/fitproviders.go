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
