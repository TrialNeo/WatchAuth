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
