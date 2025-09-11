package packet

import (
	"fmt"
	"io"
	"slices"
)

// header is the header of compressed 'batches' from Minecraft.
var header = []byte{0xfe, 0xe3, 0x01}

// Decoder handles the decoding of Minecraft packets sent through an io.Reader. These packets in turn contain
// multiple compressed packets.
type Decoder struct {
	// pr holds a packetReader (and io.Reader) that packets are read from if the io.Reader passed to
	// NewDecoder implements the packetReader interface.
	pr      packetReader
	encrypt *encrypt
}

// packetReader is used to read packets immediately instead of copying them in a buffer first. This is a
// specific case made to reduce RAM usage.
type packetReader interface {
	ReadPacket() ([]byte, error)
}

// NewDecoder returns a new decoder decoding data from the io.Reader passed. One read call from the reader is
// assumed to consume an entire packet.
func NewDecoder(reader io.Reader) (result *Decoder, err error) {
	pr, ok := reader.(packetReader)
	if !ok {
		return nil, fmt.Errorf("NewDecoder: Given reader not implements ReadPacket() method")
	}
	return &Decoder{pr: pr}, nil
}

// EnableEncryption enables encryption for the Decoder using the secret key bytes passed. Each packet received
// will be decrypted.
func (decoder *Decoder) EnableEncryption(encryptKeyBytes []byte, decryptKeyBytes []byte) (err error) {
	decoder.encrypt, err = newEncrypt(encryptKeyBytes, decryptKeyBytes)
	if err != nil {
		return fmt.Errorf("EnableEncryption: %v", err)
	}
	return nil
}

// Decode decodes one 'packet' from the io.Reader passed in NewDecoder(), producing a slice of packets that it
// held and an error if not successful.
func (decoder *Decoder) Decode() (packetData []byte, err error) {
	var data []byte

	data, err = decoder.pr.ReadPacket()
	if err != nil {
		return nil, fmt.Errorf("read batch: %w", err)
	}

	if len(data) < 3 {
		return nil, nil
	}
	if !slices.Equal(data[0:3], header) {
		return nil, fmt.Errorf("decode batch: invalid header %x, expected %x", data[0:3], header)
	}

	data = data[3:]
	if decoder.encrypt != nil {
		data = decoder.encrypt.decrypt(data)
	}

	return data, nil
}
