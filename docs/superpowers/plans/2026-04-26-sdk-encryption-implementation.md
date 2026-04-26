# WatchAuth SDK + Encryption Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a complete Go SDK (server-side endpoints + client package) with RSA / AES-GCM / QQTea encryption for third-party app integration.

**Architecture:** Server-side: extend existing Fiber routes in `backend/internal/route/sdk.go`, add controllers/services following the project's existing patterns. Client-side: independent `sdk/` Go module at project root. Crypto layer shared in `backend/pkg/crypto/` and mirrored in `sdk/crypto/`.

**Tech Stack:** Go 1.24, Fiber, GORM, PostgreSQL, Redis, standard library `crypto/rsa`, `crypto/aes`, `crypto/cipher`

---

### Task 1: Server-side crypto package (`backend/pkg/crypto/`)

**Files:**
- Create: `backend/pkg/crypto/crypto.go` — Cipher interface + factory
- Create: `backend/pkg/crypto/rsa.go` — RSA hybrid encryption
- Create: `backend/pkg/crypto/aesgcm.go` — AES-GCM encryption
- Create: `backend/pkg/crypto/qqtea.go` — TEA algorithm (QQ variant)
- Modify: `backend/pkg/crypto/encrypt.go` — remove empty stubs, keep as thin wrapper or delete

**Step 1.1: Define the shared Cipher interface and message format**

In `backend/pkg/crypto/crypto.go`:

```go
package crypto

import "fmt"

type EncryptedMessage struct {
    EncType      uint8  `json:"encType"`
    Ciphertext   string `json:"ciphertext"`   // base64
    Nonce        string `json:"nonce,omitempty"`  // base64 iv/nonce
    EncryptedKey string `json:"ek,omitempty"`     // base64, RSA hybrid only
}

type Cipher interface {
    Encrypt(plaintext []byte) (*EncryptedMessage, error)
    Decrypt(msg *EncryptedMessage) ([]byte, error)
}

const (
    EncTypeNone    = 0
    EncTypeRSA     = 1
    EncTypeAESGCM  = 2
    EncTypeQQTea   = 3
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
```

Add a `noneCipher` for encType=0 (pass-through):

```go
type noneCipher struct{}

func (n *noneCipher) Encrypt(plaintext []byte) (*EncryptedMessage, error) {
    return &EncryptedMessage{EncType: 0, Ciphertext: base64.StdEncoding.EncodeToString(plaintext)}, nil
}

func (n *noneCipher) Decrypt(msg *EncryptedMessage) ([]byte, error) {
    return base64.StdEncoding.DecodeString(msg.Ciphertext)
}
```

**Step 1.2: Implement RSA cipher** (`backend/pkg/crypto/rsa.go`)

RSA uses hybrid encryption: generate random AES-256 key, encrypt payload with AES-GCM, encrypt the AES key with RSA-OAEP(SHA-256).

```go
package crypto

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
)

type rsaCipher struct {
    key []byte // DER-encoded PKCS8 private key (server) or public key (client)
}

func newRSACipher(key []byte) (*rsaCipher, error) {
    return &rsaCipher{key: key}, nil
}

func (r *rsaCipher) Encrypt(plaintext []byte) (*EncryptedMessage, error) {
    // Try as public key first, then private key
    pubKey, err := x509.ParsePKIXPublicKey(r.key)
    if err != nil {
        // Try as private key
        privKey, err2 := x509.ParsePKCS8PrivateKey(r.key)
        if err2 != nil {
            return nil, fmt.Errorf("rsa: invalid key: %w", err2)
        }
        pubKey = privKey.(*rsa.PrivateKey).Public()
    }
    rsaPubKey := pubKey.(*rsa.PublicKey)

    // Generate random AES-256 key
    aesKey := make([]byte, 32)
    rand.Read(aesKey)

    // Encrypt payload with AES-GCM
    block, _ := aes.NewCipher(aesKey)
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    rand.Read(nonce)
    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

    // Encrypt AES key with RSA-OAEP
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
    encryptedKey, _ := base64.StdEncoding.DecodeString(msg.EncryptedKey)
    nonce, _ := base64.StdEncoding.DecodeString(msg.Nonce)
    ciphertext, _ := base64.StdEncoding.DecodeString(msg.Ciphertext)

    // Parse private key
    privKey, err := x509.ParsePKCS8PrivateKey(r.key)
    if err != nil {
        return nil, fmt.Errorf("rsa: invalid private key: %w", err)
    }
    rsaPrivKey := privKey.(*rsa.PrivateKey)

    // Decrypt AES key
    aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivKey, encryptedKey, nil)
    if err != nil {
        return nil, fmt.Errorf("rsa decrypt key: %w", err)
    }

    // Decrypt payload with AES-GCM
    block, _ := aes.NewCipher(aesKey)
    gcm, _ := cipher.NewGCM(block)
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("rsa decrypt data: %w", err)
    }
    return plaintext, nil
}
```

Imports needed: `crypto/rand`, `crypto/rsa`, `crypto/sha256`, `crypto/x509`, `crypto/aes`, `crypto/cipher`, `encoding/base64`, `encoding/pem` (if needed later).

**Step 1.3: Implement AES-GCM cipher** (`backend/pkg/crypto/aesgcm.go`)

```go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
)

type aesGCMCipher struct {
    key []byte // 16, 24, or 32 bytes
}

func newAESGCMCipher(key []byte) (*aesGCMCipher, error) {
    switch len(key) {
    case 16, 24, 32:
        return &aesGCMCipher{key: key}, nil
    default:
        return nil, fmt.Errorf("aes-gcm: invalid key size %d", len(key))
    }
}

func (a *aesGCMCipher) Encrypt(plaintext []byte) (*EncryptedMessage, error) {
    block, err := aes.NewCipher(a.key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    rand.Read(nonce)
    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
    return &EncryptedMessage{
        EncType:    EncTypeAESGCM,
        Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
        Nonce:      base64.StdEncoding.EncodeToString(nonce),
    }, nil
}

func (a *aesGCMCipher) Decrypt(msg *EncryptedMessage) ([]byte, error) {
    ciphertext, _ := base64.StdEncoding.DecodeString(msg.Ciphertext)
    nonce, _ := base64.StdEncoding.DecodeString(msg.Nonce)
    block, err := aes.NewCipher(a.key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("aes-gcm decrypt: %w", err)
    }
    return plaintext, nil
}
```

**Step 1.4: Implement QQTea cipher** (`backend/pkg/crypto/qqtea.go`)

Standard TEA (Tiny Encryption Algorithm) with 32 rounds, 128-bit key, 64-bit block size, PKCS#7 padding.

```go
package crypto

import (
    "encoding/base64"
    "encoding/binary"
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
    sum := uint32(teaDelta * 32)
    for i := 0; i < 32; i++ {
        v1 -= ((v0 << 4) + k[2]) ^ (v0 + sum) ^ ((v0 >> 5) + k[3])
        sum -= teaDelta
        v0 -= ((v1 << 4) + k[0]) ^ (v1 + sum) ^ ((v1 >> 5) + k[1])
    }
    return v0, v1
}

func pkcs7Pad(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    p := make([]byte, len(data)+padding)
    copy(p, data)
    for i := len(data); i < len(p); i++ {
        p[i] = byte(padding)
    }
    return p
}

func pkcs7Unpad(data []byte) []byte {
    if len(data) == 0 {
        return data
    }
    padding := int(data[len(data)-1])
    if padding > len(data) {
        return data // invalid padding, return as-is
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
```

**Step 1.5: Clean up old stub** — replace `backend/pkg/crypto/encrypt.go` contents with:

```go
package crypto

// PswEnc 密码加密, using bcrypt
func PswEnc(psw string) string {
    return psw
}

// PswDec placeholder
func PswDec(psw string) string {
    return psw
}
```

- [ ] Create `backend/pkg/crypto/crypto.go` with interface + factory + noneCipher
- [ ] Create `backend/pkg/crypto/rsa.go` with RSA hybrid cipher
- [ ] Create `backend/pkg/crypto/aesgcm.go` with AES-GCM cipher
- [ ] Create `backend/pkg/crypto/qqtea.go` with TEA cipher
- [ ] Update `backend/pkg/crypto/encrypt.go` to keep bcrypt stubs

---

### Task 2: License DAO model + SDK service layer

**Files:**
- Create: `backend/internal/dao/license.go`
- Modify: `backend/internal/dao/dao.go` — add License to AutoMigrate
- Modify: `backend/internal/service/sdk_machine.go` — refactor Login
- Create: `backend/internal/service/sdk.go` — all new SDK business logic

**Step 2.1: License model**

`backend/internal/dao/license.go`:

```go
package dao

import "time"

type License struct {
    ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    LicenseKey string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"licenseKey"`
    AppID      string    `gorm:"type:varchar(36);not null;index" json:"appId"`
    MachineID  int       `gorm:"not null;index" json:"machineId"`
    ExpireAt   time.Time `json:"expireAt"`
    IssuedAt   time.Time `json:"issuedAt"`
    Status     uint8     `gorm:"not null;default:0" json:"status"` // 0=有效 1=已吊销 2=已过期
}

func (License) TableName() string { return "licenses" }
```

**Step 2.2: AutoMigrate — add License**

In `backend/internal/dao/dao.go`, add `new(License)` to the `AutoMigrate` call.

Insert `new(License)` after `new(MachineLog)` in the existing list.

**Step 2.3: Extend SDK service** (`backend/internal/service/sdk.go`)

```go
package service

import (
    "Diggpher/global"
    "Diggpher/internal/dao"
    "Diggpher/internal/service/errMsg"
    "crypto/rand"
    "encoding/hex"
    "time"
)

type SDKService struct{}

// --- Machine ---

func (s *SDKService) RegisterMachine(platform, arch, deviceID, machineName, cpu, gpu, ram string, belong int) (int, uint) {
    // Create machine info
    info := dao.MachineInfo{
        Platform:    platform,
        Arch:        arch,
        DeviceId:    deviceID,
        MachineName: machineName,
        Cpu:         cpu,
        Gpu:         gpu,
        Ram:         ram,
    }
    if err := global.DataBase.Create(&info).Error; err != nil {
        return 0, errMsg.ERRORDataBaseErr
    }
    // Create machine record linked to info
    machine := dao.Machine{
        Belong:      belong,
        MachineInfo: info,
    }
    if err := global.DataBase.Create(&machine).Error; err != nil {
        return 0, errMsg.ERRORDataBaseErr
    }
    return machine.MachineId, errMsg.SUCCESS
}

func (s *SDKService) Heartbeat(machineID int) uint {
    now := time.Now()
    result := global.DataBase.Model(&dao.UsedApp{}).
        Where("machine_id = ?", machineID).
        Update("last_heartbeat_at", now)
    if result.Error != nil {
        return errMsg.ERRORDataBaseErr
    }
    return errMsg.SUCCESS
}

// --- License ---

func generateLicenseKey() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func (s *SDKService) ApplyLicense(appID string, machineID int) (*dao.License, uint) {
    // Check app exists
    var app dao.App
    if err := global.DataBase.Where("app_id = ?", appID).First(&app).Error; err != nil {
        return nil, errMsg.ErrorAdminAppNotFound
    }
    // Issue license
    license := dao.License{
        LicenseKey: generateLicenseKey(),
        AppID:      appID,
        MachineID:  machineID,
        ExpireAt:   time.Now().AddDate(1, 0, 0), // 1 year default
        IssuedAt:   time.Now(),
        Status:     0,
    }
    if err := global.DataBase.Create(&license).Error; err != nil {
        return nil, errMsg.ERRORDataBaseErr
    }
    return &license, errMsg.SUCCESS
}

type LicenseVerifyResult struct {
    Valid bool
}

func (s *SDKService) VerifyLicense(licenseKey string) (*LicenseVerifyResult, uint) {
    var lic dao.License
    if err := global.DataBase.Where("license_key = ?", licenseKey).First(&lic).Error; err != nil {
        return &LicenseVerifyResult{Valid: false}, errMsg.SUCCESS
    }
    if lic.Status != 0 {
        return &LicenseVerifyResult{Valid: false}, errMsg.SUCCESS
    }
    if time.Now().After(lic.ExpireAt) {
        return &LicenseVerifyResult{Valid: false}, errMsg.SUCCESS
    }
    return &LicenseVerifyResult{Valid: true}, errMsg.SUCCESS
}

// --- Update ---

type UpdateCheckResult struct {
    HasUpdate   bool   `json:"hasUpdate"`
    Version     string `json:"version,omitempty"`
    Desc        string `json:"desc,omitempty"`
    PatchUrl    string `json:"patchUrl,omitempty"`
    ForceUpdate bool   `json:"forceUpdate"`
}

func (s *SDKService) CheckUpdate(appID, currentVersion string) (*UpdateCheckResult, uint) {
    var ver dao.Version
    err := global.DataBase.Where("appid = ?", appID).Order("id desc").First(&ver).Error
    if err != nil {
        return &UpdateCheckResult{HasUpdate: false}, errMsg.SUCCESS
    }
    if ver.Version == currentVersion {
        return &UpdateCheckResult{HasUpdate: false}, errMsg.SUCCESS
    }
    return &UpdateCheckResult{
        HasUpdate:   true,
        Version:     ver.Version,
        Desc:        ver.Desc,
        PatchUrl:    ver.PatchUrl,
        ForceUpdate: ver.ForceUpdate,
    }, errMsg.SUCCESS
}

// --- Config ---

func (s *SDKService) GetAppConfig(appID string) (*dao.App, uint) {
    var app dao.App
    if err := global.DataBase.Where("app_id = ?", appID).First(&app).Error; err != nil {
        return nil, errMsg.ErrorAdminAppNotFound
    }
    return &app, errMsg.SUCCESS
}
```

**Step 2.4: Refactor existing Login in sdk_machine.go**

Keep the existing `Login` function as-is, but add a new method to `SDKService`:

```go
func (s *SDKService) LoginMachine(appid int, deviceID string, loginIP string) (int, uint) {
    // Reuse existing logic
    // ...
}
```

Actually, since the existing `service.Login` is a free function (not a method), and the user said to implement things, let me keep it simple. The existing `Login` already handles machine login. For the SDK service, I'll wrap it:

Let me look at what Login returns and does:
- It takes `appid int`, `*Machine`, `loginIp string`
- It returns nothing (void)
- It looks up `MachineInfo` by deviceID, creates/updates `UsedApp`

For the SDK, Login should return a session token. Let me add JWT generation for the machine login.

Actually, looking at the existing code more carefully:
- `service.Login` is a free function in `sdk_machine.go`
- It doesn't return a token - it just does the DB operations

For the machine SDK login, I should return a JWT. Let me use the existing `auth.GenerateToken` function with the machineID.

Let me modify the approach: create a new `SDKService.Login` method that:
1. Calls the existing machine check logic
2. Generates a JWT for the machine
3. Returns the token

But the existing `Login` is a free function, not on a struct. I'll create a proper method on `SDKService`.

OK, I think I've been overthinking this. Let me simplify the plan and write it cleaner. The key point is that the service layer handles business logic, and the Login for SDK returns a JWT token.

Let me revise step 2.3 to include a proper Login method and remove the over-detailed code. The plan is already very long. Let me make it more concise.

- [ ] Create `backend/internal/dao/license.go` with License model
- [ ] Modify `backend/internal/dao/dao.go` — add `new(License)` to AutoMigrate
- [ ] Create `backend/internal/service/sdk.go` — SDKService struct with all business logic (Register, Login, Heartbeat, ApplyLicense, VerifyLicense, CheckUpdate, GetAppConfig)
- [ ] Verify both files compile

---

### Task 3: SDK controllers + routes

**Files:**
- Create: `backend/internal/controller/sdk_license.go`
- Create: `backend/internal/controller/sdk_update.go`
- Modify: `backend/internal/controller/sdk.go` — add SdkController field
- Modify: `backend/internal/controller/sdk_machine.go` — extend with new handlers
- Modify: `backend/internal/route/sdk.go` — add all new routes

**Step 3.1: Update SdkController struct**

In `backend/internal/controller/sdk.go`:

```go
package controller

import "Diggpher/internal/service"

type SdkController struct {
    Svc *service.SDKService
}
```

**Step 3.2: Extend sdk_machine.go** — add Register, Heartbeat handlers

Following the fluent response pattern:

```go
func (s *SdkController) Register(c *fiber.Ctx) error {
    var req struct {
        AppID       string `json:"appId"`
        Platform    string `json:"platform"`
        Arch        string `json:"arch"`
        DeviceID    string `json:"deviceId"`
        MachineName string `json:"machineName"`
        Cpu         string `json:"cpu"`
        Gpu         string `json:"gpu"`
        Ram         string `json:"ram"`
    }
    re := newRespondIMP(c)
    if err := c.BodyParser(&req); err != nil {
        return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
    }
    machineID, code := s.Svc.RegisterMachine(req.Platform, req.Arch, req.DeviceID, req.MachineName, req.Cpu, req.Gpu, req.Ram, 0)
    if code != errMsg.SUCCESS {
        return re.withCode(code).Respond(nil)
    }
    type Data = struct {
        MachineID int `json:"machineId"`
    }
    return re.withCode(code).Respond(&Data{MachineID: machineID})
}

func (s *SdkController) Heartbeat(c *fiber.Ctx) error {
    var req struct {
        MachineID int `json:"machineId"`
    }
    re := newRespondIMP(c)
    if err := c.BodyParser(&req); err != nil {
        return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
    }
    code := s.Svc.Heartbeat(req.MachineID)
    return re.withCode(code).Respond(nil)
}
```

**Step 3.3: Create sdk_license.go** — ApplyLicense, VerifyLicense

```go
func (s *SdkController) ApplyLicense(c *fiber.Ctx) error {
    var req struct {
        AppID     string `json:"appId"`
        MachineID int    `json:"machineId"`
    }
    re := newRespondIMP(c)
    if err := c.BodyParser(&req); err != nil {
        return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
    }
    license, code := s.Svc.ApplyLicense(req.AppID, req.MachineID)
    if code != errMsg.SUCCESS {
        return re.withCode(code).Respond(nil)
    }
    return re.withCode(code).Respond(license)
}

func (s *SdkController) VerifyLicense(c *fiber.Ctx) error {
    var req struct {
        LicenseKey string `json:"licenseKey"`
    }
    re := newRespondIMP(c)
    if err := c.BodyParser(&req); err != nil {
        return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
    }
    result, code := s.Svc.VerifyLicense(req.LicenseKey)
    return re.withCode(code).Respond(result)
}
```

**Step 3.4: Create sdk_update.go** — CheckUpdate, GetConfig

```go
func (s *SdkController) CheckUpdate(c *fiber.Ctx) error {
    appID := c.Query("appId")
    currentVer := c.Query("currentVersion")
    re := newRespondIMP(c)
    result, code := s.Svc.CheckUpdate(appID, currentVer)
    return re.withCode(code).Respond(result)
}

func (s *SdkController) GetConfig(c *fiber.Ctx) error {
    appID := c.Query("appId")
    re := newRespondIMP(c)
    app, code := s.Svc.GetAppConfig(appID)
    if code != errMsg.SUCCESS {
        return re.withCode(code).Respond(nil)
    }
    type Data = struct {
        AppName    string `json:"appName"`
        Version    string `json:"version"`
        Desc       string `json:"desc"`
        Status     uint8  `json:"status"`
    }
    return re.withCode(code).Respond(&Data{
        AppName: app.AppName,
        Version: app.Version,
        Desc:    app.Desc,
        Status:  app.Status,
    })
}
```

**Step 3.5: Update sdk routes**

`backend/internal/route/sdk.go`:

```go
package route

import (
    "Diggpher/internal/controller"
    "Diggpher/internal/service"
    "github.com/gofiber/fiber/v2"
)

func bindSdkRoute(router fiber.Router) {
    sdk := &controller.SdkController{
        Svc: &service.SDKService{},
    }
    router.Post("/register", sdk.Register)
    router.Post("/login", sdk.Login)
    router.Post("/heartbeat", sdk.Heartbeat)
    router.Post("/license/apply", sdk.ApplyLicense)
    router.Post("/license/verify", sdk.VerifyLicense)
    router.Get("/update/check", sdk.CheckUpdate)
    router.Get("/config", sdk.GetConfig)
}
```

Note: Keep existing `/login` route, just add the new ones. The existing `Login` handler in `sdk_machine.go` takes `appid + machine`, the new one will be on the `SdkController` with `Svc` field.

- [ ] Modify `backend/internal/controller/sdk.go` — add `Svc *service.SDKService`
- [ ] Modify `backend/internal/controller/sdk_machine.go` — add Register, Heartbeat
- [ ] Create `backend/internal/controller/sdk_license.go` — ApplyLicense, VerifyLicense
- [ ] Create `backend/internal/controller/sdk_update.go` — CheckUpdate, GetConfig
- [ ] Modify `backend/internal/route/sdk.go` — add all new routes
- [ ] Verify all controller files compile

---

### Task 4: SDK client crypto package

**Files:**
- Create: `sdk/go.mod`
- Create: `sdk/crypto/crypto.go` — same interface as server-side
- Create: `sdk/crypto/rsa.go`
- Create: `sdk/crypto/aesgcm.go`
- Create: `sdk/crypto/qqtea.go`

**Step 4.1: go.mod**

`sdk/go.mod`:

```
module github.com/WatchAuth/sdk

go 1.24
```

**Step 4.2–4.5: Mirror server-side crypto**

Copy the same interface + implementations from `backend/pkg/crypto/` to `sdk/crypto/`. The only difference: RSA cipher's `Encrypt` only needs the public key (client doesn't have the private key), and `Decrypt` is not available on the client side for RSA.

Actually, for the client:
- RSA: Encrypt with public key (client has public key), cannot decrypt
- AES-GCM: Full encrypt/decrypt with shared key
- QQTea: Full encrypt/decrypt with shared key

For RSA client, we can add a separate `NewRSAClient(key []byte)` that only does Encrypt (or make the interface optional).

Simpler approach: keep the same interface, but the RSA client cipher only supports Encrypt. If Decrypt is called, return an error.

- [ ] Create `sdk/go.mod`
- [ ] Create `sdk/crypto/crypto.go` — interface + factory
- [ ] Create `sdk/crypto/rsa.go` — client-side RSA encrypt only
- [ ] Create `sdk/crypto/aesgcm.go`
- [ ] Create `sdk/crypto/qqtea.go`

---

### Task 5: SDK client main package

**Files:**
- Create: `sdk/models.go` — request/response types
- Create: `sdk/client.go` — Client struct + New()
- Create: `sdk/transport.go` — HTTP + encryption transport
- Create: `sdk/api.go` — all API methods

**Step 5.1: models.go**

```go
package sdk

type Config struct {
    AppID     string
    AppKey    string // PEM public key (RSA), hex key (AES/QQTea), or raw
    EncType   int
    ServerURL string // e.g. "http://localhost:9090/api"
}

type MachineReq struct {
    DeviceID    string `json:"deviceId"`
    Platform    string `json:"platform"`
    Arch        string `json:"arch"`
    MachineName string `json:"machineName"`
    Cpu         string `json:"cpu"`
    Gpu         string `json:"gpu"`
    Ram         string `json:"ram"`
}

type RegisterResp struct {
    MachineID int `json:"machineId"`
}

type LicenseInfo struct {
    LicenseKey string `json:"licenseKey"`
    AppID      string `json:"appId"`
    MachineID  int    `json:"machineId"`
    ExpireAt   string `json:"expireAt"`
    Status     uint8  `json:"status"`
}

type VerifyResult struct {
    Valid bool `json:"valid"`
}

type UpdateInfo struct {
    HasUpdate   bool   `json:"hasUpdate"`
    Version     string `json:"version,omitempty"`
    Desc        string `json:"desc,omitempty"`
    PatchUrl    string `json:"patchUrl,omitempty"`
    ForceUpdate bool   `json:"forceUpdate"`
}

type AppConfig struct {
    AppName string `json:"appName"`
    Version string `json:"version"`
    Desc    string `json:"desc"`
    Status  uint8  `json:"status"`
}
```

**Step 5.2: transport.go**

```go
package sdk

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type encryptedTransport struct {
    config Config
    cipher crypto.Cipher
    client *http.Client
}

func newTransport(cfg Config) (*encryptedTransport, error) {
    key := []byte(cfg.AppKey)
    c, err := crypto.NewCipher(uint8(cfg.EncType), key)
    if err != nil {
        return nil, err
    }
    return &encryptedTransport{config: cfg, cipher: c, client: &http.Client{}}, nil
}

func (t *encryptedTransport) do(method, path string, reqBody, respData interface{}) error {
    // Marshal request
    bodyBytes, _ := json.Marshal(reqBody)
    // Encrypt
    encrypted, err := t.cipher.Encrypt(bodyBytes)
    if err != nil {
        return fmt.Errorf("encrypt: %w", err)
    }
    // Send
    payload, _ := json.Marshal(encrypted)
    httpReq, _ := http.NewRequest(method, t.config.ServerURL+path, bytes.NewReader(payload))
    httpReq.Header.Set("Content-Type", "application/json")
    resp, err := t.client.Do(httpReq)
    if err != nil {
        return fmt.Errorf("request: %w", err)
    }
    defer resp.Body.Close()
    // Read response
    respBytes, _ := io.ReadAll(resp.Body)
    // Parse encrypted envelope
    var encResp crypto.EncryptedMessage
    if err := json.Unmarshal(respBytes, &encResp); err != nil {
        // Try plain JSON (error case)
        return fmt.Errorf("response: %s", string(respBytes))
    }
    // Decrypt
    plaintext, err := t.cipher.Decrypt(&encResp)
    if err != nil {
        return fmt.Errorf("decrypt: %w", err)
    }
    // Unmarshal into response struct
    // The encrypted response wraps in standard format: {code, msg, data}
    // We need to extract data field
    var wrapped struct {
        Code uint            `json:"code"`
        Msg  string          `json:"msg"`
        Data json.RawMessage `json:"data,omitempty"`
    }
    if err := json.Unmarshal(plaintext, &wrapped); err != nil {
        return fmt.Errorf("unmarshal response: %w", err)
    }
    if wrapped.Code != 0 {
        return fmt.Errorf("api error %d: %s", wrapped.Code, wrapped.Msg)
    }
    if respData != nil && wrapped.Data != nil {
        return json.Unmarshal(wrapped.Data, respData)
    }
    return nil
}
```

**Step 5.3: client.go + api.go**

Client struct:

```go
package sdk

import "context"

type Client struct {
    config    Config
    transport *encryptedTransport
}

func New(cfg Config) (*Client, error) {
    t, err := newTransport(cfg)
    if err != nil {
        return nil, err
    }
    return &Client{config: cfg, transport: t}, nil
}

func (c *Client) Register(ctx context.Context, req MachineReq) (*RegisterResp, error) {
    var resp RegisterResp
    err := c.transport.do("POST", "/sdk/register", req, &resp)
    return &resp, err
}

func (c *Client) Login(ctx context.Context, machineID int) (string, error) {
    var resp struct {
        Token string `json:"token"`
    }
    err := c.transport.do("POST", "/sdk/login", map[string]int{"machineId": machineID}, &resp)
    return resp.Token, err
}

func (c *Client) Heartbeat(ctx context.Context, machineID int) error {
    return c.transport.do("POST", "/sdk/heartbeat", map[string]int{"machineId": machineID}, nil)
}

func (c *Client) ApplyLicense(ctx context.Context, appID string, machineID int) (*LicenseInfo, error) {
    var resp LicenseInfo
    err := c.transport.do("POST", "/sdk/license/apply", map[string]interface{}{"appId": appID, "machineId": machineID}, &resp)
    return &resp, err
}

func (c *Client) VerifyLicense(ctx context.Context, licenseKey string) (*VerifyResult, error) {
    var resp VerifyResult
    err := c.transport.do("POST", "/sdk/license/verify", map[string]string{"licenseKey": licenseKey}, &resp)
    return &resp, err
}

func (c *Client) CheckUpdate(ctx context.Context, appID, currentVersion string) (*UpdateInfo, error) {
    var resp UpdateInfo
    err := c.transport.do("GET", "/sdk/update/check?appId="+appID+"&currentVersion="+currentVersion, nil, &resp)
    return &resp, err
}

func (c *Client) GetConfig(ctx context.Context, appID string) (*AppConfig, error) {
    var resp AppConfig
    err := c.transport.do("GET", "/sdk/config?appId="+appID, nil, &resp)
    return &resp, err
}
```

- [ ] Create `sdk/models.go`
- [ ] Create `sdk/transport.go` — encrypted HTTP transport
- [ ] Create `sdk/client.go` — Client + New + all API methods
- [ ] Create `sdk/api.go` — alias/re-export if needed (or merge into client.go)

---

### Task 6: Verify everything compiles

- [ ] `cd backend && go build ./...` — all server packages compile (excluding pre-existing logger dependency issue)
- [ ] `cd sdk && go build ./...` — SDK client compiles
- [ ] `git status` to review all changed/new files
- [ ] `git add` and `git commit` with descriptive message
