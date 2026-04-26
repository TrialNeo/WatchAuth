package sdk

import "context"

type Client struct {
	config    Config
	transport *encryptedTransport
}

func New(cfg Config) *Client {
	return &Client{
		config:    cfg,
		transport: newTransport(cfg),
	}
}

func (c *Client) Register(ctx context.Context, req MachineReq) (*RegisterResp, error) {
	var resp RegisterResp
	err := c.transport.do("POST", "/sdk/register", req, &resp)
	return &resp, err
}

func (c *Client) Login(ctx context.Context, deviceID string) (*LoginResp, error) {
	var resp LoginResp
	err := c.transport.do("POST", "/sdk/login", map[string]string{"deviceId": deviceID}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
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
	err := c.transport.do("POST", "/sdk/update/check", map[string]string{"appId": appID, "currentVersion": currentVersion}, &resp)
	return &resp, err
}

func (c *Client) GetConfig(ctx context.Context, appID string) (*AppConfig, error) {
	var resp AppConfig
	err := c.transport.do("POST", "/sdk/config", map[string]string{"appId": appID}, &resp)
	return &resp, err
}
