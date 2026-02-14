package route

import (
	"Diggpher/internal/controller"
	"Diggpher/internal/service"
	"Diggpher/pkg/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

// bindAdminRoute 绑定admin的路由组
func bindAdminRoute(admin fiber.Router) {
	var (
		adminCtrl = &controller.AdminController{
			Service: service.NewAdminService(),
		}
	)
	admin.Post("/admin/login", adminCtrl.Login)
	admin.Use(auth.MiddlewareAuth()).Route("/admin", func(router fiber.Router) {
		// 这是授权之后的使用hh
		router.Get("/permissions", adminCtrl.Permissions)
		router.Route("/app", func(router fiber.Router) {
			router.
				Get("/list", adminCtrl.GetAppList).
				Post("/create", adminCtrl.CreateApp).
				Post("/delete", adminCtrl.DeleteApp).
				Post("/info", adminCtrl.AppInfo).
				Route("/version", func(router fiber.Router) {
					router.
						Get("/appNameList", adminCtrl.GetAppNameList).
						Post("/list", adminCtrl.GetVerList).
						Post("/create", adminCtrl.NewVer)
				})
		})
	})

}
