package sde

import (
	"fmt"
	"strconv"
	"time"
)

// CLFMetadata holds the metadata portion of a CLF fit
type CLFMetadata struct {
	Title string `json:"title"`
}

// CLFPreset is an individual preset in a CLF fit
type CLFPreset struct {
	Name    string       `json:"presetname"`
	Modules []*CLFModule `json:"modules"`
}

// CLFSuit houses the "ship" portion of a CLF fit
type CLFSuit struct {
	TypeID  string  `json:"typeid"`
	SDEType SDEType `json:"-"`
}

// CLF module holds an individual module in the fit
type CLFModule struct {
	SDEType  SDEType `json:"-"`
	TypeID   string  `json:"typeid"`
	SlotType string  `json:"slottype"`
	Index    int     `json:"index"`
}

// SDEFit is a structure representing a CLF fit for DUST514 and internal
// structures for calculating stats.
type SDEFit struct {
	CLFVersion     int         `json:"clf-version"`
	CLFType        string      `json:"X-clf-type"`
	CLFGeneratedBy string      `json:"X-generatedby"`
	Metadata       CLFMetadata `json:"metadata"`
	Suit           CLFSuit     `json:"ship"`
	Fitting        CLFPreset   `json:"presets"`
}

// fillFields is an internal function used to fill all the extra non-json
// within the SDEFit structure and sub structures.
func (s *SDEFit) fillFields() {
	defer Debug(time.Now())

	if PrimarySDE == nil {
		fmt.Printf("Error filling SDEFit fields the PrimarySDE is nil.  Set it with GiveSDE()\n")
		return
	}
	typeid, _ := strconv.Atoi(s.Suit.TypeID)
	var err error
	s.Suit.SDEType, err = PrimarySDE.GetType(typeid)
	if err != nil {
		fmt.Printf("Error filling SDEFit fields: %v\n", err.Error())
	}

	for _, v := range s.Fitting.Modules {
		tid, _ := strconv.Atoi(v.TypeID)
		var err1 error
		v.SDEType, err1 = PrimarySDE.GetType(tid)
		if err1 != nil {
			fmt.Printf("Error filling SDEFit fields: %v\n", err.Error())
		}
	}
}
