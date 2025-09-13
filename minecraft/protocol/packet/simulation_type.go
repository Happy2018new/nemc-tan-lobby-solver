package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"

const (
	SimulationTypeGame byte = iota
	SimulationTypeEditor
	SimulationTypeTest
	SimulationTypeInvalid
)

// SimulationType is an in-progress packet. We currently do not know the use case.
type SimulationType struct {
	// SimulationType is the simulation type selected.
	SimulationType byte
}

// ID ...
func (*SimulationType) ID() uint32 {
	return IDSimulationType
}

func (pk *SimulationType) Marshal(io protocol.IO) {
	io.Uint8(&pk.SimulationType)
}
