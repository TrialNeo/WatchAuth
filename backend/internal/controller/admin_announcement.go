package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

func (a *AdminController) CreateAnnouncement(c *fiber.Ctx) error {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  uint8  `json:"status"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	id, code := a.Service.CreateAnnouncement(req.Title, req.Content, req.Status)
	return re.withCode(code).Respond(fiber.Map{"id": id})
}

func (a *AdminController) UpdateAnnouncement(c *fiber.Ctx) error {
	var req struct {
		ID      uint   `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  uint8  `json:"status"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	code := a.Service.UpdateAnnouncement(req.ID, req.Title, req.Content, req.Status)
	return re.withCode(code).Respond(nil)
}

func (a *AdminController) DeleteAnnouncement(c *fiber.Ctx) error {
	var req struct {
		ID uint `json:"id"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	code := a.Service.DeleteAnnouncement(req.ID)
	return re.withCode(code).Respond(nil)
}

func (a *AdminController) GetAnnouncementList(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	list, code := a.Service.GetAnnouncementList()
	return re.withCode(code).Respond(fiber.Map{"list": list})
}

func (a *AdminController) GetActiveAnnouncements(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	list, code := a.Service.GetActiveAnnouncements()
	return re.withCode(code).Respond(fiber.Map{"list": list})
}
