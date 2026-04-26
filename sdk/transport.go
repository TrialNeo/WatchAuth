package sdk

import (
	"bytes"
	ciph "github.com/WatchAuth/sdk/crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type encryptedTransport struct {
	config Config
	client *http.Client
}

func newTransport(cfg Config) *encryptedTransport {
	return &encryptedTransport{config: cfg, client: &http.Client{}}
}

// serverResponse is the standard API response wrapper
type serverResponse struct {
	Code uint            `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data,omitempty"`
}

// do sends an encrypted request and decrypts the response.
func (t *encryptedTransport) do(method, path string, reqBody, respData interface{}) error {
	bodyBytes, _ := json.Marshal(reqBody)

	// Encrypt request body
	envelope, sessionKey, err := t.encrypt(bodyBytes)
	if err != nil {
		return fmt.Errorf("encrypt request: %w", err)
	}

	// Build outer request
	outerReq := struct {
		AppID string `json:"appId"`
		EncType uint8 `json:"encType"`
		ciph.EncryptedMessage
	}{
		AppID:            t.config.AppID,
		EncType:          envelope.EncType,
		EncryptedMessage: *envelope,
	}
	payload, _ := json.Marshal(outerReq)

	// Send HTTP request
	url := t.config.ServerURL + path
	httpReq, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	// Parse standard server response
	var sr serverResponse
	if err := json.Unmarshal(respBytes, &sr); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}
	if sr.Code != 0 {
		return fmt.Errorf("api error %d: %s", sr.Code, sr.Msg)
	}
	if respData == nil || sr.Data == nil {
		return nil
	}

	// Decrypt response data
	var respEnvelope ciph.EncryptedMessage
	if err := json.Unmarshal(sr.Data, &respEnvelope); err != nil {
		return fmt.Errorf("parse encrypted response: %w", err)
	}

	plaintext, err := t.decrypt(&respEnvelope, sessionKey)
	if err != nil {
		return fmt.Errorf("decrypt response: %w", err)
	}

	if err := json.Unmarshal(plaintext, respData); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	return nil
}

// encrypt creates an encrypted envelope. For RSA, returns the session key for response decryption.
func (t *encryptedTransport) encrypt(plaintext []byte) (*ciph.EncryptedMessage, []byte, error) {
	switch t.config.EncType {
	case ciph.EncTypeRSA:
		// RSA: generate AES session key, RSA-encrypt it, AES-GCM encrypt body
		aesKey := make([]byte, 32)
		if _, err := rand.Read(aesKey); err != nil {
			return nil, nil, err
		}
		pubKey, err := x509.ParsePKIXPublicKey([]byte(t.config.AppKey))
		if err != nil {
			return nil, nil, fmt.Errorf("rsa: invalid public key: %w", err)
		}
		rsaPub := pubKey.(*rsa.PublicKey)
		encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, aesKey, nil)
		if err != nil {
			return nil, nil, err
		}
		aesCiph, _ := ciph.NewCipher(ciph.EncTypeAESGCM, aesKey)
		inner, err := aesCiph.Encrypt(plaintext)
		if err != nil {
			return nil, nil, err
		}
		return &ciph.EncryptedMessage{
			EncType:      ciph.EncTypeRSA,
			Ciphertext:   inner.Ciphertext,
			Nonce:        inner.Nonce,
			EncryptedKey: base64.StdEncoding.EncodeToString(encryptedKey),
		}, aesKey, nil

	default:
		// AES-GCM or QQTea: use shared key directly
		c, err := ciph.NewCipher(t.config.EncType, []byte(t.config.AppKey))
		if err != nil {
			return nil, nil, err
		}
		env, err := c.Encrypt(plaintext)
		return env, nil, err
	}
}

// decrypt decrypts a response envelope. For RSA, uses the session key from the request.
func (t *encryptedTransport) decrypt(msg *ciph.EncryptedMessage, sessionKey []byte) ([]byte, error) {
	switch t.config.EncType {
	case ciph.EncTypeRSA:
		if sessionKey == nil {
			return nil, fmt.Errorf("rsa: no session key for response")
		}
		aesCiph, _ := ciph.NewCipher(ciph.EncTypeAESGCM, sessionKey)
		return aesCiph.Decrypt(msg)
	default:
		c, err := ciph.NewCipher(t.config.EncType, []byte(t.config.AppKey))
		if err != nil {
			return nil, err
		}
		return c.Decrypt(msg)
	}
}
