package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

func (a *AdminController) GetSystemConfigs(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	configs, code := a.Service.GetSystemConfigs()
	return re.withCode(code).Respond(fiber.Map{"configs": configs})
}

func (a *AdminController) UpdateSystemConfig(c *fiber.Ctx) error {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	code := a.Service.UpdateSystemConfig(req.Key, req.Value)
	return re.withCode(code).Respond(nil)
}
