package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
	"Diggpher/pkg/middleware/auth"
	"crypto/rand"
	"encoding/hex"
	"time"
)

type SDKService struct{}

// --- Machine ---

// RegisterMachine creates a new machine info record and returns the machine ID.
func (s *SDKService) RegisterMachine(platform, arch, deviceID, machineName, cpu, gpu, ram string, belong int) (int, uint) {
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
	return info.MachineId, errMsg.SUCCESS
}

// GetAppConfig returns the app by appID (used by controllers to get SecretKeys for encryption).
func (s *SDKService) GetAppConfig(appID string) (*dao.App, uint) {
	var app dao.App
	if err := global.DataBase.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return nil, errMsg.ErrorAdminAppNotFound
	}
	return &app, errMsg.SUCCESS
}

// Heartbeat updates machine's last heartbeat time.
func (s *SDKService) Heartbeat(machineID int) uint {
	now := time.Now()
	result := global.DataBase.Model(&dao.UsedApp{}).
		Where("machine_id = ?", machineID).
		Update("last_heartbeat_at", now)
	if result.Error != nil {
		return errMsg.ERRORDataBaseErr
	}
	if result.RowsAffected == 0 {
		return errMsg.ErrorMachineNotFound
	}
	return errMsg.SUCCESS
}

// LoginResp contains the result of a machine login.
type LoginResp struct {
	Token     string `json:"token"`
	MachineID int    `json:"machineId"`
}

// Login authenticates a machine by deviceID and returns a JWT token.
func (s *SDKService) Login(deviceID string) (*LoginResp, uint) {
	var info dao.MachineInfo
	if err := global.DataBase.Where("device_id = ?", deviceID).First(&info).Error; err != nil {
		return nil, errMsg.ErrorMachineNotFound
	}
	var machine dao.Machine
	if err := global.DataBase.First(&machine, info.MachineId).Error; err != nil {
		return nil, errMsg.ErrorMachineNotFound
	}
	if machine.IsBan {
		return nil, errMsg.ErrorTokenInvalid
	}
	token, err := auth.GenerateToken(uint(info.MachineId))
	if err != nil {
		return nil, errMsg.ErrorAdminJWT
	}
	return &LoginResp{Token: token, MachineID: info.MachineId}, errMsg.SUCCESS
}

func generateLicenseKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ApplyLicense issues a new license for an app+machine.
func (s *SDKService) ApplyLicense(appID string, machineID int) (*dao.License, uint) {
	var app dao.App
	if err := global.DataBase.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return nil, errMsg.ErrorAdminAppNotFound
	}
	license := dao.License{
		LicenseKey: generateLicenseKey(),
		AppID:      appID,
		MachineID:  machineID,
		ExpireAt:   time.Now().AddDate(1, 0, 0),
		IssuedAt:   time.Now(),
		Status:     0,
	}
	if err := global.DataBase.Create(&license).Error; err != nil {
		return nil, errMsg.ERRORDataBaseErr
	}
	return &license, errMsg.SUCCESS
}

// LicenseVerifyResult contains the result of a license verification check.
type LicenseVerifyResult struct {
	Valid bool `json:"valid"`
}

// VerifyLicense checks if a license key is valid.
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

// UpdateCheckResult contains the result of a version update check.
type UpdateCheckResult struct {
	HasUpdate   bool   `json:"hasUpdate"`
	Version     string `json:"version,omitempty"`
	Desc        string `json:"desc,omitempty"`
	PatchUrl    string `json:"patchUrl,omitempty"`
	ForceUpdate bool   `json:"forceUpdate"`
}

// CheckUpdate checks if a newer version exists for the given app.
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
