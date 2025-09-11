package packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
)

// RefreshEntitlements is sent by the client to the server to refresh the entitlements of the player.
type RefreshEntitlements struct{}

// ID ...
func (*RefreshEntitlements) ID() uint32 {
	return IDRefreshEntitlements
}

func (*RefreshEntitlements) Marshal(protocol.IO) {}
