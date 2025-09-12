package packet

// Pool is a map holding packets indexed by a packet ID.
type Pool map[uint16]Packet

// NewClientPool returns a new pool containing packets sent by a client.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewClientPool() Pool {
	return map[uint16]Packet{
		IDTanLoginRequest:      &TanLoginRequest{},
		IDTanCreateRoomRequest: &TanCreateRoomRequest{},
		IDTanEnterRoomRequest:  &TanEnterRoomRequest{},
		IDTanLeaveRoomRequest:  &TanLeaveRoomRequest{},
		IDTanKickOutRequest:    &TanKickOutRequest{},
	}
}

// NewServerPool returns a new pool containing packets sent by a server.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewServerPool() Pool {
	return map[uint16]Packet{
		IDTanLoginResponse:      &TanLoginResponse{},
		IDTanCreateRoomResponse: &TanCreateRoomResponse{},
		IDTanEnterRoomResponse:  &TanEnterRoomResponse{},
		IDTanLeaveRoomResponse:  &TanLeaveRoomResponse{},
		IDTanKickOutResponse:    &TanKickOutResponse{},
		IDTanNotifyServerReady:  &TanNotifyServerReady{},
	}
}
