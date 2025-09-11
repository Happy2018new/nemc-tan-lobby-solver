package packet

import "Eulogist/core/minecraft/standard/protocol"

// UpdateSubChunkBlocks is essentially just UpdateBlock packet, however for a set of blocks in a sub-chunk.
type UpdateSubChunkBlocks struct {
	/*
		PhoenixBuilder specific changes.
		Author: Happy2018new

		Position is the block position of the sub-chunk being referred to.
		If Position is (x, y, z), then (x>>4, y>>4, z>>4) is the corresponding sub chunk position.

		This is a mistake of upstream.
	*/
	Position protocol.BlockPos
	// Position is the position of the sub-chunk being referred to.
	// Position protocol.SubChunkPos

	// Blocks contains each updated block change entry.
	Blocks []protocol.BlockChangeEntry
	// Extra contains each updated block change entry for the second layer, usually for waterlogged blocks.
	Extra []protocol.BlockChangeEntry
}

// ID ...
func (*UpdateSubChunkBlocks) ID() uint32 {
	return IDUpdateSubChunkBlocks
}

func (pk *UpdateSubChunkBlocks) Marshal(io protocol.IO) {
	// PhoenixBuilder specific changes.
	// Author: Happy2018new
	//
	// This is a mistake of upstream.
	{
		io.UBlockPos(&pk.Position)
		// io.SubChunkPos(&pk.Position)
	}
	protocol.Slice(io, &pk.Blocks)
	protocol.Slice(io, &pk.Extra)
}
