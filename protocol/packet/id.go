package packet

// Client -> Server
const (
	IDTanLoginRequest uint16 = iota
	IDTanCreateRoomRequest
	_
	IDTanEnterRoomRequest
	_
	IDTanLeaveRoomRequest
	IDTanKickOutRequest
	_
	_
	IDTanChangeRoomPrivacyRequest // TODO
	IDTanExtendWhiteListRequest   // TODO
	IDTanShrinkWhiteListRequest   // TODO
	IDTanSetTagListRequest
	IDChangeRoomInfoRequest // TODO
	_
	_
	_
	_
	_
	_
	IDTanSetRoomDisplayModListRequest // TODO
	IDTanUpdateRoomPerformance        // TODO
)

// Server -> Client
const (
	IDTanLoginResponse uint16 = iota
	IDTanCreateRoomResponse
	_
	IDTanEnterRoomResponse
	IDTanNewGuestResponse
	IDTanLeaveRoomResponse
	IDTanKickOutResponse
	_
	_
	_
	_
	_
	IDTanSetTagListResponse
	IDChangeRoomInfoResponse // TODO
	_
	_
	_
	_
	_
	_
	IDTanSetRoomDisplayModListResponse // TODO
)

// Client <-> Server
const IDTanNotifyServerReady uint16 = 7
