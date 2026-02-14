package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

// GetAppNameList 获取现有app名称列表
func (a *AdminController) GetAppNameList(c *fiber.Ctx) error {
	resp := a.Service.GetAppNameList()
	return Success(c, fiber.Map{
		"appNames": resp.AppNames,
	})
}

func (a *AdminController) GetVerList(c *fiber.Ctx) error {
	req := struct {
		Appid string `json:"appid"`
	}{}
	if err := c.BodyParser(&req); err != nil {
		req.Appid = ""
	}
	resp := a.Service.GetVerList(req.Appid)
	return Success(c, fiber.Map{"infos": resp.VerList})
}

func (a *AdminController) NewVer(c *fiber.Ctx) error {
	req := struct {
		Appid       string `json:"appid"`
		Version     string `json:"version"`
		Desc        string `json:"desc"`
		Sign        string `json:"sign"`
		ForceUpdate bool   `json:"forceUpdate"`
		Status      bool   `json:"status"`
		PatchUrl    string `json:"patchUrl"`
	}{}
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, errMsg.ERROR, err.Error())
	}
	resp := a.Service.NewVer(req.Appid, req.Version, req.Desc, req.Sign, req.PatchUrl, req.ForceUpdate, req.Status)
	if resp.Code != errMsg.SUCCESS {
		return Fail(c, resp.Code, resp.ErrMsg)
	}
	return Success(c, nil)
}
