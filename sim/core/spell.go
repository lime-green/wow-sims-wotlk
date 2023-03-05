package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/stats"
)

type ApplySpellResults func(sim *Simulation, target *Unit, spell *Spell)
type ExpectedDamageCalculator func(sim *Simulation, target *Unit, spell *Spell, useSnapshot bool) *SpellResult
type CanCastCondition func(sim *Simulation, target *Unit) bool

type SpellConfig struct {
	// See definition of Spell (below) for comments on these.
	ActionID
	SpellSchool  SpellSchool
	ProcMask     ProcMask
	Flags        SpellFlag
	MissileSpeed float64
	BaseCost     float64
	MetricSplits int

	ManaCost   ManaCostOptions
	EnergyCost EnergyCostOptions
	RageCost   RageCostOptions
	RuneCost   RuneCostOptions
	FocusCost  FocusCostOptions

	Cast               CastConfig
	ExtraCastCondition CanCastCondition

	BonusHitRating       float64
	BonusCritRating      float64
	BonusSpellPower      float64
	BonusExpertiseRating float64
	BonusArmorPenRating  float64

	DamageMultiplier         float64
	DamageMultiplierAdditive float64
	CritMultiplier           float64

	ThreatMultiplier float64

	FlatThreatBonus float64

	// Performs the actions of this spell.
	ApplyEffects ApplySpellResults

	// Optional field. Calculates expected average damage.
	ExpectedDamage ExpectedDamageCalculator

	Dot DotConfig
	Hot DotConfig
}

type Spell struct {
	// ID for this spell.
	ActionID

	// The unit who will perform this spell.
	Unit *Unit

	// Fire, Frost, Shadow, etc.
	SpellSchool SpellSchool
	SchoolIndex stats.SchoolIndex

	// Controls which effects can proc from this spell.
	ProcMask ProcMask

	// Flags
	Flags SpellFlag

	// Speed in yards/second. Spell missile speeds can be found in the game data.
	// Example: https://wow.tools/dbc/?dbc=spellmisc&build=3.4.0.44996
	MissileSpeed float64

	ResourceMetrics *ResourceMetrics
	healthMetrics   []*ResourceMetrics

	Cost               SpellCost // Cost for the spell.
	DefaultCast        Cast      // Default cast parameters with all static effects applied.
	CD                 Cooldown
	SharedCD           Cooldown
	ExtraCastCondition CanCastCondition

	// Performs a cast of this spell.
	castFn CastSuccessFunc

	SpellMetrics      []SpellMetrics
	splitSpellMetrics [][]SpellMetrics // Used to split metrics by some condition.

	// Performs the actions of this spell.
	ApplyEffects ApplySpellResults

	// Optional field. Calculates expected average damage.
	expectedDamageInternal ExpectedDamageCalculator

	// The current or most recent cast data.
	CurCast Cast

	BonusHitRating           float64
	BonusCritRating          float64
	BonusSpellPower          float64
	BonusExpertiseRating     float64
	BonusArmorPenRating      float64
	CastTimeMultiplier       float64
	CostMultiplier           float64
	DamageMultiplier         float64
	DamageMultiplierAdditive float64
	CritMultiplier           float64

	// Multiplier for all threat generated by this effect.
	ThreatMultiplier float64

	// Adds a fixed amount of threat to this spell, before multipliers.
	FlatThreatBonus float64

	initialBonusHitRating           float64
	initialBonusCritRating          float64
	initialBonusSpellPower          float64
	initialDamageMultiplier         float64
	initialDamageMultiplierAdditive float64
	initialCritMultiplier           float64
	initialThreatMultiplier         float64
	// Note that bonus expertise and armor pen are static, so we don't bother resetting them.

	resultCache SpellResult

	dots        DotArray
	aoeDot      *Dot
	IgnoreHaste bool
}

func (unit *Unit) OnSpellRegistered(handler SpellRegisteredHandler) {
	for _, spell := range unit.Spellbook {
		handler(spell)
	}
	unit.spellRegistrationHandlers = append(unit.spellRegistrationHandlers, handler)
}

// Registers a new spell to the unit. Returns the newly created spell.
func (unit *Unit) RegisterSpell(config SpellConfig) *Spell {
	if len(unit.Spellbook) > 100 {
		panic(fmt.Sprintf("Over 100 registered spells when registering %s! There is probably a spell being registered every iteration.", config.ActionID))
	}

	// Default the other damage multiplier to 1 if only one or the other is set.
	if config.DamageMultiplier != 0 && config.DamageMultiplierAdditive == 0 {
		config.DamageMultiplierAdditive = 1
	}
	if config.DamageMultiplierAdditive != 0 && config.DamageMultiplier == 0 {
		config.DamageMultiplier = 1
	}

	spell := &Spell{
		ActionID:     config.ActionID,
		Unit:         unit,
		SpellSchool:  config.SpellSchool,
		ProcMask:     config.ProcMask,
		Flags:        config.Flags,
		MissileSpeed: config.MissileSpeed,

		DefaultCast:        config.Cast.DefaultCast,
		CD:                 config.Cast.CD,
		SharedCD:           config.Cast.SharedCD,
		ExtraCastCondition: config.ExtraCastCondition,

		ApplyEffects: config.ApplyEffects,

		expectedDamageInternal: config.ExpectedDamage,

		BonusHitRating:           config.BonusHitRating,
		BonusCritRating:          config.BonusCritRating,
		BonusSpellPower:          config.BonusSpellPower,
		BonusExpertiseRating:     config.BonusExpertiseRating,
		BonusArmorPenRating:      config.BonusArmorPenRating,
		CastTimeMultiplier:       1,
		CostMultiplier:           1,
		DamageMultiplier:         config.DamageMultiplier,
		DamageMultiplierAdditive: config.DamageMultiplierAdditive,
		CritMultiplier:           config.CritMultiplier,

		ThreatMultiplier: config.ThreatMultiplier,
		FlatThreatBonus:  config.FlatThreatBonus,

		splitSpellMetrics: make([][]SpellMetrics, MaxInt(1, config.MetricSplits)),
		IgnoreHaste:       config.Cast.IgnoreHaste,
	}

	if (spell.DamageMultiplier != 0 || spell.ThreatMultiplier != 0) && spell.ProcMask == ProcMaskUnknown {
		panic("Unknown proc mask on " + spell.ActionID.String())
	}

	if (spell.DamageMultiplier != 0 || spell.ThreatMultiplier != 0) && spell.SpellSchool == SpellSchoolNone {
		panic("SpellSchool for spell " + spell.ActionID.String() + " not set")
	}

	switch {
	case spell.SpellSchool.Matches(SpellSchoolPhysical):
		spell.SchoolIndex = stats.SchoolIndexPhysical
	case spell.SpellSchool.Matches(SpellSchoolArcane):
		spell.SchoolIndex = stats.SchoolIndexArcane
	case spell.SpellSchool.Matches(SpellSchoolFire):
		spell.SchoolIndex = stats.SchoolIndexFire
	case spell.SpellSchool.Matches(SpellSchoolFrost):
		spell.SchoolIndex = stats.SchoolIndexFrost
	case spell.SpellSchool.Matches(SpellSchoolHoly):
		spell.SchoolIndex = stats.SchoolIndexHoly
	case spell.SpellSchool.Matches(SpellSchoolNature):
		spell.SchoolIndex = stats.SchoolIndexNature
	case spell.SpellSchool.Matches(SpellSchoolShadow):
		spell.SchoolIndex = stats.SchoolIndexShadow
	}

	if config.ManaCost.BaseCost != 0 || config.ManaCost.FlatCost != 0 {
		spell.Cost = newManaCost(spell, config.ManaCost)
	} else if config.EnergyCost.Cost != 0 {
		spell.Cost = newEnergyCost(spell, config.EnergyCost)
	} else if config.RageCost.Cost != 0 {
		spell.Cost = newRageCost(spell, config.RageCost)
	} else if config.RuneCost.BloodRuneCost != 0 || config.RuneCost.FrostRuneCost != 0 || config.RuneCost.UnholyRuneCost != 0 || config.RuneCost.RunicPowerCost != 0 || config.RuneCost.RunicPowerGain != 0 {
		spell.Cost = newRuneCost(spell, config.RuneCost)
	} else if config.FocusCost.Cost != 0 {
		spell.Cost = newFocusCost(spell, config.FocusCost)
	}

	spell.createDots(config.Dot, false)
	spell.createDots(config.Hot, true)

	spell.castFn = spell.makeCastFunc(config.Cast, spell.applyEffects)

	if spell.ApplyEffects == nil {
		spell.ApplyEffects = func(*Simulation, *Unit, *Spell) {}
	}

	unit.Spellbook = append(unit.Spellbook, spell)

	for _, handler := range unit.spellRegistrationHandlers {
		handler(spell)
	}

	if unit.Env != nil && unit.Env.IsFinalized() {
		spell.finalize()
	}

	return spell
}

// Returns the first registered spell with the given ID, or nil if there are none.
func (unit *Unit) GetSpell(actionID ActionID) *Spell {
	for _, spell := range unit.Spellbook {
		if spell.ActionID.SameAction(actionID) {
			return spell
		}
	}
	return nil
}

// Retrieves an existing spell with the same ID as the config uses, or registers it if there is none.
func (unit *Unit) GetOrRegisterSpell(config SpellConfig) *Spell {
	registered := unit.GetSpell(config.ActionID)
	if registered == nil {
		return unit.RegisterSpell(config)
	} else {
		return registered
	}
}

func (spell *Spell) Dot(target *Unit) *Dot {
	return spell.dots.Get(target)
}
func (spell *Spell) CurDot() *Dot {
	return spell.dots.Get(spell.Unit.CurrentTarget)
}
func (spell *Spell) AOEDot() *Dot {
	return spell.aoeDot
}
func (spell *Spell) Hot(target *Unit) *Dot {
	return spell.dots.Get(target)
}
func (spell *Spell) CurHot() *Dot {
	return spell.dots.Get(spell.Unit.CurrentTarget)
}
func (spell *Spell) AOEHot() *Dot {
	return spell.aoeDot
}
func (spell *Spell) SelfHot() *Dot {
	return spell.aoeDot
}

// Metrics for the current iteration
func (spell *Spell) CurDamagePerCast() float64 {
	if spell.SpellMetrics[0].Casts == 0 {
		return 0
	} else {
		casts := int32(0)
		damage := 0.0
		for _, opponent := range spell.Unit.GetOpponents() {
			casts += spell.SpellMetrics[opponent.UnitIndex].Casts
			damage += spell.SpellMetrics[opponent.UnitIndex].TotalDamage
		}
		return damage / float64(casts)
	}
}

func (spell *Spell) finalize() {
	// Assert that user doesn't set dynamic fields during static initialization.
	if spell.CastTimeMultiplier != 1 {
		panic(spell.ActionID.String() + " has non-default CastTimeMultiplier during finalize!")
	}
	if spell.CostMultiplier != 1 {
		panic(spell.ActionID.String() + " has non-default CostMultiplier during finalize!")
	}
	spell.initialBonusHitRating = spell.BonusHitRating
	spell.initialBonusCritRating = spell.BonusCritRating
	spell.initialBonusSpellPower = spell.BonusSpellPower
	spell.initialDamageMultiplier = spell.DamageMultiplier
	spell.initialDamageMultiplierAdditive = spell.DamageMultiplierAdditive
	spell.initialCritMultiplier = spell.CritMultiplier
	spell.initialThreatMultiplier = spell.ThreatMultiplier

	if len(spell.splitSpellMetrics) > 1 && spell.ActionID.Tag != 0 {
		panic(spell.ActionID.String() + " has split metrics and a non-zero tag, can only have one!")
	}
	for i := range spell.splitSpellMetrics {
		spell.splitSpellMetrics[i] = make([]SpellMetrics, len(spell.Unit.Env.AllUnits))
	}
	spell.SpellMetrics = spell.splitSpellMetrics[0]
}

func (spell *Spell) reset(_ *Simulation) {
	for i := range spell.splitSpellMetrics {
		for j := range spell.SpellMetrics {
			spell.splitSpellMetrics[i][j] = SpellMetrics{}
		}
	}

	// Reset dynamic effects.
	spell.BonusHitRating = spell.initialBonusHitRating
	spell.BonusCritRating = spell.initialBonusCritRating
	spell.BonusSpellPower = spell.initialBonusSpellPower
	spell.CastTimeMultiplier = 1
	spell.CostMultiplier = 1
	spell.DamageMultiplier = spell.initialDamageMultiplier
	spell.DamageMultiplierAdditive = spell.initialDamageMultiplierAdditive
	spell.CritMultiplier = spell.initialCritMultiplier
	spell.ThreatMultiplier = spell.initialThreatMultiplier
}

func (spell *Spell) SetMetricsSplit(splitIdx int32) {
	spell.SpellMetrics = spell.splitSpellMetrics[splitIdx]
	spell.ActionID.Tag = splitIdx
}

func (spell *Spell) CalculateDamage() (damage float64) {
	if spell.Flags.Matches(SpellFlagNoMetrics) {
		return 0
	}

	if len(spell.splitSpellMetrics) == 1 {
		damage += spell.Unit.Metrics.CalculateDamage(spell, spell.ActionID, spell.SpellMetrics)
	} else {
		for i, spellMetrics := range spell.splitSpellMetrics {
			damage += spell.Unit.Metrics.CalculateDamage(spell, spell.ActionID.WithTag(int32(i)), spellMetrics)
		}
	}

	return damage
}

func (spell *Spell) DoneIteration() {
	if spell.Flags.Matches(SpellFlagNoMetrics) {
		return
	}

	if len(spell.splitSpellMetrics) == 1 {
		spell.Unit.Metrics.addSpellMetrics(spell, spell.ActionID, spell.SpellMetrics)
	} else {
		for i, spellMetrics := range spell.splitSpellMetrics {
			spell.Unit.Metrics.addSpellMetrics(spell, spell.ActionID.WithTag(int32(i)), spellMetrics)
		}
	}
}

func (spell *Spell) HealthMetrics(target *Unit) *ResourceMetrics {
	if spell.healthMetrics == nil {
		spell.healthMetrics = make([]*ResourceMetrics, len(spell.Unit.AttackTables))
	}
	if spell.healthMetrics[target.UnitIndex] == nil {
		spell.healthMetrics[target.UnitIndex] = target.NewHealthMetrics(spell.ActionID)
	}
	return spell.healthMetrics[target.UnitIndex]
}

func (spell *Spell) ReadyAt() time.Duration {
	return BothTimersReadyAt(spell.CD.Timer, spell.SharedCD.Timer)
}

func (spell *Spell) IsReady(sim *Simulation) bool {
	if spell == nil {
		return false
	}
	return BothTimersReady(spell.CD.Timer, spell.SharedCD.Timer, sim)
}

func (spell *Spell) TimeToReady(sim *Simulation) time.Duration {
	return MaxTimeToReady(spell.CD.Timer, spell.SharedCD.Timer, sim)
}

// Returns whether a call to Cast() would be successful, without actually
// doing a cast.
func (spell *Spell) CanCast(sim *Simulation, target *Unit) bool {
	if spell == nil {
		return false
	}

	if spell.ExtraCastCondition != nil && !spell.ExtraCastCondition(sim, target) {
		if sim.Log != nil {
			sim.Log("Cant cast because of extra condition")
		}
		return false
	}

	// While casting or channeling, no other action is possible
	if spell.Unit.Hardcast.Expires > sim.CurrentTime {
		if sim.Log != nil {
			sim.Log("Cant cast because already casting/channeling")
		}
		return false
	}

	if spell.DefaultCast.GCD > 0 && !spell.Unit.GCD.IsReady(sim) {
		if sim.Log != nil {
			sim.Log("Cant cast because of GCD")
		}
		return false
	}

	if !BothTimersReady(spell.CD.Timer, spell.SharedCD.Timer, sim) {
		if sim.Log != nil {
			sim.Log("Cant cast because of CDs")
		}
		return false
	}

	if spell.Cost != nil {
		// temp hack
		spell.CurCast.Cost = spell.DefaultCast.Cost
		if !spell.Cost.MeetsRequirement(spell) {
			if sim.Log != nil {
				sim.Log("Cant cast because of resource cost")
			}
			return false
		}
	}

	return true
}

func (spell *Spell) Cast(sim *Simulation, target *Unit) bool {
	if target == nil {
		target = spell.Unit.CurrentTarget
	}
	return spell.castFn(sim, target)
}

// Skips the actual cast and applies spell effects immediately.
func (spell *Spell) SkipCastAndApplyEffects(sim *Simulation, target *Unit) {
	if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
		spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s)",
			spell.ActionID, spell.DefaultCast.Cost, time.Duration(0))
		spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
	}
	spell.applyEffects(sim, target)
}

func (spell *Spell) applyEffects(sim *Simulation, target *Unit) {
	if spell.SpellMetrics == nil {
		spell.reset(sim)
	}
	if target == nil {
		target = spell.Unit.CurrentTarget
	}
	// target can still be null in individual sims when the caster is the enemy target
	if target != nil {
		spell.SpellMetrics[target.UnitIndex].Casts++
	}
	spell.ApplyEffects(sim, target, spell)
}

func (spell *Spell) ApplyAOEThreatIgnoreMultipliers(threatAmount float64) {
	numTargets := spell.Unit.Env.GetNumTargets()
	for i := int32(0); i < numTargets; i++ {
		spell.SpellMetrics[i].TotalThreat += threatAmount
	}
}
func (spell *Spell) ApplyAOEThreat(threatAmount float64) {
	spell.ApplyAOEThreatIgnoreMultipliers(threatAmount * spell.Unit.PseudoStats.ThreatMultiplier)
}

func (spell *Spell) expectedDamageHelper(sim *Simulation, target *Unit, useSnapshot bool) float64 {
	result := spell.expectedDamageInternal(sim, target, spell, useSnapshot)
	if !spell.SpellSchool.Matches(SpellSchoolPhysical) {
		result.Damage /= result.ResistanceMultiplier
		result.Damage *= AverageMagicPartialResistMultiplier
		result.ResistanceMultiplier = AverageMagicPartialResistMultiplier
	}
	result.inUse = false
	return result.Damage
}
func (spell *Spell) ExpectedDamage(sim *Simulation, target *Unit) float64 {
	return spell.expectedDamageHelper(sim, target, false)
}
func (spell *Spell) ExpectedDamageFromCurrentSnapshot(sim *Simulation, target *Unit) float64 {
	return spell.expectedDamageHelper(sim, target, true)
}

// Handles computing the cost of spells and checking whether the Unit
// meets them.
type SpellCost interface {
	// Whether the Unit associated with the spell meets the resource cost
	// requirements to cast the spell.
	MeetsRequirement(*Spell) bool

	// Logs a message for when the cast fails due to lack of resources.
	LogCostFailure(*Simulation, *Spell)

	// Subtracts the resources used from a cast from the Unit.
	SpendCost(*Simulation, *Spell)

	// Space for handling refund mechanics. Not all spells provide refunds.
	IssueRefund(*Simulation, *Spell)
}

func (spell *Spell) IssueRefund(sim *Simulation) {
	spell.Cost.IssueRefund(sim, spell)
}
