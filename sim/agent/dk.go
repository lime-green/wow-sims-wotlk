package main

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight/dps"
	"log"
	"strings"
	"time"
)

type DpsDeathknightAgent struct {
	BaseAgent
	simAgent        *dps.DpsDeathknight
	spellOptionsMap map[string]*core.Spell
}

func (dpsDeathknightAgent *DpsDeathknightAgent) Init(session *Session) {
	dk := dpsDeathknightAgent.simAgent

	dpsDeathknightAgent.spellOptionsMap = map[string]*core.Spell{
		"Pestilence":             dk.Pestilence,
		"BloodBoil":              dk.BloodBoil,
		"DeathCoil":              dk.DeathCoil,
		"BloodStrike":            dk.BloodStrike,
		"PlagueStrike":           dk.PlagueStrike,
		"IcyTouch":               dk.IcyTouch,
		"HornOfWinter":           dk.HornOfWinter,
		"DeathStrike":            dk.DeathStrike,
		"Obliterate":             dk.Obliterate,
		"HowlingBlast":           dk.HowlingBlast,
		"FrostStrike":            dk.FrostStrike,
		"EmpowerRuneWeapon":      dk.EmpowerRuneWeapon,
		"RaiseDead":              dk.RaiseDead,
		"DeathAndDecay":          dk.DeathAndDecay,
		"UnbreakableArmor":       dk.UnbreakableArmor,
		"BloodTap":               dk.BloodTap,
		"ArmyOfTheDead":          dk.ArmyOfTheDead,
		"BloodFury":              dk.GetSpell(core.ActionID{SpellID: 33697}),
		"HyperspeedAcceleration": dk.GetSpell(core.ActionID{SpellID: 54758}),
	}

}
func (dpsDeathknightAgent *DpsDeathknightAgent) Cast(spell string, session *Session) Response {
	dk := dpsDeathknightAgent.simAgent
	target := dk.CurrentTarget

	var Spell *core.Spell = nil
	v, ok := dpsDeathknightAgent.spellOptionsMap[spell]
	Spell = v

	if !ok {
		log.Println("unknown spell: ", spell)
		return Response{Success: false}
	}

	castHit := false
	canCast := false

	if Spell != nil {
		canCast = Spell.CanCast(session.sim, target)
	}

	//if !canCast {
	//	log.Println("Unable to cast spell: ", spell)
	//	return Response{Success: false}
	//}

	if canCast {
		//log.Println("Casting spell: ", spell)
		Spell.Cast(session.sim, dk.CurrentTarget)
		castHit = dk.LastOutcome.Matches(core.OutcomeLanded)
	}

	return Response{
		Success: true,
		Body: map[string]interface{}{
			"castPossible": canCast,
			"castHit":      castHit,
		},
	}
}

var defaultRuneTypes = []string{
	"Blood",
	"Blood",
	"Frost",
	"Frost",
	"Unholy",
	"Unholy",
}

var debuffsToTrack = []string{
	"BloodPlague",
	"FrostFever",
}

var buffsToTrack = []string{
	"Blood Tap",
	"Bloodlust",
	"Icy Talons",
	"DMC Greatness Strength Proc",
	"Mjolnir Runestone Proc",
	"Killing Machine Proc",
	"Rime",
	"Potion of Speed",
	"Hyperspeed Acceleration",
	"Blood Fury",
}

func (dpsDeathknightAgent *DpsDeathknightAgent) GetState(session *Session) Response {
	dk := dpsDeathknightAgent.simAgent
	target := dk.CurrentTarget

	// get rune CD remaining for each rune
	var runeCDs = make([]int64, 6)
	for i := 0; i < 6; i++ {
		runeCDs[i] = (dk.RuneReadyAt(session.sim, int8(i)) - session.sim.CurrentTime).Milliseconds()
	}

	// get rune types for each rune
	var runeTypes = make([]string, 6)
	for i := 0; i < 6; i++ {
		if dk.RuneIsDeath(int8(i)) {
			runeTypes[i] = "Death"
		} else {
			runeTypes[i] = defaultRuneTypes[i]
		}
	}

	// get rune grace for each rune
	var runeGraces = make([]int64, 6)
	for i := 0; i < 6; i++ {
		if runeCDs[i] > 0 {
			runeGraces[i] = 0
		} else {
			runeGraces[i] = dk.RuneGraceAt(int8(i), session.sim.CurrentTime).Milliseconds()
		}
	}

	buffs := dk.GetActiveAuras()
	var auraStates []map[string]interface{}
	for _, aura := range buffs {
		name := aura.Tag
		if name == "" {
			name = aura.Label
		}

		auraStates = append(auraStates, map[string]interface{}{
			"name":     name,
			"duration": aura.RemainingDuration(session.sim).Milliseconds(),
			"isActive": true,
		})
	}

	if buffs == nil {
		auraStates = make([]map[string]interface{}, 0)
	}

	// if auraStates does not contain buffsToTrack, add them with duration 0
	for _, buffName := range buffsToTrack {
		found := false
		for _, aura := range auraStates {
			if aura["name"] == buffName {
				found = true
				break
			}
		}
		if !found {
			auraStates = append(auraStates, map[string]interface{}{
				"name":     buffName,
				"duration": 0,
				"isActive": false,
			})
		}
	}

	//iterate over all DK abilities and get their remaining cooldowns
	var abilityCDs []map[string]interface{}

	for spellName, spell := range dpsDeathknightAgent.spellOptionsMap {
		if spell == nil {
			continue
		}

		abilityCD := spell.TimeToReady(session.sim).Milliseconds()
		gcdCost := spell.DefaultCast.GCD

		if !spell.IgnoreHaste {
			gcdCost = spell.Unit.ApplyCastSpeed(spell.DefaultCast.GCD)
		}

		gcdCost = core.MaxDuration(gcdCost, core.GCDMin)

		abilityCDs = append(abilityCDs, map[string]interface{}{
			"name":        spellName,
			"cdRemaining": abilityCD,
			"gcdCost":     gcdCost.Milliseconds(),
			"canCast":     spell.CanCast(session.sim, target),
		})
	}

	// iterate over currentTarget debuffs
	var debuffs []map[string]interface{}
	for _, aura := range target.GetActiveAuras() {
		name := aura.Tag
		if name == "" {
			name = aura.Label
		}

		// remove any trailing hyphen number from the aura name example: "Frost Fever-2" should be "Frost Fever"
		if strings.Contains(name, "-") {
			name = strings.Split(name, "-")[0]
		}

		debuffs = append(debuffs, map[string]interface{}{
			"name":     name,
			"duration": aura.RemainingDuration(session.sim).Milliseconds(),
			"isActive": true,
		})
	}

	if debuffs == nil {
		debuffs = make([]map[string]interface{}, 0)
	}

	// if debuffs doesn't contain debuffs in debuffsToTrack, add them with a duration of 0
	for _, debuffName := range debuffsToTrack {
		found := false
		for _, debuff := range debuffs {
			if debuff["name"] == debuffName {
				found = true
				break
			}
		}

		if !found {
			debuffs = append(debuffs, map[string]interface{}{
				"name":     debuffName,
				"duration": 0,
				"isActive": false,
			})
		}
	}

	gcdRemaining := dk.GCD.TimeToReady(session.sim).Milliseconds()
	var (
		damage        float64 = 0
		meleeDamage   float64 = 0
		diseaseDamage float64 = 0
		abilityDamage float64 = 0
	)

	for _, spell := range dk.Spellbook {
		if spell.OtherID == proto.OtherAction_OtherActionAttack {
			meleeDamage += spell.CalculateDamage()
		} else if spell.SpellID == 55095 || spell.SpellID == 55078 {
			diseaseDamage += spell.CalculateDamage()
		} else {
			// anything other than melee or diseases is considered an ability
			abilityDamage += spell.CalculateDamage()
		}
		damage += spell.CalculateDamage()
	}

	durationSeconds := core.MaxDuration(time.Second, session.sim.CurrentTime).Seconds()
	timeRemaining := session.sim.GetRemainingDuration().Milliseconds()
	dps := damage / durationSeconds

	return Response{
		Success: true,
		Body: map[string]interface{}{
			"gcdAvailable":  gcdRemaining <= 0,
			"isExecute35":   session.sim.IsExecutePhase35(),
			"gcdRemaining":  gcdRemaining,
			"runeCDs":       runeCDs,
			"runeTypes":     runeTypes,
			"runeGraces":    runeGraces,
			"currentTime":   session.sim.CurrentTime.Milliseconds(),
			"buffs":         auraStates,
			"timeRemaining": timeRemaining,
			"abilities":     abilityCDs,
			"debuffs":       debuffs,
			"dps":           dps,
			"isDone":        timeRemaining <= 0,
			"totalDamage":   damage,
			"runicPower":    uint8(dk.CurrentRunicPower()),
			"abilityDamage": abilityDamage,
			"abilityDPS":    abilityDamage / durationSeconds,
			"meleeDamage":   meleeDamage,
			"meleeDPS":      meleeDamage / durationSeconds,
			"diseaseDamage": diseaseDamage,
			"diseaseDPS":    diseaseDamage / durationSeconds,
		},
	}
}
