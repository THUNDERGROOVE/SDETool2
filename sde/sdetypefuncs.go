package sde

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// SDEType holds hopefully all of the information you will need about a type.
type SDEType struct {
	ParentSDE  *SDE                   `json:"-"`
	TypeID     int                    `json:"typeid"`
	TypeName   string                 `json:"typename"`
	Attributes map[string]interface{} `json:"attributes"`
}

// GetAttributes grabs the attributes for the type and applied them.  This is
// used to speed up querries for simple lookups.
func (s *SDEType) GetAttributes() error {
	if s.TypeName == "" {
		rows, err := s.ParentSDE.DB.Query(fmt.Sprintf("SELECT typeName FROM CatmaTypes WHERE TypeID == '%v';", s.TypeID))
		if err != nil {
			return err
		}
		if rows.Next() {
			var nTypeName string
			rows.Scan(&nTypeName)
			s.TypeName = nTypeName
		}
	}

	rows, err := s.ParentSDE.DB.Query(fmt.Sprintf("SELECT catmaAttributeName, catmaValueInt, catmaValueReal, catmaValueText FROM CatmaAttributes WHERE TypeID == '%v'", s.TypeID))
	if err != nil {
		return err
	}
	for rows.Next() {
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string
		var catmaValueText string

		rows.Scan(&catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)
		if catmaValueInt != "None" {
			v, _ := strconv.Atoi(catmaValueInt)
			s.Attributes[catmaAttributeName] = v
		}
		if catmaValueReal != "None" {
			v, _ := strconv.ParseFloat(catmaValueReal, 64)
			s.Attributes[catmaAttributeName] = v
		}
		if catmaValueText != "None" {
			s.Attributes[catmaAttributeName] = catmaValueText
		}
	}
	return nil
}

// GetName returns the display name of a type.
func (s *SDEType) GetName() string {
	if name, ok := s.Attributes["mDisplayName"]; ok {
		return name.(string)
	}

	return s.TypeName
}

// GetRoF gets the ROF of any weapon in rounds per minute.  Must vall GetAttributes first
func (s *SDEType) GetRoF() int {
	// if v, ok := s.Attributes["mFireMode0.m_eFireMode"]; ok {
	// 	if v == "DWFM_SingleBurst" {
	bi := s.Attributes["m_BurstInfo.m_fBurstInterval"].(float64)
	fi := s.Attributes["mFireMode0.fireInterval"].(float64)
	return int((bi + fi + 0.01) * 10000)
	// 	}
	// }
	// if i, ok := s.Attributes["mFireMode0.fireInterval"]; ok {
	// 	interval := i.(float64)
	// 	return int(interval * 10000)
	// }
	return 0
}

// GetDPS returns the DPS of a type, if it can.
// Notice: CCP has some all kinds of fucked up shit with bursts and intervals
// don't expect these numbers to be accurate until I can finally fix all of it.
func (s *SDEType) GetDPS() float64 {
	RoF := s.GetRoF()
	var damage float64
	if d, ok := s.Attributes["mFireMode0.instantHitDamage"]; ok {
		damage = d.(float64)
	} else {
		damage = 0
	}

	fmt.Printf("RoF: %v\n", RoF)
	return float64((damage * float64(RoF)) / 60)
}

// IsWeapon returns true if a type has a weapon tag.
func (s *SDEType) IsWeapon() bool {
	for k, v := range s.Attributes {
		if strings.Contains(k, "tag.") {
			if v.(int) == 352335 {
				return true
			}
		}
	}
	return false
}

// IsAurum returns if the item is puchased with aurum.
// Be the soldiar of tomorrow, today with Aurum(C)(TM)(LOLCCP)
func (s *SDEType) IsAurum() bool {
	if strings.Contains(s.TypeName, "aur") {
		return true
	}

	return false
}

// IsObtainable returns True if the item is consumable.
// The name is misleading but it should be used to check if an item is
// obtainable by a player.
func (s *SDEType) IsObtainable() bool {
	if _, ok := s.Attributes["consumable"]; ok {
		return true
	}
	return false
}

// getFromTags is an internal method to return all types that share have a tag
// of the type provided.
func (s *SDEType) getFromTags(t SDEType) ([]*SDEType, error) {
	types := make([]*SDEType, 0)
	rows, err := s.ParentSDE.DB.Query(fmt.Sprintf("SELECT typeID FROM CatmaAttributes WHERE catmaValueInt == '%v';", t.TypeID))
	if err != nil {
		return types, err
	}
	for rows.Next() {
		var nTypeID int
		rows.Scan(&nTypeID)
		types = append(types, &SDEType{
			s.ParentSDE,
			nTypeID,
			"",
			make(map[string]interface{})})
	}
	return types, nil
}

// GetSharedTagTypes returns a slice of SDETypes that share the 'main' tag
// of a type.
func (s *SDEType) GetSharedTagTypes() ([]*SDEType, error) {
	types := make([]*SDEType, 0)
	if s.IsWeapon() {
		for k, v := range s.Attributes {
			if strings.Contains(k, "tag.") {
				tag, _ := s.ParentSDE.GetType(v.(int))
				tag.GetAttributes()
				if strings.Contains(tag.TypeName, "tag_weapon_") { // if s is a scrambler rifle, return all scrambler rifles.
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_core") { // Return all dropsuits since.
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_module_") {
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_amarr") {
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_caldari") {
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_gallente") {
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_minmatar") {
					types, err := s.getFromTags(tag)
					return types, err
				}
			}
		}
	}
	return types, nil
}

// ToJSON returns a Marshaled and indented version of our SDEType.
func (s *SDEType) ToJSON() (string, error) {
	v, err := json.MarshalIndent(s, "", "    ")
	return string(v), err
}
