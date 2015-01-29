package sde

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"strings"
)

var WorthyAttributes map[string]AtterSet

type AtterSet struct {
	SetName       string
	AttributeName string
	DoRangeFilter bool
}

func init() {
	WorthyAttributes = make(map[string]AtterSet, 0)

	// Biotic stuff
	WorthyAttributes["mCharProp.meleeDamage"] = AtterSet{SetName: "Biotics", AttributeName: "melee damage"}
	WorthyAttributes["mCharProp.maxStamina"] = AtterSet{SetName: "Biotics", AttributeName: "stamina"}
	WorthyAttributes["mCharProp.staminaRecoveryPerSecond"] = AtterSet{SetName: "Biotics", AttributeName: "stamina recovery"}
	WorthyAttributes["mVICProp.groundSpeed"] = AtterSet{SetName: "Biotics", AttributeName: "speed"}

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

	//Misc
	WorthyAttributes["metaLevel"] = AtterSet{SetName: "Misc", AttributeName: "meta level"}
}

func PrintWorthyStats(t SDEType) {
	p := make(map[string][]string)
	// Iterate attributes for matches
	for k, v := range WorthyAttributes {
		if val, ok := t.Attributes[k]; ok {
			if _, kk := p[v.SetName]; !kk {
				p[v.SetName] = make([]string, 0)
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
						p[vv.SetName] = append(p[vv.SetName], fmt.Sprintf("modifies: %v by %v using %v", vv.AttributeName, val, mod))
					}
				} else {
					log.LogError("Attribute name wasn't a stirng? o:")
				}
			}
		}
	}

	for k, v := range p {
		fmt.Printf("=== %v ===\n", k)
		for _, vv := range v {
			fmt.Println("  ", vv)
		}
	}
}
