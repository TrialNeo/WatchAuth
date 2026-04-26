package controller

import (
	"Diggpher/internal/service"
	"Diggpher/internal/service/errMsg"
	"Diggpher/pkg/crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type SdkController struct {
	Svc *service.SDKService
}

// decryptBody parses the encrypted envelope and decrypts into target struct.
// Returns the cipher for reusing in encrypting the response.
// For RSA: the session AES key from the request is reused for the response.
// For AES-GCM/QQTea: the app's stored key is used directly for both directions.
func (s *SdkController) decryptBody(c *fiber.Ctx, target interface{}) (crypto.Cipher, uint) {
	var env struct {
		AppID   string `json:"appId"`
		EncType uint8  `json:"encType"`
		crypto.EncryptedMessage
	}
	if err := c.BodyParser(&env); err != nil {
		return nil, errMsg.ERRORInvalidParams
	}
	app, code := s.Svc.GetAppConfig(env.AppID)
	if code != errMsg.SUCCESS {
		return nil, code
	}

	var ciph crypto.Cipher
	var plaintext []byte
	var err error

	if env.EncType == crypto.EncTypeRSA {
		ciph, plaintext, err = decryptRSABody([]byte(app.SecretKeys), &env.EncryptedMessage)
	} else {
		ciph, err = crypto.NewCipher(env.EncType, []byte(app.SecretKeys))
		if err == nil {
			plaintext, err = ciph.Decrypt(&env.EncryptedMessage)
		}
	}
	if err != nil {
		return nil, errMsg.ERROR
	}
	if err := json.Unmarshal(plaintext, target); err != nil {
		return nil, errMsg.ERRORInvalidParams
	}
	return ciph, errMsg.SUCCESS
}

// decryptRSABody decrypts the RSA encrypted message and returns an AES-GCM cipher
// using the session key derived from the RSA key exchange.
func decryptRSABody(key []byte, msg *crypto.EncryptedMessage) (crypto.Cipher, []byte, error) {
	privKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, nil, err
	}
	rsaPriv, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("invalid RSA private key")
	}
	ek, err := base64.StdEncoding.DecodeString(msg.EncryptedKey)
	if err != nil {
		return nil, nil, err
	}
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPriv, ek, nil)
	if err != nil {
		return nil, nil, err
	}
	ciph, err := crypto.NewCipher(crypto.EncTypeAESGCM, aesKey)
	if err != nil {
		return nil, nil, err
	}
	plaintext, err := ciph.Decrypt(msg)
	return ciph, plaintext, err
}

// encryptResp encrypts the response data using the cipher and returns it through RespondIMP.
func encryptResp(re *RespondIMP, ciph crypto.Cipher, code uint, data interface{}) error {
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	dataBytes, _ := json.Marshal(data)
	encrypted, err := ciph.Encrypt(dataBytes)
	if err != nil {
		return re.withCode(errMsg.ERROR).Respond(nil)
	}
	return re.withCode(code).Respond(encrypted)
}
