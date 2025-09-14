package login

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
)

// readRaknetPacket ..
func (d *Dialer) readRaknetPacket(decoder *packet.Decoder) (pk packet.Packet, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("readRaknetPacket: %v", r)
		}
	}()

	pkData, err := decoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("readRaknetPacket: %v", err)
	}

	buf := bytes.NewBuffer(pkData)
	reader := encoding.NewReader(buf)

	header := packet.Header{}
	err = header.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("readRaknetPacket: %v", err)
	}

	pk = packet.NewServerPool()[header.PacketID]
	if pk == nil {
		return nil, fmt.Errorf("readRaknetPacket: unsupported packet %d; payload = %#v", header.PacketID, buf.Bytes())
	}

	pk.Marshal(reader)
	return pk, nil
}

// writeRaknetPacket ..
func (d *Dialer) writeRaknetPacket(encoder *packet.Encoder, pk packet.Packet) error {
	if pk == nil {
		return nil
	}

	buf := bytes.NewBuffer(nil)
	writer := encoding.NewWriter(buf, 0)
	pk.Marshal(writer)

	header := packet.Header{PacketID: pk.ID()}
	headerBuf := bytes.NewBuffer(nil)
	if err := header.Write(headerBuf); err != nil {
		return fmt.Errorf("writeRaknetPacket: %v", err)
	}

	err := encoder.Encode(append(headerBuf.Bytes(), buf.Bytes()...))
	if err != nil {
		return fmt.Errorf("writeRaknetPacket: %v", err)
	}

	return nil
}
