package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

// CheckUpdate 检查版本更新
func (s *SdkController) CheckUpdate(c *fiber.Ctx) error {
	var req struct {
		AppID          string `json:"appId"`
		CurrentVersion string `json:"currentVersion"`
	}
	re := newRespondIMP(c)
	ciph, code := s.decryptBody(c, &req)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	result, code := s.Svc.CheckUpdate(req.AppID, req.CurrentVersion)
	return encryptResp(re, ciph, code, result)
}

// GetConfig 获取应用配置
func (s *SdkController) GetConfig(c *fiber.Ctx) error {
	var req struct {
		AppID string `json:"appId"`
	}
	re := newRespondIMP(c)
	ciph, code := s.decryptBody(c, &req)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	app, code := s.Svc.GetAppConfig(req.AppID)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	type Data = struct {
		AppName string `json:"appName"`
		Version string `json:"version"`
		Desc    string `json:"desc"`
		Status  uint8  `json:"status"`
	}
	return encryptResp(re, ciph, code, &Data{
		AppName: app.AppName,
		Version: app.Version,
		Desc:    app.Desc,
		Status:  app.Status,
	})
}
