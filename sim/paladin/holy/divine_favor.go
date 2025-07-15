package holy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

/*
Divine Favor
Increases spell casting haste by 20% and spell critical chance by 20% for 20 sec.
3 min cooldown.
*/
func (holy *HolyPaladin) registerDivineFavor() {
	actionID := core.ActionID{SpellID: 31842}

	holy.DivineFavorAura = holy.RegisterAura(core.Aura{
		Label:    "Divine Favor" + holy.Label,
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// +20% spell haste
			holy.MultiplyCastSpeed(sim, 1.2)
			
			// +20% spell critical chance
			holy.AddStatDynamic(sim, stats.SpellCritPercent, 20*core.CritRatingPerCritPercent)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// Remove +20% spell haste
			holy.MultiplyCastSpeed(sim, 1/1.2)
			
			// Remove +20% spell critical chance
			holy.AddStatDynamic(sim, stats.SpellCritPercent, -20*core.CritRatingPerCritPercent)
		},
	})

	divineFavor := holy.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: paladin.SpellMaskDivineFavor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    holy.NewTimer(),
				Duration: 3 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Activate the buff
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: holy.DivineFavorAura,
	})

	holy.AddMajorCooldown(core.MajorCooldown{
		Spell: divineFavor,
		Type:  core.CooldownTypeDPS,
	})
}
