package controller

import (
	"Diggpher/internal/service"
	"Diggpher/internal/service/errMsg"
	"Diggpher/pkg/middleware/auth"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (a *AdminController) GetAppList(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	resp := a.Service.GetAppList()
	type Data = struct {
		Apps []*service.App `json:"apps,omitempty"`
	}
	return re.withCode(resp.Code).Respond(&Data{Apps: resp.Apps})
}

func (a *AdminController) CreateApp(c *fiber.Ctx) error {
	var req struct {
		AppId       string  `json:"appid"`
		AppName     string  `json:"appName"`
		Description string  `json:"description"`
		EncType     uint8   `json:"encType"`
		FeeType     uint8   `json:"feeType"`
		Fee         float64 `json:"fee"`
		Status      uint8   `json:"status"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	a.Service.CreateApp(req.AppId, req.AppName, req.Description, req.EncType, req.FeeType, req.Status, req.Fee)
	adminID, _ := auth.GetUserIDFromContext(c)
	a.Service.LogOperation(adminID, "", "create", "app", req.AppId,
		"创建应用: "+req.AppName, c.IP())
	return re.withCode(errMsg.SUCCESS).Respond(nil)
}

func (a *AdminController) DeleteApp(c *fiber.Ctx) error {
	var req struct {
		AppIDs []string `json:"appids"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	resp := a.Service.DelApp(req.AppIDs)
	if resp.Code == errMsg.SUCCESS {
		adminID, _ := auth.GetUserIDFromContext(c)
		a.Service.LogOperation(adminID, "", "delete", "app", "",
			fmt.Sprintf("删除应用: %v", req.AppIDs), c.IP())
	}
	return re.withCode(resp.Code).Respond(resp.Apps)
}

func (a *AdminController) AppInfo(c *fiber.Ctx) error {
	var req struct {
		AppID string `json:"appid"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	resp := a.Service.AppInfo(req.AppID)
	type Data = struct {
		App *service.App `json:"app,omitempty"`
	}
	return re.withCode(resp.Code).Respond(&Data{App: resp.App})
}
