package crypto

import (
	"encoding/base64"
	"fmt"
)

type EncryptedMessage struct {
	EncType      uint8  `json:"encType"`
	Ciphertext   string `json:"ciphertext"`
	Nonce        string `json:"nonce,omitempty"`
	EncryptedKey string `json:"ek,omitempty"`
}

type Cipher interface {
	Encrypt(plaintext []byte) (*EncryptedMessage, error)
	Decrypt(msg *EncryptedMessage) ([]byte, error)
}

const (
	EncTypeNone   = 0
	EncTypeRSA    = 1
	EncTypeAESGCM = 2
	EncTypeQQTea  = 3
)

func NewCipher(encType uint8, key []byte) (Cipher, error) {
	switch encType {
	case EncTypeNone:
		return &noneCipher{}, nil
	case EncTypeRSA:
		return newRSACipher(key)
	case EncTypeAESGCM:
		return newAESGCMCipher(key)
	case EncTypeQQTea:
		return newQQTeaCipher(key)
	default:
		return nil, fmt.Errorf("unknown encType: %d", encType)
	}
}

type noneCipher struct{}

func (n *noneCipher) Encrypt(plaintext []byte) (*EncryptedMessage, error) {
	return &EncryptedMessage{EncType: 0, Ciphertext: base64.StdEncoding.EncodeToString(plaintext)}, nil
}

func (n *noneCipher) Decrypt(msg *EncryptedMessage) ([]byte, error) {
	return base64.StdEncoding.DecodeString(msg.Ciphertext)
}
