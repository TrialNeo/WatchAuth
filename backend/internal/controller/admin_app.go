package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

// GetAppList 获取现有app详情列表
func (a *AdminController) GetAppList(c *fiber.Ctx) error {
	resp := a.Service.GetAppList()
	return Respond(c, resp.Code, fiber.Map{
		"apps": resp.Apps,
	})
}

// CreateApp 更新app配置信息
func (a *AdminController) CreateApp(c *fiber.Ctx) error {
	reqParam := struct {
		AppId       string  `json:"appid"`
		AppName     string  `json:"appName"`
		Description string  `json:"description"`
		EncType     uint8   `json:"encType"`
		FeeType     uint8   `json:"feeType"`
		Fee         float64 `json:"fee"`
		Status      uint8   `json:"status"`
	}{}
	if err := c.BodyParser(&reqParam); err != nil {
		return Respond(c, errMsg.ERROR, "")
	}
	a.Service.CreateApp(reqParam.AppId, reqParam.AppName, reqParam.Description, reqParam.EncType, reqParam.FeeType, reqParam.Status, reqParam.Fee)
	return Respond(c, errMsg.SUCCESS, fiber.Map{})
}

// DeleteApp 删除选中的一些app
func (a *AdminController) DeleteApp(c *fiber.Ctx) error {
	reqParam := struct {
		AppIDs []string `json:"appids"`
	}{}
	if err := c.BodyParser(&reqParam); err != nil {
		return Respond(c, errMsg.ERRORInvalidParams, nil)
	}
	resp := a.Service.DelApp(reqParam.AppIDs)
	return Respond(c, resp.Code, resp.Apps)
}

// AppInfo 获取单个app详细信息
func (a *AdminController) AppInfo(c *fiber.Ctx) error {
	reqParam := struct {
		AppID string `json:"appid"`
	}{}
	if err := c.BodyParser(&reqParam); err != nil {
		return Respond(c, errMsg.ERRORInvalidParams, nil)
	}
	resp := a.Service.AppInfo(reqParam.AppID)
	return Respond(c, resp.Code, fiber.Map{
		"app": resp.App,
	})
}
