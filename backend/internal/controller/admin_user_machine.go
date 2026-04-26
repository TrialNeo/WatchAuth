package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

// Auth 给机器进行授权
func (u *UserMachineController) Auth(c *fiber.Ctx) error {
	var req struct {
		UserId   int    `json:"userId"`
		Appid    int    `json:"appId"`
		DeviceId string `json:"deviceId"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	if err := u.Service.Auth(req.Appid, req.DeviceId); err != nil {
		return re.withCode(errMsg.ERROR).Respond(nil)
	}
	return re.withCode(errMsg.SUCCESS).Respond(nil)
}

func (u *UserMachineController) List(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	code, data := u.Service.List()
	return re.withCode(code).Respond(data)
}

func (u *UserMachineController) Ban(c *fiber.Ctx) error {
	var req struct {
		MachineId int `json:"machineId"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	code := u.Service.Ban(req.MachineId)
	return re.withCode(code).Respond(nil)
}

// ReadLog 读取机器的上报日志
func (u *UserMachineController) ReadLog(c *fiber.Ctx) error {
	var req struct {
		MachineId int `json:"machineId"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	code, data := u.Service.ReadLog(req.MachineId, 50)
	return re.withCode(code).Respond(data)
}
