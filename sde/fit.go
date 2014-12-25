package sde

import (
	"fmt"
	"strconv"
)

type CLFMetadata struct {
	Title string `json:"title"`
}

type CLFPreset struct {
	Name    string       `json:"presetname"`
	Modules []*CLFModule `json:"modules"`
}

type CLFSuit struct {
	TypeID  string `json:"typeid"`
	SDEType SDEType
}

type CLFModule struct {
	SDEType  SDEType
	TypeID   string `json:"typeid"`
	SlotType string `json:"slottype"`
	Index    int    `json:"index"`
}

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
