package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

// Auth 给机器进行授权
func (u *UserMachineController) Auth(c *fiber.Ctx) error {
	reqParam := struct {
		UserId   int    `json:"userId"`
		Appid    int    `json:"appId"`
		DeviceId string `json:"deviceId"`
	}{}
	if err := c.BodyParser(&reqParam); err != nil {
		return Respond(c, errMsg.ERRORInvalidParams, nil)
	}
	if err := u.Service.Auth(reqParam.Appid, reqParam.DeviceId); err != nil {
		return Respond(c, errMsg.ERROR, nil)
	}
	return Respond(c, errMsg.SUCCESS, nil)
}

func (u *UserMachineController) List(c *fiber.Ctx) error {
	code, response := u.Service.List()
	return Respond(c, code, response)
}

func (u *UserMachineController) Ban(c *fiber.Ctx) error {
	reqParam := struct {
		MachineId int `json:"machineId"`
	}{}
	if err := c.BodyParser(&reqParam); err != nil {
		return Respond(c, errMsg.ERRORInvalidParams, nil)
	}
	code := u.Service.Ban(reqParam.MachineId)
	return Respond(c, code, nil)
}

// ReadLog 读取机器的上报日志
func (u *UserMachineController) ReadLog(c *fiber.Ctx) error {
	reqParam := struct {
		MachineId int `json:"machineId"`
	}{}
	if err := c.BodyParser(&reqParam); err != nil {
		return Respond(c, errMsg.ERRORInvalidParams, nil)
	}
	code := u.Service.ReadLog(reqParam.MachineId)
	return Respond(c, code, nil)
}
