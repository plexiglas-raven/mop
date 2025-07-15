package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const RendAndTearBonusCritPercent = 35.0
const RendAndTearDamageMultiplier = 1.2

// Modifies the Bleed aura to apply the bonus.
func (druid *Druid) applyRendAndTear(aura core.Aura) core.Aura {
	if druid.AssumeBleedActive {
		return aura
	}

	aura.ApplyOnGain(func(aura *core.Aura, _ *core.Simulation) {
		druid.BleedsActive[aura.Unit]++
	})
	aura.ApplyOnExpire(func(aura *core.Aura, _ *core.Simulation) {
		druid.BleedsActive[aura.Unit]--
	})

	return aura
}

func (druid *Druid) ApplyPrimalFury() {
	actionID := core.ActionID{SpellID: 16961}
	rageMetrics := druid.NewRageMetrics(actionID)
	const autoRageGen = 15.0
	mangleRageGen := autoRageGen * core.TernaryFloat64((druid.Spec == proto.Spec_SpecGuardianDruid) && druid.Talents.SoulOfTheForest, 1.3, 1)
	cpMetrics := druid.NewComboPointMetrics(actionID)

	druid.RegisterAura(core.Aura{
		Label:    "Primal Fury",
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}

			if druid.InForm(Bear) {
				if spell == druid.MHAutoSpell {
					druid.AddRage(sim, autoRageGen, rageMetrics)
				} else if druid.MangleBear.IsEqual(spell) {
					druid.AddRage(sim, mangleRageGen, rageMetrics)
				}
			} else if druid.InForm(Cat) {
				if spell.Matches(DruidSpellBuilder) {
					druid.AddComboPoints(sim, 1, cpMetrics)
				}
			}
		},
	})
}

func (druid *Druid) ApplyLeaderOfThePack() {
	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 68285})
	manaRestore := 0.08
	healthRestore := 0.04

	icd := core.Cooldown{
		Timer:    druid.NewTimer(),
		Duration: time.Second * 6,
	}

	healingSpell := druid.RegisterSpell(Cat|Bear, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 34299},
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagIgnoreModifiers,
		DamageMultiplier: 1,
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, target, healthRestore*spell.Unit.MaxHealth(), spell.OutcomeHealing)
		},
	})

	druid.RegisterAura(core.Aura{
		Icd:      &icd,
		Label:    "Improved Leader of the Pack",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) || !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			if !druid.InForm(Cat | Bear) {
				return
			}
			icd.Use(sim)
			druid.AddMana(sim, druid.MaxMana()*manaRestore, manaMetrics)
			healingSpell.Cast(sim, &druid.Unit)
		},
	})
}

func (druid *Druid) ApplyNurturingInstinct() {
	druid.GetSpellPowerValue = func(spell *core.Spell) float64 {
		sp := druid.GetStat(stats.SpellPower) + spell.BonusSpellPower

		if spell.ProcMask.Matches(core.ProcMaskSpellHealing) || (spell.SpellSchool == core.SpellSchoolNature) {
			sp += druid.GetStat(stats.Agility)
		}

		return sp
	}
}
