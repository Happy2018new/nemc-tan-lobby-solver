package service

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/packet"
)

// readPacketWithContext ..
func readPacketWithContext(ctx context.Context, conn net.Conn, decoder *packet.Decoder) (packet.Packet, error) {
	pkChannel := make(chan packet.Packet, 1)
	errChannel := make(chan error, 1)
	defer close(pkChannel)
	defer close(errChannel)

	go func() {
		pk, err := readPacket(decoder)
		if err != nil {
			errChannel <- err
			return
		}
		pkChannel <- pk
	}()

	select {
	case <-ctx.Done():
		_ = conn.Close()
		return nil, fmt.Errorf("readPacketWithContext: %v", ctx.Err())
	case err := <-errChannel:
		return nil, fmt.Errorf("readPacketWithContext: %v", err)
	case pk := <-pkChannel:
		return pk, nil
	}
}

// readPacket ..
func readPacket(decoder *packet.Decoder) (pk packet.Packet, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("readPacket: %v", r)
		}
	}()

	pkData, err := decoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("readPacket: %v", err)
	}

	buf := bytes.NewBuffer(pkData)
	reader := encoding.NewReader(buf)

	header := packet.Header{}
	err = header.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("readPacket: %v", err)
	}

	pk = packet.NewServerPool()[header.PacketID]
	if pk == nil {
		return nil, fmt.Errorf("readPacket: unsupported packet %d; payload = %#v", header.PacketID, buf.Bytes())
	}

	pk.Marshal(reader)
	return pk, nil
}

// writePacket ..
func writePacket(encoder *packet.Encoder, pk packet.Packet) error {
	if pk == nil {
		return nil
	}

	buf := bytes.NewBuffer(nil)
	writer := encoding.NewWriter(buf, 0)
	pk.Marshal(writer)

	header := packet.Header{PacketID: pk.ID()}
	headerBuf := bytes.NewBuffer(nil)
	if err := header.Write(headerBuf); err != nil {
		return fmt.Errorf("writePacket: %v", err)
	}

	err := encoder.Encode(append(headerBuf.Bytes(), buf.Bytes()...))
	if err != nil {
		return fmt.Errorf("writePacket: %v", err)
	}

	return nil
}
