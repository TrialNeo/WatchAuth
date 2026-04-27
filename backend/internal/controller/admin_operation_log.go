package controller

import (
	"github.com/gofiber/fiber/v2"
)

func (a *AdminController) GetOperationLogs(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)
	action := c.Query("action")
	target := c.Query("target")
	resp, code := a.Service.GetOperationLogs(page, pageSize, action, target)
	return re.withCode(code).Respond(resp)
}

func (a *AdminController) GetTodayStats(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	stats, code := a.Service.GetTodayStats()
	return re.withCode(code).Respond(stats)
}

func (a *AdminController) GetRecentLogs(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	logs, code := a.Service.GetRecentLogs(10)
	return re.withCode(code).Respond(fiber.Map{"list": logs})
}
