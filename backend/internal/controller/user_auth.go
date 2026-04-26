package controller

import (
	"Diggpher/internal/service"
	"Diggpher/internal/service/errMsg"
	"Diggpher/pkg/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Svc *service.UserService
}

func (u *UserController) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	resp, code := u.Svc.Register(req.Username, req.Password, req.Email)
	return re.withCode(code).Respond(resp)
}

func (u *UserController) Login(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	resp, code := u.Svc.Login(req.Username, req.Password)
	return re.withCode(code).Respond(resp)
}

func (u *UserController) Profile(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	userID, ok := auth.GetUserIDFromContext(c)
	if !ok {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	user, code := u.Svc.GetProfile(userID)
	return re.withCode(code).Respond(user)
}

func (u *UserController) PurchaseLicense(c *fiber.Ctx) error {
	var req struct {
		AppID string `json:"appId"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	userID, ok := auth.GetUserIDFromContext(c)
	if !ok {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	license, code := u.Svc.PurchaseLicense(userID, req.AppID)
	return re.withCode(code).Respond(license)
}

func (u *UserController) ListLicenses(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	userID, ok := auth.GetUserIDFromContext(c)
	if !ok {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	licenses, code := u.Svc.ListLicenses(userID)
	return re.withCode(code).Respond(licenses)
}

func (u *UserController) RevokeLicense(c *fiber.Ctx) error {
	var req struct {
		LicenseID uint `json:"licenseId"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	userID, ok := auth.GetUserIDFromContext(c)
	if !ok {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	code := u.Svc.RevokeLicense(userID, req.LicenseID)
	return re.withCode(code).Respond(nil)
}
