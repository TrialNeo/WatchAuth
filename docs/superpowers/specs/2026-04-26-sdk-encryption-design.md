# WatchAuth SDK + 加密功能设计

## 概述

为 WatchAuth 提供一套完整的 Go SDK，供第三方客户端应用集成，实现机器注册、登录认证、许可证管理、版本更新检查等功能。通信数据支持 RSA / AES-GCM / QQTea 三种加密方式。

## 加密层（EncType）

| 值 | 算法 | 类型 | 说明 |
|----|------|------|------|
| 0  | 无   | -    | 明文传输 |
| 1  | RSA  | 非对称 | 服务端持有私钥，客户端持有公钥 |
| 2  | AES-GCM | 对称 | 带认证加密，需要 nonce/iv |
| 3  | QQTea | 对称 | TEA 算法变种，分组密码 |

所有加密实现放在 `backend/pkg/crypto/`（服务端复用），`sdk/crypto/`（客户端独立包）。

服务端路由统一加解密中间件或由 controller 按需调用。

## 通信协议

每次 SDK 请求体使用加密信封：

```json
{
  "encType": 1,
  "ciphertext": "<base64 密文>",
  "nonce": "<base64 iv/nonce>"
}
```

服务端收到后根据 App 的 EncType 解密明文 `[]byte`，处理后再加密返回。

## SDK API 列表（服务端 + 客户端）

| 操作 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 注册 | POST | `/api/sdk/register` | 注册机器，返回 machineId + secret |
| 登录 | POST | `/api/sdk/login` | 机器认证，返回 session token |
| 心跳 | POST | `/api/sdk/heartbeat` | 维持在线状态 |
| 申请许可证 | POST | `/api/sdk/license/apply` | 为机器申请 app 许可证 |
| 验证许可证 | POST | `/api/sdk/license/verify` | 验证许可证有效性 |
| 检查更新 | GET | `/api/sdk/update/check` | 检查版本更新 |
| 获取配置 | GET | `/api/sdk/config` | 获取应用远程配置 |

## 项目结构变化

```
WatchAuth/
├── sdk/                          # 新建：Go SDK 客户端包
│   ├── go.mod                    # 独立 module，供外部导入
│   ├── client.go                 # Client 主结构 + New()
│   ├── models.go                 # 公共数据模型
│   ├── api.go                    # 全部 API 方法
│   ├── crypto/
│   │   ├── crypto.go             # Encrypter/Decrypter 接口
│   │   ├── rsa.go
│   │   ├── aesgcm.go
│   │   └── qqtea.go
│   └── transport.go              # HTTP 请求封装（加密/解密）
│
├── backend/
│   ├── pkg/crypto/               # 服务端加密实现
│   │   ├── crypto.go             # 统一的 Encrypt/Decrypt 接口
│   │   ├── rsa.go
│   │   ├── aesgcm.go
│   │   └── qqtea.go
│   └── internal/
│       ├── route/sdk.go          # 扩展 SDK 路由
│       ├── controller/sdk.go     # SDK 控制器
│       ├── controller/sdk_machine.go
│       ├── service/sdk.go        # SDK 业务逻辑
│       ├── service/sdk_machine.go # 已存在，重构
│       └── dao/                  # 需要新增 License 等模型
```

## SDK 客户端使用示例

```go
client := sdk.New(sdk.Config{
    AppID:     "your-app-id",
    AppKey:    "rsa-public-key-or-aes-key-or-tea-key",
    EncType:   sdk.RSA,
    ServerURL: "http://localhost:9090/api",
})

// 注册机器
machine, err := client.Register(ctx, sdk.MachineReq{
    DeviceID:    "unique-machine-id",
    Platform:    "windows",
    Arch:        "x64",
    MachineName: "My PC",
})

// 登录
session, err := client.Login(ctx, machine.MachineID)

// 心跳
err := client.Heartbeat(ctx, session.Token)

// 申请许可证
license, err := client.ApplyLicense(ctx, session.Token, appID)

// 验证许可证
valid, err := client.VerifyLicense(ctx, session.Token, licenseID)

// 检查版本更新
update, err := client.CheckUpdate(ctx, appID, currentVersion)

// 获取配置
config, err := client.GetConfig(ctx, session.Token)
```

## 数据库模型变更

```go
// 许可证模型（新增）
type License struct {
    ID        uint      `gorm:"primaryKey"`
    AppID     string    `gorm:"not null;index"`
    MachineID int       `gorm:"not null;index"`
    LicenseKey string   `gorm:"type:text;not null"`
    ExpireAt  time.Time
    IssuedAt  time.Time
    Status    uint8     // 0=有效 1=已吊销 2=已过期
}
```

完整 License 生命周期管理包括签发、验证、吊销、过期检测。

## 实现顺序

1. `backend/pkg/crypto/` — RSA、AES-GCM、QQTea 三种加密实现
2. `backend/internal/` — 扩展 SDK 路由 / controller / service（login 增强，新增 register/heartbeat/license/update/config）
3. `backend/internal/dao/` — License 模型 + AutoMigrate
4. `sdk/` — 独立 Go SDK 客户端（crypto + client + api）
