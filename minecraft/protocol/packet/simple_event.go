package packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
)

const (
	SimpleEventCommandsEnabled = iota + 1
	SimpleEventCommandsDisabled
	SimpleEventUnlockWorldTemplateSettings
)

// SimpleEvent is sent by the server to send a 'simple event' to the client, meaning an event without any
// additional event data. The event is typically used by the client for telemetry.
type SimpleEvent struct {
	/*
		PhoenixBuilder specific changes.
		Author: Happy2018new

		EventType is the type of the event to be called. It is one of the constants that may be found above.

		This is a mistake of upstream.
	*/
	EventType uint16
	// EventType int16
}

// ID ...
func (*SimpleEvent) ID() uint32 {
	return IDSimpleEvent
}

func (pk *SimpleEvent) Marshal(io protocol.IO) {
	// PhoenixBuilder specific changes.
	// Author: Happy2018new
	//
	// This is a mistake of upstream.
	{
		io.Uint16(&pk.EventType)
		// io.Int16(&pk.EventType)
	}
}
