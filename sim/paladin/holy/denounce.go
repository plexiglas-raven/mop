package holy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/paladin"
)

func (holy *HolyPaladin) registerDenounce() {
	// Base cast time is 1.5 seconds
	// Glyph of Denounce: Holy Shock reduces cast time by 0.5 seconds per stack (max 3 stacks for instant cast)

	holy.Denounce = holy.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2812},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: paladin.SpellMaskDenounce,
		Flags:          core.SpellFlagAPL,

		MaxRange: 30,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 1500,
				GCD:      core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				// Apply glyph cast time reduction if active
				if holy.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDenounce) {
					if aura := holy.GetAura("Glyph of Denounce"); aura != nil && aura.IsActive() {
						reduction := time.Millisecond * 500 * time.Duration(aura.GetStacks())
						cast.CastTime = cast.CastTime - reduction
						if cast.CastTime < 0 {
							cast.CastTime = 0
						}
					}
				}
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   holy.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.807, // Spell power coefficient for Denounce

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := holy.CalcScalingSpellDmg(0.807)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			// Consume the glyph stacks on cast
			if holy.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDenounce) {
				if aura := holy.GetAura("Glyph of Denounce"); aura != nil && aura.IsActive() {
					aura.Deactivate(sim)
				}
			}

			spell.DealDamage(sim, result)
		},
	})
}
