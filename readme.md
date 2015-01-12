SDETool2
=======

[![GoDoc](https://godoc.org/github.com/THUNDERGROOVE/SDETool2?status.png)](https://godoc.org/github.com/THUNDERGROOVE/SDETool2)

SDETool2 uses the Dust514 Static Data Export to poll for info.

Building
========
When SDETool2 becomes more useful I will provide precompiled binaries for Windows.

You need [Go](http://golang.org) with your GOPATH environment variable setup with our dependencies

On Windows you will need GCC which you can usually get from either [Cygwin](http://www.cygwin.com/) or [Mingw](http://www.mingw.org/)  I use Mingw so that's your best bet.  It must match the CPU architecture as what you're building for.

In Linux you'll need build-essentials.  Installing it can depend on your distro however you should be able to figure this out if you need.
``` bash
make deps
```
Should download our dependencies for you

Then you should be able to
``` bash
make
```

Dependencies
============

SDETool2 currently uses
```
go-sqlite3 by mattn
mux and handlers from Gorilla
clipboard by atotto
```
All of these can be found on Github

Arguments
=====
```
-t
	Used to select a type
	-dps
		Print the DPS of a weapon
	-c
		Used to print a list of all types that share a distinct tag.
	-json
		Prints the type in JSON format.  Used for debugging mainly.
-v
	Used to set the version.  For a list of version use the -versions flag
-dump
	Used internally to generate the types package
-dl
	Download all versions of the SDE
-versions
	Lists all available SDE versions.
-http
	Starts an http server that acts as a JSON endpoint to the SDE for more information see the server section
	-port
		Set the port.  Defaults to 80
-pf
    If given a string, it will pull a fit from Protofits.com.  Currently this doesn't work as viewing fits requires you to be logged in.
-clip
    If given a string, it will get a CLF fit from your clipboard.
```

You may find more hidden gems in ``` SDETool2 -help ```

Fits
====
I just recently implemented basic CLF fitting imports and added ```SDEType``` fields to the structures.  I plan on adding more methods to lookup information on fits.

My end goal is to be able to create a list of statistics and compare the entire fit across multiple versions of DUST using previous versions of the SDE, including hotfixes in the event they are uploaded.

Server
======

Searching:
```
/search/{search:(.*)}
```
Returns a list of SDETypes.  Until I can either add some caching or somehow improve the speed of getting massive lists of attributes this will take 1-5 seconds to load.

Getting a Type:
```
/type/{TypeID:[0-9]+}
```
It will return some JSON that looks like

``` json
{
    "parentSDE": {
        "version": "1.9"
    },
    "typeId": 364029,
    "typeName": "arm_scout_am_pro_ak0",
    "attributes": {
        "basePrice": 57690,
        "categoryID": 354390,
        "consumable": "True",
        "displayAttributes.0": 353969,
        "displayAttributes.1": 353968,
        "displayAttributes.10": 354138,
        "displayAttributes.11": 354190,
        "displayAttributes.12": 354030,
        "displayAttributes.13": 356885,
        "displayAttributes.14": 356889,
        "displayAttributes.15": 356884,
        "displayAttributes.16": 353835,
        "displayAttributes.17": 354189,
        "displayAttributes.2": 353970,
        "displayAttributes.3": 353971,
        "displayAttributes.4": 353954,
        "displayAttributes.5": 355579,
        "displayAttributes.6": 355580,
        "displayAttributes.7": 353972,
        "displayAttributes.8": 353973,
        "displayAttributes.9": 354137,
        "mBleedOutDuration": 30,
        "mBleedOutHealth": 200,
        "mCharMeleeProp.meleeDamage": 80,
        "mCharMeleeProp.meleeRange": 200,
        "mCharMeleeProp.meleeStaminaCost": 30,
        "mCharProp.maxStamina": 275,
        "mCharProp.movementCrouch.backwardSpeedScale": 1,
        "mCharProp.movementCrouch.groundSpeedScale": 0.3,
        "mCharProp.movementCrouch.jumpStaminaCost": 65,
        "mCharProp.movementCrouch.strafeSpeedScale": 1.1,
        "mCharProp.movementFreeWalk.backwardSpeedScale": 0.25,
        "mCharProp.movementFreeWalk.groundSpeedScale": 0.35,
        "mCharProp.movementFreeWalk.jumpStaminaCost": 65,
        "mCharProp.movementFreeWalk.strafeSpeedScale": 0.7,
        "mCharProp.movementRun.backwardSpeedScale": 1,
        "mCharProp.movementRun.groundSpeedScale": 1,
        "mCharProp.movementRun.jumpStaminaCost": 65,
        "mCharProp.movementRun.strafeSpeedScale": 0.9,
        "mCharProp.movementSprint.backwardSpeedScale": 1,
        "mCharProp.movementSprint.groundSpeedScale": 1.4,
        "mCharProp.movementSprint.jumpStaminaCost": 65,
        "mCharProp.movementSprint.strafeSpeedScale": 0.5,
        "mCharProp.sprintStaminaCostPerSecond": 10,
        "mCharProp.sprintStaminaRegenPenaltyModifier": 0,
        "mCharProp.sprintStaminaRegenPenaltyRecoveryTime": 1,
        "mCharProp.staminaRecoveryPerSecond": 40,
        "mDescription": "The Scout dropsuit is a lightweight suit optimized for enhanced mobility, multi-spectrum stealth, and heightened awareness. Augmented joint servo motors give every movement extra speed and flexibility, while integrated friction and impact dampening materials reduce the overall sound signature. \r\n\r\nBuilding on recent advancements in biotic technology, this suit incorporates an array of cardiovascular augmentations that are automatically administered to the user in battle, improving overall stamina and reducing fatigue. \r\n\r\nWhen missions call for speed and stealth, situations in which heavily armored suits would be more of a burden than an advantage, a scout dropsuit is the best option. The enhanced mobility it provides makes up for its relatively low protection, and when combined with stealth technology modules, the scout suit is the obvious choice for infiltration, counter-espionage, and assassination.",
        "mDisplayName": "Scout ak.0",
        "mHackSpeedFactor": 1.05,
        "mItemTier": "PRO - Scout",
        "mModuleSlots.0.mandatory": "False",
        "mModuleSlots.0.slotType": "IE",
        "mModuleSlots.0.visible": "True",
        "mModuleSlots.1.mandatory": "False",
        "mModuleSlots.1.slotType": "IE",
        "mModuleSlots.1.visible": "True",
        "mModuleSlots.10.mandatory": "False",
        "mModuleSlots.10.slotType": "WS",
        "mModuleSlots.10.visible": "True",
        "mModuleSlots.2.mandatory": "False",
        "mModuleSlots.2.slotType": "GS",
        "mModuleSlots.2.visible": "True",
        "mModuleSlots.3.mandatory": "False",
        "mModuleSlots.3.slotType": "IH",
        "mModuleSlots.3.visible": "True",
        "mModuleSlots.4.mandatory": "False",
        "mModuleSlots.4.slotType": "IH",
        "mModuleSlots.4.visible": "True",
        "mModuleSlots.5.mandatory": "False",
        "mModuleSlots.5.slotType": "WP",
        "mModuleSlots.5.visible": "True",
        "mModuleSlots.6.mandatory": "False",
        "mModuleSlots.6.slotType": "IL",
        "mModuleSlots.6.visible": "True",
        "mModuleSlots.7.mandatory": "False",
        "mModuleSlots.7.slotType": "IL",
        "mModuleSlots.7.visible": "True",
        "mModuleSlots.8.mandatory": "False",
        "mModuleSlots.8.slotType": "IL",
        "mModuleSlots.8.visible": "True",
        "mModuleSlots.9.mandatory": "False",
        "mModuleSlots.9.slotType": "IL",
        "mModuleSlots.9.visible": "True",
        "mShortDescription": "The ak.0 features an expanded low slot loadout as well as improved PG and CPU output",
        "mVICProp.airSpeed": 500,
        "mVICProp.armorDamageScale": 1,
        "mVICProp.groundSpeed": 525,
        "mVICProp.healArmorRate": 1,
        "mVICProp.healHealthRate": 0,
        "mVICProp.healShieldRate": 30,
        "mVICProp.healthDamageScale": 1,
        "mVICProp.maxArmor": 170,
        "mVICProp.maxCpuReserve": 340,
        "mVICProp.maxHealth": 10,
        "mVICProp.maxPowerReserve": 70,
        "mVICProp.maxShield": 60,
        "mVICProp.minDamageToCauseShieldRechargePause": 0,
        "mVICProp.rateOfFireMultiplier": 1,
        "mVICProp.reloadTimeMultiplier": 1,
        "mVICProp.shieldDamageScale": 1,
        "mVICProp.shieldRechargeDelay": 4,
        "mVICProp.shieldRechargePauseOnShieldDepleted": 6,
        "mVICProp.signatureScanPrecision": 40,
        "mVICProp.signatureScanProfile": 35,
        "mVICProp.signatureScanRadius": 2000,
        "metaLevel": 7,
        "requiredSkills.0.skillLevel": 5,
        "requiredSkills.0.skillTypeID": 364594,
        "tag.0": 353506,
        "tag.1": 352339,
        "tag.2": 352332,
        "tag.3": 353502
    }
}

```

I plan to add more later probably for tags and such.
