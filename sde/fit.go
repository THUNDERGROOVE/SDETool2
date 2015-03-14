package sde

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

const (
	// SkillLevels
	LevelZero SkillLevel = iota
	LevelOne
	LevelTwo
	LevelThree
	LevelFour
	LevelFive
)

func init() {
	//gob.Register(SkillLevel{})
	gob.Register(CLFMetadata{})
	gob.Register(CLFPreset{})
	gob.Register(CLFSuit{})
	gob.Register(CLFModule{})
	gob.Register(Fit{})
	gob.Register(Stats{})
}

// SkillLevel is a type for each skill level available
type SkillLevel uint8

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
	TypeID  string   `json:"typeid"`
	SDEType *SDEType `json:"-"`
}

// CLFModule holds an individual module in the fit
type CLFModule struct {
	SDEType  SDEType `json:"-"`
	TypeID   string  `json:"typeid"`
	SlotType string  `json:"slottype"`
	Index    int     `json:"index"`
}

// Fit is a structure representing a CLF fit for DUST514 and internal
// structures for calculating stats.
type Fit struct {
	CLFVersion     int         `json:"clf-version"`
	CLFType        string      `json:"X-clf-type"`
	CLFGeneratedBy string      `json:"X-generatedby"`
	Metadata       CLFMetadata `json:"metadata"`
	Suit           CLFSuit     `json:"ship"`
	Fitting        CLFPreset   `json:"presets"`
}

// Stats is a general structure to output all of the stats of a fit.
// Fields values are automatically inserted via ApplySuitBonuses.
// This structure is ready to be exported via JSON.
type Stats struct {
	HealArmorRate       int64   `sde:"mVICProp.healArmorRate"                       json:"repairRate"`
	Shields             int64   `sde:"mVICProp.maxShield"                           json:"shield"`
	Armor               int64   `sde:"mVICProp.maxArmor"                            json:"armor"`
	HealShieldRate      int64   `sde:"mVICProp.healShieldRate"                      json:"shieldRecharge"`
	ShieldDepletedDelay int64   `sde:"mVICProp.shieldRechargePauseOnShieldDepleted" json:"depletedDelay"`
	CPU                 int64   `sde:"mVICProp.maxPowerReserve"                     json:"cpu"`
	CPUUsed             int64   `json:"cpuUsed`
	CPUPercent          int     `json:"cpuPercent"`
	PG                  int64   `sde:"mVICProp.maxPowerReserve"                     json:"pg"`
	PGUsed              int64   `json:"pgUsed"`
	PGPercent           int     `json:"pgPercent"`
	Stamina             float64 `sde:"mCharProp.maxStamina"                         json:"stamina"`
	StaminaRecovery     float64 `sde:"mCharProp.staminaRecoveryPerSecond"           json:"staminaRecovery"`
	ScanPrecision       int64   `sde:"mVICProp.signatureScanPrecision"              json:"scanPrecision"`
	ScanProfile         int64   `sde:"mVICProp.signatureScanProfile"                json:"scanProfile"`
	ScanRadius          int64   `sde:"mVICProp.signatureScanRadius"                 json:"scanRadius"`
	MetaLevel           int64   `sde:"metaLevel"                                    json:"metaLevel"`
}

// FillFields is an internal function used to fill all the extra non-json
// within the SDEFit structure and sub structures.
func (s *Fit) FillFields() {
	log.Info("Filling fields for fit with type", s.Suit.TypeID)
	defer Debug(time.Now())

	if PrimarySDE == nil {
		log.LogError("Error filling SDEFit fields the PrimarySDE is nil.	Set it with GiveSDE()\n")
		return
	}
	typeid, _ := strconv.Atoi(s.Suit.TypeID)
	if typeid <= 0 {
		log.LogError("Fill fields called with no suit")
	}
	t, err := PrimarySDE.GetType(typeid)
	s.Suit.SDEType = &t
	if err != nil {
		log.LogError("Error filling SDEFit fields:", err.Error())
	}

	for _, v := range s.Fitting.Modules {
		tid, _ := strconv.Atoi(v.TypeID)
		if tid <= 0 {
			continue
		} else {
			log.Info(v.TypeID)
		}
		var err1 error
		v.SDEType, err1 = PrimarySDE.GetType(tid)
		if err1 != nil {
			log.LogError("Error filling SDEFit fields:", err.Error())
		}
	}
}

// Stats returns a Stats for a fit.  Assumes all Lvl5 skills.
func (s *Fit) Stats() *Stats {
	ss := &Stats{}
	for _, v := range s.Fitting.Modules {
		log.Info("Fit:", s.Metadata.Title, "module found", v.SDEType.GetName())
	}
	s.ApplySuitBonus(ss, LevelFive)

	t := reflect.TypeOf(Stats{})

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		value := f.Tag.Get("sde")

		if value == "" {
			continue
		}

		if _, ok := s.Suit.SDEType.Attributes[value]; ok {
			attAssert(s.Suit.SDEType, i, ss, value)

		} else {
			log.LogError("Unable to unpack Stats field", value, "from attributes")
		}
	}

	// Static stats
	if ss.PGUsed == 0 {
		ss.PGPercent = 0
	} else {
		ss.PGPercent = int((ss.PGUsed / ss.PG) * 100)
	}
	if ss.CPUUsed == 0 {
		ss.CPUUsed = 0
	} else {
		ss.CPUPercent = int((ss.CPUUsed / ss.PGUsed) * 100)
	}
	return ss
}

// ApplySuitBonus applies bonuses for the suit using the skill level provided.
func (f *Fit) ApplySuitBonus(stats *Stats, skill SkillLevel) {
	f.Suit.SDEType.GetAttributes()
	if v, ok := f.Suit.SDEType.Attributes["requiredSkills.0.skillTypeID"]; ok {
		s, _ := PrimarySDE.GetType(v.(int))
		s.GetAttributes()

		var modcount int
		var bonuscount int

		reg, err := regexp.Compile("(modifier.)(\\d+)(.modifierValue)")
		breg, err2 := regexp.Compile("(bonusModifiers.)(\\d+)(.modifierValue)")
		if err != nil {
			log.LogError("Error compiling regex", err.Error())
		}
		if err2 != nil {
			log.LogError("Error compiling regex", err.Error())
		}

		for k, _ := range s.Attributes {
			if reg.Match([]byte(k)) {
				log.Info("Found modifier", k)
				modcount++
			}
			if breg.Match([]byte(k)) {
				log.Info("Found bonus modifier", k)
				bonuscount++
			}
		}

		log.Info("Attempting to apply modifiers")

		modcount--
		bonuscount--
		log.Info("modcount: %v bonuscount: %v\n", modcount, bonuscount)
		for i := 0; i <= modcount; i++ {
			log.Info("Applying modifier", i)
			var (
				err    error
				ival   interface{}
				itype  interface{}
				istack interface{}
				iname  interface{}
			)
			ival, err = getAttribute(s, fmt.Sprintf("modifier.%v.modifierValue", i))
			itype, err = getAttribute(s, fmt.Sprintf("modifier.%v.modifierType", i))
			istack, err = getAttribute(s, fmt.Sprintf("modifier.%v.stackingPenalized", i))
			iname, err = getAttribute(s, fmt.Sprintf("modifier.%v.attributeName", i))
			if err != nil {
				log.LogError("Unable to get attributes for modifer.  :/", err.Error())
				continue
			}
			val := ival.(float64)
			typ := itype.(string)
			stack := istack.(string)
			name := iname.(string)

			f.applySkill(val, typ, stack, name, skill)

		}

	} else {
		log.LogError("No required skills for type.", v)
		// log.Info("Dumping attributes:")
		// for k, v := range f.Suit.SDEType.Attributes {
		// 	fmt.Printf("\t[%v] => %v\n", k, v)
		// }
	}
}

// ToJSON returns an indented marshaled JSON string of our Stats object.
func (s *Stats) ToJSON() string {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.LogError("Error marshaling stats", err.Error())
	}
	return string(data)
}

// Private helper methods

// applySkill applys a raw skill value
func (f *Fit) applySkill(val float64, typ, stack, name string, level SkillLevel) {
	if f.Suit.SDEType == nil || f.Suit.SDEType.Attributes == nil {
		log.LogError("applySkill called with nil suit SDEType")
		return
	}
	switch reflect.TypeOf(f.Suit.SDEType.Attributes[name]).Kind() {
	case reflect.Int:
		ov := f.Suit.SDEType.Attributes[name].(int)

		for i := 1; i <= int(level); i++ {
			ov = ov * int(val)
		}
		log.Info("Setting attribute", name, "to", ov)
		f.Suit.SDEType.Attributes[name] = ov
	case reflect.Float64:
		ov := f.Suit.SDEType.Attributes[name].(float64)

		for i := 1; i <= int(level); i++ {
			ov = ov * val
			log.Info("Applying level", i, ov)
		}
		log.Info("Setting attribute", name, "to", ov)
		f.Suit.SDEType.Attributes[name] = ov
	default:
		log.LogError("Unsupported type in applySkill")
	}

}

// Private helpers

// getAttribute is a helper to get an attribute from an SDEType.
func getAttribute(s SDEType, a string) (interface{}, error) {
	if v, ok := s.Attributes[a]; ok {
		return v, nil
	} else {
		log.LogError("Wasn't cable to find attribute", a)
	}
	return nil, errors.New("attribute not found" + a)
}

// attAssert is a helper to get a value from SDEType and insert it into
// out Stats using reflection.  It's ugly... No really don't look...
func attAssert(s *SDEType, index int, stats *Stats, value string) {
	t := reflect.TypeOf(Stats{})
	ss := reflect.ValueOf(stats).Elem()
	a := s.Attributes[value]
	log.Info("attAssert with index", index)
	if ss.Field(index).CanSet() == false {
		log.LogError("For whatever reason we cannot set the value at index", index)
		return
	} else {
		switch reflect.TypeOf(a).Kind() {
		case reflect.Float64:
			switch ss.Field(index).Type().Kind() {
			case reflect.Float64: // Field is float.  Can set without conversion
				ss.Field(index).SetFloat(a.(float64))
			case reflect.Int: // Field is int.  Must convert
				ss.Field(index).SetInt(int64(a.(float64)))
			case reflect.Int64:
				ss.Field(index).SetInt(int64(a.(float64)))
			default:
				log.Info("Unsupported type in attAssert", t.Field(index).Type.Kind())
			}
		case reflect.Int:
			switch t.Field(index).Type.Kind() {
			case reflect.Float64: // Field is float.  Convert
				ss.Field(index).SetFloat(float64(a.(int)))
			case reflect.Int: // Field is int.  Must convert int64
				ss.Field(index).SetInt(int64(a.(int)))
			case reflect.Int64:
				ss.Field(index).SetInt(int64(a.(int)))
			default:
				log.Info("Unsupported type in attAssert", t.Field(index).Type.Kind())
			}
		case reflect.Int64:
			switch t.Field(index).Type.Kind() {
			case reflect.Float64: // Field is float.  Convert
				ss.Field(index).SetInt(int64(a.(float64)))
			case reflect.Int: // Field is int.  Must convert int64
				ss.Field(index).SetInt(int64(a.(int)))
			default:
				log.Info("Unsupported type in attAssert", t.Field(index).Type.Kind())
			}
		default:
			log.Info("Unsupported type in main switch attAssert", reflect.TypeOf(a).Kind())
		case reflect.String:
		}
	}
}
