package crypto

import (
	"crypto/aes"
	"crypto/cipher"
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
		privKey, err2 := x509.ParsePKCS8PrivateKey(r.key)
		if err2 != nil {
			return nil, fmt.Errorf("rsa: invalid key: %w", err2)
		}
		pubKey = privKey.(*rsa.PrivateKey).Public()
	}
	rsaPubKey := pubKey.(*rsa.PublicKey)

	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("rsa: generate aes key: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa encrypt key: %w", err)
	}

	return &EncryptedMessage{
		EncType:      EncTypeRSA,
		Ciphertext:   base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:        base64.StdEncoding.EncodeToString(nonce),
		EncryptedKey: base64.StdEncoding.EncodeToString(encryptedKey),
	}, nil
}

func (r *rsaCipher) Decrypt(msg *EncryptedMessage) ([]byte, error) {
	encryptedKey, err := base64.StdEncoding.DecodeString(msg.EncryptedKey)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(msg.Nonce)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(msg.Ciphertext)
	if err != nil {
		return nil, err
	}

	privKey, err := x509.ParsePKCS8PrivateKey(r.key)
	if err != nil {
		return nil, fmt.Errorf("rsa: invalid private key: %w", err)
	}
	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("rsa: key is not a private key")
	}

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivKey, encryptedKey, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa decrypt key: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa decrypt data: %w", err)
	}
	return plaintext, nil
}
