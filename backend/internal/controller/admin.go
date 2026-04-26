package controller

import (
	"Diggpher/internal/service"
	"Diggpher/internal/service/errMsg"
	"Diggpher/pkg/middleware/auth"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"os"
)

type AdminController struct {
	Service *service.AdminService
}

// Login 登录
func (a *AdminController) Login(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	re := newRespondIMP(c)
	if err := c.BodyParser(&req); err != nil {
		return re.withCode(errMsg.ERRORInvalidParams).Respond(nil)
	}
	resp := a.Service.Login(req.Username, req.Password, c.IP())
	return re.withCode(resp.Code).Respond(resp)
}

// Permissions 获取用户权限（通过用户身份获取菜单和按钮列表）
func (a *AdminController) Permissions(c *fiber.Ctx) error {
	type menu struct {
		Id         string      `json:"id"`
		Type       string      `json:"type"`
		Path       string      `json:"path"`
		Title      string      `json:"title"`
		Icon       string      `json:"icon"`
		ParentId   interface{} `json:"parentId"`
		Order      int         `json:"order"`
		Status     string      `json:"status"`
		CreateTime string      `json:"createTime"`
		UpdateTime string      `json:"updateTime"`
		IsBuiltIn  bool        `json:"isBuiltIn"`
		Children   []struct {
			Id         string `json:"id"`
			Type       string `json:"type"`
			Path       string `json:"path"`
			Title      string `json:"title"`
			Icon       string `json:"icon"`
			ParentId   string `json:"parentId"`
			Order      int    `json:"order"`
			Status     string `json:"status"`
			CreateTime string `json:"createTime"`
			UpdateTime string `json:"updateTime"`
			IsBuiltIn  bool   `json:"isBuiltIn"`
			Children   []struct {
				Id         string        `json:"id"`
				Type       string        `json:"type"`
				Path       string        `json:"path"`
				Title      string        `json:"title"`
				Icon       string        `json:"icon"`
				ParentId   string        `json:"parentId"`
				Order      int           `json:"order"`
				Status     string        `json:"status"`
				CreateTime string        `json:"createTime"`
				UpdateTime string        `json:"updateTime"`
				IsBuiltIn  bool          `json:"isBuiltIn"`
				Children   []interface{} `json:"children"`
			} `json:"children"`
		} `json:"children"`
	}

	re := newRespondIMP(c)
	_, _ = auth.GetUserIDFromContext(c)
	menus := make([]menu, 0)
	file, _ := os.ReadFile("./configs/menus.json")
	_ = json.Unmarshal(file, &menus)
	type Data = struct {
		ButtonPermissions []string `json:"buttonPermissions"`
		Menus             []menu   `json:"menus"`
	}
	return re.withCode(errMsg.SUCCESS).Respond(&Data{
		ButtonPermissions: []string{"user:add", "user:edit", "user:delete", "user:view", "role:add", "role:edit", "role:delete", "role:view", "menu:add", "menu:edit", "menu:delete", "menu:view"},
		Menus:             menus,
	})
}
