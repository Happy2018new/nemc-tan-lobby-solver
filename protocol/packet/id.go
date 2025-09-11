package packet

const (
	IDTanLoginRequest uint16 = iota
	_
	_
	IDTanEnterRoomRequest
)

const (
	IDTanLoginResponse uint16 = iota
	_
	_
	IDTanEnterRoomResponse
	_
	_
	_
	IDTanNotifyServerReady
)
