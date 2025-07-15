package feral

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/druid"
)

func (cat *FeralDruid) registerShredSpell() {
	flatDamageBonus := 0.07100000232 * cat.ClassSpellScaling // ~77.7265
	hasGlyph := cat.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfShred)

	cat.Shred = cat.RegisterSpell(druid.Cat, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5221},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: druid.DruidSpellShred,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   40,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},

			IgnoreHaste: true,
		},

		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return (!cat.PseudoStats.InFrontOfTarget && !cat.CannotShredTarget) || (hasGlyph && (cat.BerserkCatAura.IsActive() || cat.TigersFuryAura.IsActive()))
		},

		DamageMultiplier: 5,
		CritMultiplier:   cat.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			if cat.AssumeBleedActive || (cat.BleedsActive[target] > 0) {
				baseDamage *= druid.RendAndTearDamageMultiplier
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				cat.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				cat.ApplyBloodletting(target)
			} else {
				spell.IssueRefund(sim)
			}
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := flatDamageBonus + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())

			if cat.AssumeBleedActive || (cat.BleedsActive[target] > 0) {
				baseDamage *= druid.RendAndTearDamageMultiplier
			}

			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMeleeWeaponSpecialHitAndCrit)
		},
	})
}

func (cat *FeralDruid) CurrentShredCost() float64 {
	return cat.Shred.Cost.GetCurrentCost()
}
