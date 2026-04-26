package controller

import (
	"Diggpher/internal/service"
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

// GetAppNameList 获取现有app名称列表
func (a *AdminController) GetAppNameList(c *fiber.Ctx) error {
	re := newRespondIMP(c)
	resp := a.Service.GetAppNameList()
	type Data = struct {
		AppNames []service.IKAppIDVAppName `json:"appNames,omitempty"`
	}
	return re.withCode(resp.Code).Respond(&Data{AppNames: resp.AppNames})
}

func (a *AdminController) GetVerList(c *fiber.Ctx) error {
	var req struct {
		Appid string `json:"appid"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		req.Appid = ""
	}
	resp := a.Service.GetVerList(req.Appid)
	type Data = struct {
		Infos []*service.VerInfo `json:"infos,omitempty"`
	}
	return re.withCode(resp.Code).Respond(&Data{Infos: resp.VerList})
}

func (a *AdminController) NewVer(c *fiber.Ctx) error {
	var req struct {
		Appid       string `json:"appid"`
		Version     string `json:"version"`
		Desc        string `json:"desc"`
		Sign        string `json:"sign"`
		ForceUpdate bool   `json:"forceUpdate"`
		Status      bool   `json:"status"`
		PatchUrl    string `json:"patchUrl"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	resp := a.Service.NewVer(req.Appid, req.Version, req.Desc, req.Sign, req.PatchUrl, req.ForceUpdate, req.Status)
	return re.withCode(resp.Code).Respond(nil)
}
