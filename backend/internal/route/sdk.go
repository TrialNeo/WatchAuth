package route

import (
	"Diggpher/internal/controller"
	"Diggpher/internal/service"
	"github.com/gofiber/fiber/v2"
)

func bindSdkRoute(router fiber.Router) {
	sdk := &controller.SdkController{
		Svc: &service.SDKService{},
	}
	router.Post("/register", sdk.Register)
	router.Post("/login", sdk.Login)
	router.Post("/heartbeat", sdk.Heartbeat)
	router.Post("/license/apply", sdk.ApplyLicense)
	router.Post("/license/verify", sdk.VerifyLicense)
	router.Post("/update/check", sdk.CheckUpdate)
	router.Post("/config", sdk.GetConfig)
}
