package packet

import (
	"fmt"

	"github.com/database64128/chacha8-go/chacha8"
)

// WebRTCCipherNonce ..
var WebRTCCipherNonce = []byte{0x31, 0x36, 0x33, 0x20, 0x4e, 0x65, 0x74, 0x45, 0x61, 0x73, 0x65, 0x0a}

// encrypt holds an encryption session with several fields required to encrypt and/or decrypt incoming
// packets. It may be initialised using the encrypt key bytes and decrypt key bytes.
type encrypt struct {
	encrypter *chacha8.Cipher
	decrypter *chacha8.Cipher
}

// newEncrypt init a new encrypt/decrypt implements for current session.
// encryptKeyBytes and decryptKeyBytes is the key to encrypt and decrypt data.
func newEncrypt(encryptKeyBytes []byte, decryptKeyBytes []byte) (result *encrypt, err error) {
	encrypter, err := chacha8.NewUnauthenticatedCipher(encryptKeyBytes, WebRTCCipherNonce)
	if err != nil {
		return nil, fmt.Errorf("newEncrypt: %v", err)
	}
	decrypter, err := chacha8.NewUnauthenticatedCipher(decryptKeyBytes, WebRTCCipherNonce)
	if err != nil {
		return nil, fmt.Errorf("newEncrypt: %v", err)
	}
	return &encrypt{
		encrypter: encrypter,
		decrypter: decrypter,
	}, nil
}

// encrypt encrypts the data passed.
func (encrypt *encrypt) encrypt(data []byte) []byte {
	encrypt.encrypter.XORKeyStream(data, data)
	return data
}

// decrypt decrypts the data passed.
func (encrypt *encrypt) decrypt(data []byte) []byte {
	encrypt.decrypter.XORKeyStream(data, data)
	return data
}
