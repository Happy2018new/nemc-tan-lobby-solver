package packet

import (
	"fmt"
	"io"
)

// Encoder handles the encoding of Minecraft packets that are sent to an io.Writer. The packets are compressed
// and optionally encoded before they are sent to the io.Writer.
type Encoder struct {
	w       io.Writer
	encrypt *encrypt
}

// NewEncoder returns a new Encoder for the io.Writer passed. Each final packet produced by the Encoder is
// sent with a single call to io.Writer.Write().
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// EnableEncryption enables encryption for the Encoder using the secret key bytes passed. Each packet sent
// after encryption is enabled will be encrypted.
func (encoder *Encoder) EnableEncryption(encryptKeyBytes []byte, decryptKeyBytes []byte) (err error) {
	encoder.encrypt, err = newEncrypt(encryptKeyBytes, decryptKeyBytes)
	if err != nil {
		return fmt.Errorf("EnableEncryption: %v", err)
	}
	return nil
}

// Encode encodes the packets passed. It writes all of them as a single packet which is  compressed and
// optionally encrypted.
func (encoder *Encoder) Encode(packetData []byte) error {
	prepend := append([]byte{}, header...)
	data := append(prepend, packetData...)
	if encoder.encrypt != nil {
		// If the encryption session is not nil, encryption is enabled, meaning we should encrypt the
		// compressed data of this packet.
		data = encoder.encrypt.encrypt(data)
	}
	if _, err := encoder.w.Write(data); err != nil {
		return fmt.Errorf("write batch: %w", err)
	}
	return nil
}
