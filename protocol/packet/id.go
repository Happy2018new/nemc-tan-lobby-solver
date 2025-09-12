package packet

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
	IDTanSetTagListRequest        // TODO
	IDChangeRoomInfoRequest       // TODO
	_
	_
	_
	_
	_
	_
	IDTanSetRoomDisplayModListRequest // TODO
	IDTanUpdateRoomPerformance        // TODO
)

const (
	IDTanLoginResponse uint16 = iota
	IDTanCreateRoomResponse
	_
	IDTanEnterRoomResponse
	IDTanNewGuestResponse // TODO
	IDTanLeaveRoomResponse
	IDTanKickOutResponse
	IDTanNotifyServerReady
	_
	_
	_
	_
	IDTanSetTagListResponse  // TODO
	IDChangeRoomInfoResponse // TODO
	_
	_
	_
	_
	_
	_
	IDTanSetRoomDisplayModListResponse // TODO
)
