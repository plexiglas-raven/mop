package holy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

func (holy *HolyPaladin) registerHolyShock() {
	// Based on MOP Classic data:
	// - Instant cast, 6-second cooldown
	// - 40 yard range
	// - Generates 1 Holy Power
	// - Spell power coefficient: 1.36 (from wowhead)
	// - Can damage enemies or heal allies

	// Create damage version
	holyShockDamage := holy.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 20473},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: paladin.SpellMaskHolyShockDamage,
		Flags:          core.SpellFlagNone,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    holy.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCritPercent: 25.0,
		DamageMultiplier: 1,
		CritMultiplier:   holy.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			actionID := core.ActionID{SpellID: 20473}
			holy.CanTriggerHolyAvengerHpGain(actionID)

			baseDamage := holy.CalcScalingSpellDmg(2.73) // 2.73 spell power coefficient for MOP 5.5
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				holy.HolyPower.Gain(sim, 1, actionID)
			}

			spell.DealDamage(sim, result)
		},
	})

	// Create healing version
	holyShockHeal := holy.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 20473},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellHealing,
		ClassSpellMask: paladin.SpellMaskHolyShockHeal,
		Flags:          core.SpellFlagNone,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    holy.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCritPercent: 25.0,
		DamageMultiplier: 1,
		CritMultiplier:   holy.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			actionID := core.ActionID{SpellID: 20473}
			holy.CanTriggerHolyAvengerHpGain(actionID)

			baseHealing := holy.CalcScalingSpellDmg(2.73) // Same coefficient for healing
			result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			if result.Landed() {
				holy.HolyPower.Gain(sim, 1, actionID)
			}

			spell.DealHealing(sim, result)
		},
	})

	// Create main Holy Shock spell that chooses between damage and healing
	holy.HolyShock = holy.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 20473},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellHealing | core.ProcMaskSpellDamage,
		ClassSpellMask: paladin.SpellMaskHolyShock,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    holy.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if target.IsOpponent(&holy.Unit) {
				// Cast damage version
				holyShockDamage.Cast(sim, target)
			} else {
				// Cast healing version
				holyShockHeal.Cast(sim, target)
			}
		},
	})
}
