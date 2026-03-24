package crypto

import (
	"Diggpher/global"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
)

// PswEnc 密码加密，
func PswEnc(psw string) string {
	hash := crypto.Hash.New(crypto.SHA256)
	hash.Write([]byte(psw))
	signature, err := rsa.SignPKCS1v15(rand.Reader, global.RsaPriPem, crypto.SHA256, hash.Sum(hash.Sum(nil)))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// PswDec 密码解密，
func PswDec(psw string) string {
	hash := crypto.Hash.New(crypto.SHA256)
	hash.Write([]byte(psw))
	signature, err := rsa.SignPKCS1v15(rand.Reader, global.RsaPriPem, crypto.SHA256, hash.Sum(hash.Sum(nil)))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(signature)
}
