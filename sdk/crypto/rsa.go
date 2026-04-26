package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

type rsaCipher struct {
	key []byte
}

func newRSACipher(key []byte) (*rsaCipher, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("rsa: empty key")
	}
	return &rsaCipher{key: key}, nil
}

func (r *rsaCipher) Encrypt(plaintext []byte) (*EncryptedMessage, error) {
	pubKey, err := x509.ParsePKIXPublicKey(r.key)
	if err != nil {
		return nil, fmt.Errorf("rsa: invalid public key: %w", err)
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("rsa: key is not an RSA public key")
	}
	// Generate random AES-256 key
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, err
	}
	// RSA-OAEP encrypt the AES key
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa encrypt key: %w", err)
	}
	// AES-GCM encrypt the payload with the AES key
	aesCiph, err := newAESGCMCipher(aesKey)
	if err != nil {
		return nil, err
	}
	inner, err := aesCiph.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return &EncryptedMessage{
		EncType:      EncTypeRSA,
		Ciphertext:   inner.Ciphertext,
		Nonce:        inner.Nonce,
		EncryptedKey: base64.StdEncoding.EncodeToString(encryptedKey),
	}, nil
}

func (r *rsaCipher) Decrypt(msg *EncryptedMessage) ([]byte, error) {
	return nil, fmt.Errorf("rsa: client does not have private key for decryption; use client-side session key for response")
}
