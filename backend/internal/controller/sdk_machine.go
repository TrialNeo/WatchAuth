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
		return Fail(c, errMsg.ERROR, "")
	}
	service.Login(reqParam.Appid, &reqParam.Machine, c.IP())
	return Success(c, nil)
}
