package packet

import (
	"Eulogist/core/minecraft/standard/protocol"
)

// RequestPermissions is a packet sent from the client to the server to request permissions that the client does not
// currently have. It can only be sent by operators and host in vanilla Minecraft.
type RequestPermissions struct {
	// EntityUniqueID is the unique ID of the player. The unique ID is unique for the entire world and is
	// often used in packets. Most servers send an EntityUniqueID equal to the EntityRuntimeID.
	EntityUniqueID int64

	/*
		PhoenixBuilder specific changes.
		Author: Happy2018new

		PermissionLevel is the current permission level of the player. This is one of the constants that may be found
		in the AdventureSettings packet.

		This is a mistake of upstream.
	*/
	PermissionLevel int32
	// PermissionLevel uint8

	// RequestedPermissions contains the requested permission flags.
	RequestedPermissions uint16
}

// ID ...
func (*RequestPermissions) ID() uint32 {
	return IDRequestPermissions
}

func (pk *RequestPermissions) Marshal(io protocol.IO) {
	io.Int64(&pk.EntityUniqueID)

	// PhoenixBuilder specific changes.
	// Author: Happy2018new
	//
	// This is a mistake of upstream.
	{
		io.Varint32(&pk.PermissionLevel)
		// io.Uint8(&pk.PermissionLevel)
	}

	io.Uint16(&pk.RequestedPermissions)
}
