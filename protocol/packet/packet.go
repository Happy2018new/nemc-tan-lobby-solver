package packet

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/encoding"
)

const (
	BoundTypeServer uint8 = iota
	BoundTypeClient
)

// Packet represents a packet that may be sent over a Minecraft network connection.
// The packet hold some methods to allows you to do the following things.
//   - Get the ID of this packet.
//   - Get the resource string of this packet.
//   - Figure out this packet is used to send to server or send to client.
//   - Encode this packet to binary and decode this packet from binary.
type Packet interface {
	// ID returns the ID of the packet.
	// All of these identifiers of packets may be found in id.go.
	// Note that there a multiple id.go.
	ID() uint16
	// BoundType returns the bound type of the packet.
	// If return 0 (BoundTypeServer), it means this packet is send from client to server.
	// Otherwise (1; BoundTypeClient) this packet is send from server to client.
	BoundType() uint8
	// Marshal encodes or decodes a Packet, depending on the encoding.IO
	// implementation passed. When passing a encoding.Writer, Marshal will
	// encode the Packet into its binary representation and write it to the
	// encoding.Writer. On the other hand, when passing a encoding.Reader,
	// Marshal will decode the bytes from the reader into the Packet.
	Marshal(io encoding.IO)
}

// Header is the header of a packet. It exists out of a single varuint32 which is composed of a packet ID and
// a sender and target sub client ID. These IDs are used for split screen functionality.
type Header struct {
	PacketID uint16
}

// Write writes the header as a single varuint32 to buf.
func (header *Header) Write(w io.Writer) error {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, header.PacketID)
	if _, err := w.Write(result); err != nil {
		return fmt.Errorf("Write: %v", err)
	}
	return nil
}

// Read reads a varuint32 from buf and sets the corresponding values to the Header.
func (header *Header) Read(r io.Reader) error {
	headerBytes := make([]byte, 2)
	if _, err := r.Read(headerBytes); err != nil {
		return fmt.Errorf("Read: %v", err)
	}
	header.PacketID = binary.BigEndian.Uint16(headerBytes)
	return nil
}
