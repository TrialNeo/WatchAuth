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
		return Fail(c, errMsg.ERROR, "")
	}
	if err := u.Service.Auth(reqParam.Appid, reqParam.DeviceId); err != nil {
		return Fail(c, errMsg.ERROR, err.Error())
	}
	return Success(c, nil)
}
