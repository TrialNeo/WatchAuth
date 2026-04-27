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
		userCtrl = &controller.AdminUserController{}
	)
	admin.Post("/admin/login", adminCtrl.Login)
	admin.Use(auth.MiddlewareAuth()).Route("/admin", func(router fiber.Router) {
		// 这是授权之后的使用hh
		router.Get("/permissions", adminCtrl.Permissions)
		router.Post("/statistics", adminCtrl.Statistics)
		router.Get("/license/list", adminCtrl.LicenseList)
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

		router.Route("/user", func(router fiber.Router) {

		})
		router.Route("/agent", func(router fiber.Router) {
			router.
				Get("/list", adminCtrl.GetAgentList).
				Get("/:id", adminCtrl.GetAgent).
				Post("/create", adminCtrl.CreateAgent).
				Post("/update", adminCtrl.UpdateAgent).
				Post("/delete", adminCtrl.DeleteAgent)
		})

		router.Route("/config", func(router fiber.Router) {
			router.
				Get("/list", adminCtrl.GetSystemConfigs).
				Post("/update", adminCtrl.UpdateSystemConfig)
		})

		router.Route("/log", func(router fiber.Router) {
			router.
				Get("/list", adminCtrl.GetOperationLogs).
				Get("/recent", adminCtrl.GetRecentLogs)
		})

		router.Route("/announcement", func(router fiber.Router) {
			router.
				Get("/list", adminCtrl.GetAnnouncementList).
				Get("/active", adminCtrl.GetActiveAnnouncements).
				Post("/create", adminCtrl.CreateAnnouncement).
				Post("/update", adminCtrl.UpdateAnnouncement).
				Post("/delete", adminCtrl.DeleteAnnouncement)
		})

		router.Get("/stats/today", adminCtrl.GetTodayStats)

		router.Route("/machine", func(router fiber.Router) {
			router.
				Post("/auth", userCtrl.Auth).
				Get("/list", userCtrl.List).
				Post("/ban", userCtrl.Ban).
				Post("/readLog", userCtrl.ReadLog)
		})
	})

}
