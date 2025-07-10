package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Divine Favor
Instantly restores 100% of your maximum mana and increases your spell power by 25% for 20 sec.
3 min cooldown.
*/
func (paladin *Paladin) registerDivineFavor() {
	actionID := core.ActionID{SpellID: 20216}

	paladin.DivineFavorAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Favor" + paladin.Label,
		ActionID: actionID,
		Duration: time.Second * 20,
	}).AttachMultiplicativePseudoStatBuff(
		&paladin.Unit.PseudoStats.HealingDealtMultiplier, 1.25,
	)

	core.RegisterPercentDamageModifierEffect(paladin.DivineFavorAura, 1.25)

	divineFavor := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: SpellMaskDivineFavor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 3 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Restore 100% mana
			paladin.AddMana(sim, paladin.MaxMana(), spell.ResourceMetrics)

			// Activate the spell power buff
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: paladin.DivineFavorAura,
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: divineFavor,
		Type:  core.CooldownTypeMana,
	})
}
