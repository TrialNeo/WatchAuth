package controller

import (
	"Diggpher/internal/service/errMsg"
	"Diggpher/pkg/middleware/auth"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (a *AdminController) CreateAgent(c *fiber.Ctx) error {
	var req struct {
		Name     string  `json:"name"`
		Contact  string  `json:"contact"`
		ParentID uint    `json:"parentId"`
		Discount float64 `json:"discount"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	id, code := a.Service.CreateAgent(req.Name, req.Contact, req.ParentID, req.Discount)
	if code == errMsg.SUCCESS {
		adminID, _ := auth.GetUserIDFromContext(c)
		a.Service.LogOperation(adminID, "", "create", "agent", fmt.Sprintf("%d", id),
			fmt.Sprintf("创建代理: %s (父级ID: %d)", req.Name, req.ParentID), c.IP())
	}
	return re.withCode(code).Respond(fiber.Map{"id": id})
}

func (a *AdminController) UpdateAgent(c *fiber.Ctx) error {
	var req struct {
		ID       uint    `json:"id"`
		Name     string  `json:"name"`
		Contact  string  `json:"contact"`
		Discount float64 `json:"discount"`
		Status   uint8   `json:"status"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	code := a.Service.UpdateAgent(req.ID, req.Name, req.Contact, req.Discount, req.Status)
	if code == errMsg.SUCCESS {
		adminID, _ := auth.GetUserIDFromContext(c)
		a.Service.LogOperation(adminID, "", "update", "agent", fmt.Sprintf("%d", req.ID),
			fmt.Sprintf("更新代理: %s (状态: %d)", req.Name, req.Status), c.IP())
	}
	return re.withCode(code).Respond(nil)
}

func (a *AdminController) DeleteAgent(c *fiber.Ctx) error {
	var req struct {
		ID uint `json:"id"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	code := a.Service.DeleteAgent(req.ID)
	if code == errMsg.SUCCESS {
		adminID, _ := auth.GetUserIDFromContext(c)
		a.Service.LogOperation(adminID, "", "delete", "agent", fmt.Sprintf("%d", req.ID),
			fmt.Sprintf("删除代理ID: %d", req.ID), c.IP())
	}
	return re.withCode(code).Respond(nil)
}

func (a *AdminController) GetAgentList(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	tree, code := a.Service.GetAgentList()
	return re.withCode(code).Respond(fiber.Map{"agents": tree})
}

func (a *AdminController) GetAgent(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	agent, code := a.Service.GetAgent(uint(id))
	return re.withCode(code).Respond(fiber.Map{"agent": agent})
}
