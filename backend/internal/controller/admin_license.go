package controller

import "github.com/gofiber/fiber/v2"

func (a *AdminController) LicenseList(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	licenses, code := a.Service.LicenseList()
	return re.withCode(code).Respond(licenses)
}
