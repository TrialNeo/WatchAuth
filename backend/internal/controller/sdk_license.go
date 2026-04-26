package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

// ApplyLicense 申请许可证
func (s *SdkController) ApplyLicense(c *fiber.Ctx) error {
	var req struct {
		AppID     string `json:"appId"`
		MachineID int    `json:"machineId"`
	}
	re := newRespondIMP(c)
	ciph, code := s.decryptBody(c, &req)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	license, code := s.Svc.ApplyLicense(req.AppID, req.MachineID)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	return encryptResp(re, ciph, code, license)
}

// VerifyLicense 验证许可证
func (s *SdkController) VerifyLicense(c *fiber.Ctx) error {
	var req struct {
		LicenseKey string `json:"licenseKey"`
	}
	re := newRespondIMP(c)
	ciph, code := s.decryptBody(c, &req)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	result, code := s.Svc.VerifyLicense(req.LicenseKey)
	return encryptResp(re, ciph, code, result)
}
