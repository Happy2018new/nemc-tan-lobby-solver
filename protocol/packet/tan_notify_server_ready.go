package packet

import "github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"

// TanNotifyServerReady ..
type TanNotifyServerReady struct {
	ServerAddress         string
	ServerRaknetGuid      string
	RTCRoomID             string
	NetherNetID           string
	WebRTCCompressEnabled bool
}

func (*TanNotifyServerReady) ID() uint16 {
	return IDTanNotifyServerReady
}

func (*TanNotifyServerReady) BoundType() uint8 {
	return BoundTypeBoth
}

func (t *TanNotifyServerReady) Marshal(io encoding.IO) {
	io.StringUTF(&t.ServerAddress)
	io.StringUTF(&t.ServerRaknetGuid)
	io.StringUTF(&t.RTCRoomID)
	io.StringUTF(&t.NetherNetID)
	io.Bool(&t.WebRTCCompressEnabled)
}
