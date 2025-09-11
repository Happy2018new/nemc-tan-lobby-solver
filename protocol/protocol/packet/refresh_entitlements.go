package packet

import (
	"Eulogist/core/minecraft/standard/protocol"
)

// RefreshEntitlements is sent by the client to the server to refresh the entitlements of the player.
type RefreshEntitlements struct{}

// ID ...
func (*RefreshEntitlements) ID() uint32 {
	return IDRefreshEntitlements
}

func (*RefreshEntitlements) Marshal(protocol.IO) {}
