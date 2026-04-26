package sdk

type Config struct {
	AppID     string
	AppKey    string // PEM public key (RSA), raw bytes (AES/QQTea)
	EncType   uint8
	ServerURL string // e.g. "http://localhost:9090/api/sdk"
}

type MachineReq struct {
	Platform    string `json:"platform"`
	Arch        string `json:"arch"`
	DeviceID    string `json:"deviceId"`
	MachineName string `json:"machineName"`
	Cpu         string `json:"cpu"`
	Gpu         string `json:"gpu"`
	Ram         string `json:"ram"`
}

type RegisterResp struct {
	MachineID int `json:"machineId"`
}

type LoginResp struct {
	Token     string `json:"token"`
	MachineID int    `json:"machineId"`
}

type LicenseInfo struct {
	LicenseKey string `json:"licenseKey"`
	AppID      string `json:"appId"`
	MachineID  int    `json:"machineId"`
	ExpireAt   string `json:"expireAt"`
	IssuedAt   string `json:"issuedAt"`
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
