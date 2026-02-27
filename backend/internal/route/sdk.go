package route

import (
	"Diggpher/internal/controller"
	"github.com/gofiber/fiber/v2"
)

// bindSdkRoute
func bindSdkRoute(router fiber.Router) {
	sdk := new(controller.SdkController)
	router.Post("/login", sdk.Login)
}
