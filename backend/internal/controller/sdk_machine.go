package controller

import (
	"Diggpher/internal/service"
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

func (s *SdkController) Login(c *fiber.Ctx) error {
	reqParam := struct {
		Appid   int             `json:"appid"`
		Machine service.Machine `json:"machine"`
	}{}
	if err := c.BodyParser(&reqParam); err != nil {
		return Respond(c, errMsg.ERRORInvalidParams, nil)
	}
	service.Login(reqParam.Appid, &reqParam.Machine, c.IP())
	return Respond(c, errMsg.SUCCESS, nil)
}
