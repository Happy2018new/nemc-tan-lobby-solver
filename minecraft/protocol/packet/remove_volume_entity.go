package packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
)

// RemoveVolumeEntity indicates a volume entity to be removed from server to client.
type RemoveVolumeEntity struct {
	// PhoenixBuilder specific changes.
	// Author: Happy2018new
	//
	// This is a mistake of upstream.
	EntityRuntimeID uint32
	// EntityRuntimeID uint64

	// Dimension ...
	Dimension int32
}

// ID ...
func (*RemoveVolumeEntity) ID() uint32 {
	return IDRemoveVolumeEntity
}

func (pk *RemoveVolumeEntity) Marshal(io protocol.IO) {
	// PhoenixBuilder specific changes.
	// Author: Happy2018new
	//
	// This is a mistake of upstream.
	{
		io.Varuint32(&pk.EntityRuntimeID)
		// io.Uint64(&pk.EntityRuntimeID)
	}
	io.Varint32(&pk.Dimension)
}
