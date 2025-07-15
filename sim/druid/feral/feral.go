package feral

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
)

func RegisterFeralDruid() {
	core.RegisterAgentFactory(
		proto.Player_FeralDruid{},
		proto.Spec_SpecFeralDruid,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFeralDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FeralDruid)
			if !ok {
				panic("Invalid spec value for Feral Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFeralDruid(character *core.Character, options *proto.Player) *FeralDruid {
	feralOptions := options.GetFeralDruid()
	selfBuffs := druid.SelfBuffs{}

	cat := &FeralDruid{
		Druid: druid.New(character, druid.Cat, selfBuffs, options.TalentsString),
	}

	// cat.SelfBuffs.InnervateTarget = &proto.UnitReference{}
	// if feralOptions.Options.ClassOptions.InnervateTarget != nil {
	// 	cat.SelfBuffs.InnervateTarget = feralOptions.Options.ClassOptions.InnervateTarget
	// }

	cat.AssumeBleedActive = feralOptions.Options.AssumeBleedActive
	cat.CannotShredTarget = feralOptions.Options.CannotShredTarget

	cat.EnableEnergyBar(core.EnergyBarOptions{
		MaxComboPoints: 5,
		MaxEnergy:      100.0,
		UnitClass:      proto.Class_ClassDruid,
	})
	cat.EnableRageBar(core.RageBarOptions{BaseRageMultiplier: 2.5})

	cat.EnableAutoAttacks(cat, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand:       cat.GetCatWeapon(),
		AutoSwingMelee: true,
	})

	cat.RegisterCatFormAura()
	cat.RegisterBearFormAura()

	return cat
}

type FeralDruid struct {
	*druid.Druid

	// Aura references
	ClearcastingAura        *core.Aura
	PredatorySwiftnessAura  *core.Aura
	SavageRoarBuff          *core.Dot
	SavageRoarDurationTable [6]time.Duration
	TigersFuryAura          *core.Aura

	// Spell references
	SavageRoar *druid.DruidSpell
	Shred      *druid.DruidSpell
	TigersFury *druid.DruidSpell

	tempSnapshotAura   *core.Aura
}

func (cat *FeralDruid) GetDruid() *druid.Druid {
	return cat.Druid
}

func (cat *FeralDruid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.LeaderOfThePack = true
}

func (cat *FeralDruid) Initialize() {
	cat.Druid.Initialize()
	cat.RegisterFeralCatSpells()
	cat.registerSavageRoarSpell()
	cat.registerShredSpell()
	cat.registerTigersFurySpell()
	cat.ApplyPrimalFury()
	cat.ApplyLeaderOfThePack()
	cat.ApplyNurturingInstinct()
	cat.applyOmenOfClarity()
	cat.applyPredatorySwiftness()

	snapshotHandler := func(aura *core.Aura, sim *core.Simulation) {
		previousRipSnapshotPower := cat.Rip.NewSnapshotPower
		previousRakeSnapshotPower := cat.Rake.NewSnapshotPower
		cat.UpdateBleedPower(cat.Rip, sim, cat.CurrentTarget, false, true)
		cat.UpdateBleedPower(cat.Rake, sim, cat.CurrentTarget, false, true)

		if cat.Rip.NewSnapshotPower > previousRipSnapshotPower+0.001 {
			if !cat.tempSnapshotAura.IsActive() || (aura.ExpiresAt() < cat.tempSnapshotAura.ExpiresAt()) {
				cat.tempSnapshotAura = aura

				if sim.Log != nil {
					cat.Log(sim, "New bleed snapshot aura found: %s", aura.ActionID)
				}
			}
		} else if cat.tempSnapshotAura.IsActive() {
			cat.Rip.NewSnapshotPower = previousRipSnapshotPower
			cat.Rake.NewSnapshotPower = previousRakeSnapshotPower
		} else {
			cat.tempSnapshotAura = nil
		}
	}

	// cat.TigersFuryAura.ApplyOnGain(snapshotHandler)
	// cat.TigersFuryAura.ApplyOnExpire(snapshotHandler)
	cat.AddOnTemporaryStatsChange(func(sim *core.Simulation, buffAura *core.Aura, _ stats.Stats) {
		snapshotHandler(buffAura, sim)
	})
}

func (cat *FeralDruid) ApplyTalents() {
	cat.Druid.ApplyTalents()
	cat.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeLeather, 86097)
	cat.applyMastery()
}

func (cat *FeralDruid) applyMastery() {
	const baseMasteryPoints = 8.0
	const masteryModPerPoint = 0.0313 // TODO: We expect 0.03125, possibly bugged?
	const baseMasteryMod = baseMasteryPoints * masteryModPerPoint

	razorClaws := cat.AddDynamicMod(core.SpellModConfig{
		ClassMask:  druid.DruidSpellThrashCat | druid.DruidSpellRake | druid.DruidSpellRip,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: baseMasteryMod + masteryModPerPoint * cat.GetMasteryPoints(),
	})

	cat.AddOnMasteryStatChanged(func(_ *core.Simulation, _ float64, newMasteryRating float64) {
		razorClaws.UpdateFloatValue(baseMasteryMod + masteryModPerPoint * core.MasteryRatingToMasteryPoints(newMasteryRating))
	})

	razorClaws.Activate()
}

func (cat *FeralDruid) Reset(sim *core.Simulation) {
	cat.Druid.Reset(sim)
	cat.Druid.ClearForm(sim)
	cat.CatFormAura.Activate(sim)

	// Reset snapshot power values until first cast
	cat.Rip.CurrentSnapshotPower = 0
	cat.Rip.NewSnapshotPower = 0
	cat.Rake.CurrentSnapshotPower = 0
	cat.Rake.NewSnapshotPower = 0
	cat.tempSnapshotAura = nil
}
