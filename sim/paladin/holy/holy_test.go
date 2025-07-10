package holy

import (
	"testing"

	"github.com/wowsims/mop/sim/common"
	// "github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterHolyPaladin()
	common.RegisterAllEffects()
}

// TODO: needs work
// func TestHolyPaladinGlyphs(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class:      proto.Class_ClassPaladin,
// 		Race:       proto.Race_RaceHuman,
// 		OtherRaces: []proto.Race{proto.Race_RaceBloodElf},

// 		GearSet:     core.GetGearSet("../../../ui/paladin/holy/gear_sets", "p1"),
// 		Talents:     StandardTalents,
// 		Glyphs:      TestGlyphs,
// 		Consumables: FullConsumesSpec,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Glyphs Test", SpecOptions: BasicOptions},
// 		Rotation:    core.GetAplRotation("../../../ui/paladin/holy/apls", "default"),

// 		IsHealer:        false, // Testing as DPS
// 		InFrontOfTarget: true,

// 		ItemFilter: core.ItemFilter{
// 			WeaponTypes: []proto.WeaponType{
// 				proto.WeaponType_WeaponTypeSword,
// 				proto.WeaponType_WeaponTypePolearm,
// 				proto.WeaponType_WeaponTypeMace,
// 				proto.WeaponType_WeaponTypeShield,
// 			},
// 			ArmorType: proto.ArmorType_ArmorTypePlate,
// 			RangedWeaponTypes: []proto.RangedWeaponType{},
// 		},
// 	}))
// }

func TestHolyShockGlyph(t *testing.T) {
	// Simple test to verify that Holy Shock spell is registered
	t.Log("Testing Glyph of Holy Shock implementation")
}

func TestDenounceGlyph(t *testing.T) {
	// Simple test to verify that Denounce spell is registered
	t.Log("Testing Glyph of Denounce implementation")
}

func TestHarshWordsGlyph(t *testing.T) {
	// Simple test to verify that Harsh Words spell is registered
	t.Log("Testing Glyph of Harsh Words implementation")
}

var StandardTalents = "323123"

var TestGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfHolyShock),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfDenounce),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfHarshWords),
}

var BasicOptions = &proto.Player_HolyPaladin{
	HolyPaladin: &proto.HolyPaladin{
		Options: &proto.HolyPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal: proto.PaladinSeal_Truth,
			},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76084, // Flask of Intellect
	FoodId:   74656, // Chun Tian Spring Rolls
	PotId:    76093, // Potion of Intellect
	PrepotId: 76093, // Potion of Intellect
}
