package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
)

type OnceReader struct {
	*bytes.Buffer
}

func NewOnceReader(data []byte) *OnceReader {
	return &OnceReader{bytes.NewBuffer(data)}
}

func (o *OnceReader) ReadPacket() ([]byte, error) {
	return o.Bytes(), nil
}

func main() {
	wrapper := bunker.NewAccessWrapper(
		"AUTH SERVER ADDRESS",
		"YOUR FB TOKEN",
		"PE AUTH (CAN BE EMPTY)",
		"SA AUTH (CAN BE EMPTY)",
	)

	str := ``
	bs, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	dec, err := packet.NewDecoder(NewOnceReader(bs))
	if err != nil {
		panic(err)
	}
	bs, err = dec.Decode()
	if err != nil {
		panic(err)
	}
	r := encoding.NewReader(bytes.NewBuffer(bs[2:]))
	loginPacket := new(packet.TanLoginRequest)
	loginPacket.Marshal(r)

	tanLobbyDebugResp, err := wrapper.GetDebug(
		``,
		loginPacket.Rand,
	)
	if err != nil {
		panic(err)
	}
	if !tanLobbyDebugResp.Success {
		panic(tanLobbyDebugResp.ErrorInfo)
	}

	{
		str := ``
		bs, err := hex.DecodeString(str)
		if err != nil {
			panic(err)
		}
		dec, err := packet.NewDecoder(NewOnceReader(bs))
		if err != nil {
			panic(err)
		}
		dec.EnableEncryption(tanLobbyDebugResp.EncryptKeyBytes, tanLobbyDebugResp.EncryptKeyBytes)
		bs, err = dec.Decode()
		if err != nil {
			panic(err)
		}
		r := encoding.NewReader(bytes.NewBuffer(bs[2:]))
		pk := new(packet.TanCreateRoomRequest)
		pk.Marshal(r)
		fmt.Println(pk)
	}
}
