package packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
)

const (
	InventoryLayoutNone = iota
	InventoryLayoutSurvival
	InventoryLayoutRecipeBook
	InventoryLayoutCreative
)

const (
	InventoryLeftTabNone = iota
	InventoryLeftTabConstruction
	InventoryLeftTabEquipment
	InventoryLeftTabItems
	InventoryLeftTabNature
	InventoryLeftTabSearch
	InventoryLeftTabSurvival
)

const (
	InventoryRightTabNone = iota
	InventoryRightTabFullScreen
	InventoryRightTabCrafting
	InventoryRightTabArmour
)

// SetPlayerInventoryOptions is a bidirectional packet that can be used to update the inventory options of a player.
type SetPlayerInventoryOptions struct {
	/*
		PhoenixBuilder specific changes.
		Author: Happy2018new

		LeftInventoryTab is the tab that is selected on the left side of the inventory. This is usually for the creative
		inventory. It is one of the InventoryLeftTab constants above.

		This is a mistake of upstream, and they have been corrected this in the future version.
	*/
	LeftInventoryTab int32
	// LeftInventoryTab byte

	/*
		PhoenixBuilder specific changes.
		Author: Happy2018new

		RightInventoryTab is the tab that is selected on the right side of the inventory. This is usually for the player's
		own inventory. It is one of the InventoryRightTab constants above.

		This is a mistake of upstream, and they have been corrected this in the future version.
	*/
	RightInventoryTab int32
	// RightInventoryTab byte

	// Filtering is whether the player has enabled the filtering between recipes they have unlocked or not.
	Filtering bool

	/*
		PhoenixBuilder specific changes.
		Author: Happy2018new

		InventoryLayout is the layout of the inventory. It is one of the InventoryLayout constants above.

		This is a mistake of upstream, and they have been corrected this in the future version.
	*/
	// InventoryLayout is the layout of the inventory. It is one of the InventoryLayout constants above.
	InventoryLayout int32
	// InventoryLayout byte

	/*
		PhoenixBuilder specific changes.
		Author: Liliya233, Happy2018new

		CraftingLayout is the layout of the crafting inventory. It is one of the InventoryLayout constants above.

		This is a mistake of upstream, and they have been corrected this in the future version.
	*/
	CraftingLayout int32
	// CraftingLayout byte
}

// ID ...
func (*SetPlayerInventoryOptions) ID() uint32 {
	return IDSetPlayerInventoryOptions
}

func (pk *SetPlayerInventoryOptions) Marshal(io protocol.IO) {
	// PhoenixBuilder specific changes.
	// Author: Happy2018new
	//
	// This is some mistakes of upstream,
	// and they have been corrected this
	// in the future version.
	{
		io.Varint32(&pk.LeftInventoryTab)
		io.Varint32(&pk.RightInventoryTab)
		io.Bool(&pk.Filtering)
		io.Varint32(&pk.InventoryLayout)
		io.Varint32(&pk.CraftingLayout)
		/*
			io.Uint8(&pk.LeftInventoryTab)
			io.Uint8(&pk.RightInventoryTab)
			io.Bool(&pk.Filtering)
			io.Uint8(&pk.InventoryLayout)
			io.Uint8(&pk.CraftingLayout)
		*/
	}
}
