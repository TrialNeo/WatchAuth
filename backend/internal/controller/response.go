package controller

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Code    uint        `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(&Response{
		Code:    0,
		Message: "",
		Data:    data,
	})
}

// Fail 业务失败（仍返回 HTTP 200）
func Fail(c *fiber.Ctx, code uint, message string) error {
	return c.Status(fiber.StatusOK).JSON(&Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
