package sde

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"reflect"
	"sort"
	"strings"
	"time"
)

var WorthyAttributes map[string]AtterSet

type AtterSet struct {
	SetName          string
	AttributeName    string
	DoRangeFilter    bool
	DoIntervalFilter bool
	ValueFunc        func(t SDEType, val interface{}) interface{}
}

func init() {
	defer Debug(time.Now())
	WorthyAttributes = make(map[string]AtterSet, 0)

	// Biotic stuff
	WorthyAttributes["mCharProp.meleeDamage"] = AtterSet{SetName: "Biotics", AttributeName: "melee damage"}
	WorthyAttributes["mCharProp.maxStamina"] = AtterSet{SetName: "Biotics", AttributeName: "stamina"}
	WorthyAttributes["mCharProp.staminaRecoveryPerSecond"] = AtterSet{SetName: "Biotics", AttributeName: "stamina recovery"}
	WorthyAttributes["mVICProp.groundSpeed"] = AtterSet{SetName: "Biotics", AttributeName: "speed", DoRangeFilter: true}
	WorthyAttributes["mCharProp.movementSprint.groundSpeedScale"] = AtterSet{SetName: "Biotics", AttributeName: "sprint speed", DoRangeFilter: true,
		ValueFunc: func(t SDEType, val interface{}) interface{} {
			if v, ok := t.Attributes["mVICProp.groundSpeed"]; ok {
				if speed, kk := val.(float64); kk {
					if scale, kkk := v.(float64); kkk {
						log.Info("Speed:", speed, "scale:", scale)
						return interface{}(float64(speed * scale))
					}
				}
			} else {
				log.LogError("Type assertion error. speed:", reflect.TypeOf(v), "val:", reflect.TypeOf(val))
			}

			return interface{}(float64(-1))
		}}
	WorthyAttributes["mCharProp.movementRun.strafeSpeedScale"] = AtterSet{SetName: "Biotics", AttributeName: "strafe speed", DoRangeFilter: true,
		ValueFunc: func(t SDEType, val interface{}) interface{} {
			if v, ok := t.Attributes["mVICProp.groundSpeed"]; ok {
				if speed, kk := val.(float64); kk {
					if scale, kkk := v.(float64); kkk {
						log.Info("Speed:", speed, "scale:", scale)
						return interface{}(float64(speed * scale))
					}
				}
			} else {
				log.LogError("Type assertion error.  speed:", reflect.TypeOf(v), "val:", reflect.TypeOf(val))
			}
			return interface{}(float64(-1))
		}}

	// Regen
	WorthyAttributes["mVICProp.healArmorRate"] = AtterSet{SetName: "Regeneration", AttributeName: "armor repair rate"}
	WorthyAttributes["mVICProp.healShieldRate"] = AtterSet{SetName: "Regeneration", AttributeName: "shield recharge rate"}
	WorthyAttributes["mVICProp.shieldRechargeDelay"] = AtterSet{SetName: "Regeneration", AttributeName: "shield recharge delay"}
	WorthyAttributes["mVICProp.shieldRechargePauseOnShieldDepleted"] = AtterSet{SetName: "Regeneration", AttributeName: "shield depleted delay"}

	//HP
	WorthyAttributes["mVICProp.maxArmor"] = AtterSet{SetName: "HP", AttributeName: "armor"}
	WorthyAttributes["mVICProp.maxShield"] = AtterSet{SetName: "HP", AttributeName: "shield"}

	//Fitting
	WorthyAttributes["mVICProp.maxPowerReserve"] = AtterSet{SetName: "Fitting", AttributeName: "PG"}
	WorthyAttributes["mVICProp.maxCpuReserve"] = AtterSet{SetName: "Fitting", AttributeName: "CPU"}
	WorthyAttributes["mVICProp.amountCpuUsage"] = AtterSet{SetName: "Fitting", AttributeName: "CPU usage"}
	WorthyAttributes["mVICProp.amountPowerUsage"] = AtterSet{SetName: "Fitting", AttributeName: "PG usage"}

	//EWAR
	WorthyAttributes["mVICProp.signatureScanPrecision"] = AtterSet{SetName: "EWAR", AttributeName: "scan precision"}
	WorthyAttributes["mVICProp.signatureScanProfile"] = AtterSet{SetName: "EWAR", AttributeName: "scan profile"}
	WorthyAttributes["mVICProp.signatureScanRadius"] = AtterSet{SetName: "EWAR", AttributeName: "scan radius", DoRangeFilter: true}

	// Weapon
	WorthyAttributes["mFireMode0.instantHitDamage"] = AtterSet{SetName: "Weapon", AttributeName: "damage"}
	WorthyAttributes["mFireMode0.fireInterval"] = AtterSet{SetName: "Weapon", AttributeName: "rate of fire", DoIntervalFilter: true}
	WorthyAttributes["mFireMode0.magazineSize"] = AtterSet{SetName: "Weapon", AttributeName: "magazine size"}
	WorthyAttributes["mFireMode0.maxAmmoCount"] = AtterSet{SetName: "Weapon", AttributeName: "total ammo"}
	WorthyAttributes["mFireMode0.reloadTime"] = AtterSet{SetName: "Weapon", AttributeName: "reload time"}

	//Misc
	WorthyAttributes["metaLevel"] = AtterSet{SetName: "Misc", AttributeName: "meta level"}
}

func PrintWorthyStats(t SDEType) {
	defer Debug(time.Now())
	p := make(map[string][]string)
	// Iterate attributes for matches
	for k, v := range WorthyAttributes {
		if val, ok := t.Attributes[k]; ok {
			if _, kk := p[v.SetName]; !kk {
				p[v.SetName] = make([]string, 0)
			}
			if v.ValueFunc != nil {
				val = v.ValueFunc(t, val)
			}
			if v.DoRangeFilter {
				val = DoRangeFilter(val)
			}
			if v.DoIntervalFilter {
				val = IntervalFilter(val)
			}
			p[v.SetName] = append(p[v.SetName], fmt.Sprintf("%v: %v", v.AttributeName, val))
		}
	}
	// Check modifiers.
	for k, v := range t.Attributes {
		if strings.Contains(k, ".attributeName") {
			index := strings.Split(strings.Split(k, "modifier.")[1], ".")[0]
			for kk, vv := range WorthyAttributes {
				if vstr, ok := v.(string); ok {
					log.Info("Attribute", k, "is of index", index, "?")
					if kk == vstr {
						log.Info("Holy tits found a match")
						val := t.Attributes[fmt.Sprintf("modifier.%v.modifierValue", index)]
						mod := t.Attributes[fmt.Sprintf("modifier.%v.modifierType", index)]
						if vv.DoRangeFilter {
							p[vv.SetName] = append(p[vv.SetName], fmt.Sprintf("modifies: %v by %v using %v", vv.AttributeName, DoRangeFilter(val), mod))
						} else {
							p[vv.SetName] = append(p[vv.SetName], fmt.Sprintf("modifies: %v by %v using %v", vv.AttributeName, val, mod))
						}
					}
				} else {
					log.LogError("Attribute name wasn't a stirng? o:")
				}
			}
		}
	}

	keys := make([]string, 0)

	for k, _ := range p {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("=== %v ===\n", key)
		for _, vv := range p[key] {
			fmt.Printf("  %v\n", vv)
		}
	}
}

func DoRangeFilter(i interface{}) float64 {
	if v, ok := i.(float64); ok {
		return float64(v / 100)
	}

	log.Info("Do range filter had no float in interface :/ got", reflect.TypeOf(i))

	return float64(-1)
}

func IntervalFilter(i interface{}) float64 {
	if v, ok := i.(float64); ok {
		return float64(v * 6000)
	}

	log.Info("Interval filter had no float in interface :/ got", reflect.TypeOf(i))

	return float64(-1)
}
