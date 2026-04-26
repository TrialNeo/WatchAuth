package crypto

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

type qqteaCipher struct {
	key [4]uint32
}

func newQQTeaCipher(key []byte) (*qqteaCipher, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("qqtea: key must be 16 bytes, got %d", len(key))
	}
	var k [4]uint32
	k[0] = binary.LittleEndian.Uint32(key[0:4])
	k[1] = binary.LittleEndian.Uint32(key[4:8])
	k[2] = binary.LittleEndian.Uint32(key[8:12])
	k[3] = binary.LittleEndian.Uint32(key[12:16])
	return &qqteaCipher{key: k}, nil
}

const teaDelta = 0x9E3779B9

func teaEncrypt(v0, v1 uint32, k [4]uint32) (uint32, uint32) {
	sum := uint32(0)
	for i := 0; i < 32; i++ {
		v0 += ((v1 << 4) + k[0]) ^ (v1 + sum) ^ ((v1 >> 5) + k[1])
		sum += teaDelta
		v1 += ((v0 << 4) + k[2]) ^ (v0 + sum) ^ ((v0 >> 5) + k[3])
	}
	return v0, v1
}

func teaDecrypt(v0, v1 uint32, k [4]uint32) (uint32, uint32) {
	var sum uint32 = 0xC6EF3720 // delta * 32 with natural uint32 overflow
	for i := 0; i < 32; i++ {
		v1 -= ((v0 << 4) + k[2]) ^ (v0 + sum) ^ ((v0 >> 5) + k[3])
		sum -= teaDelta
		v0 -= ((v1 << 4) + k[0]) ^ (v1 + sum) ^ ((v1 >> 5) + k[1])
	}
	return v0, v1
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+padding)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(padding)
	}
	return out
}

func pkcs7Unpad(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return data
	}
	return data[:len(data)-padding]
}

func (q *qqteaCipher) Encrypt(plaintext []byte) (*EncryptedMessage, error) {
	padded := pkcs7Pad(plaintext, 8)
	out := make([]byte, len(padded))
	for i := 0; i < len(padded); i += 8 {
		v0 := binary.LittleEndian.Uint32(padded[i : i+4])
		v1 := binary.LittleEndian.Uint32(padded[i+4 : i+8])
		v0, v1 = teaEncrypt(v0, v1, q.key)
		binary.LittleEndian.PutUint32(out[i:i+4], v0)
		binary.LittleEndian.PutUint32(out[i+4:i+8], v1)
	}
	return &EncryptedMessage{
		EncType:    EncTypeQQTea,
		Ciphertext: base64.StdEncoding.EncodeToString(out),
	}, nil
}

func (q *qqteaCipher) Decrypt(msg *EncryptedMessage) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(msg.Ciphertext)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%8 != 0 {
		return nil, fmt.Errorf("qqtea: ciphertext length not multiple of 8")
	}
	out := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += 8 {
		v0 := binary.LittleEndian.Uint32(ciphertext[i : i+4])
		v1 := binary.LittleEndian.Uint32(ciphertext[i+4 : i+8])
		v0, v1 = teaDecrypt(v0, v1, q.key)
		binary.LittleEndian.PutUint32(out[i:i+4], v0)
		binary.LittleEndian.PutUint32(out[i+4:i+8], v1)
	}
	return pkcs7Unpad(out), nil
}
