package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

// Register 注册新机器
func (s *SdkController) Register(c *fiber.Ctx) error {
	var req struct {
		Platform    string `json:"platform"`
		Arch        string `json:"arch"`
		DeviceID    string `json:"deviceId"`
		MachineName string `json:"machineName"`
		Cpu         string `json:"cpu"`
		Gpu         string `json:"gpu"`
		Ram         string `json:"ram"`
	}
	re := newRespondIMP(c)
	ciph, code := s.decryptBody(c, &req)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	machineID, code := s.Svc.RegisterMachine(req.Platform, req.Arch, req.DeviceID, req.MachineName, req.Cpu, req.Gpu, req.Ram, 0)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	return encryptResp(re, ciph, code, map[string]int{"machineId": machineID})
}

// Login 机器登录认证
func (s *SdkController) Login(c *fiber.Ctx) error {
	var req struct {
		DeviceID string `json:"deviceId"`
	}
	re := newRespondIMP(c)
	ciph, code := s.decryptBody(c, &req)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	resp, code := s.Svc.Login(req.DeviceID)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	return encryptResp(re, ciph, code, resp)
}

// Heartbeat 心跳维持
func (s *SdkController) Heartbeat(c *fiber.Ctx) error {
	var req struct {
		MachineID int `json:"machineId"`
	}
	re := newRespondIMP(c)
	ciph, code := s.decryptBody(c, &req)
	if code != errMsg.SUCCESS {
		return re.withCode(code).Respond(nil)
	}
	code = s.Svc.Heartbeat(req.MachineID)
	return encryptResp(re, ciph, code, nil)
}
