package controller

import (
	"github.com/gofiber/fiber/v2"
)

func (a *AdminController) Statistics(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	resp, code := a.Service.Statistics()
	return re.withCode(code).Respond(resp)
}
