package route

import (
	"Diggpher/internal/controller"
	"Diggpher/internal/service"
	"Diggpher/pkg/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func bindUserRoute(api fiber.Router) {
	user := &controller.UserController{
		Svc: &service.UserService{},
	}
	// No auth required
	api.Post("/user/register", user.Register)
	api.Post("/user/login", user.Login)

	// Auth required
	api.Use("/user", auth.MiddlewareAuth()).Route("/user", func(router fiber.Router) {
		router.Get("/profile", user.Profile)
		router.Post("/license/purchase", user.PurchaseLicense)
		router.Get("/license/list", user.ListLicenses)
		router.Post("/license/revoke", user.RevokeLicense)
	})
}
