package controller

import (
	"Diggpher/internal/service/errMsg"
	"github.com/gofiber/fiber/v2"
)

type Response[T any] struct {
	Code    uint        `json:"code"`
	Message string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Respond 统一返回，汇总
func Respond(c *fiber.Ctx, code uint, data interface{}) error {
	if code != errMsg.SUCCESS {
		return c.Status(fiber.StatusOK).JSON(&Response[any]{
			Code:    code,
			Message: errMsg.GetErrMsg(code),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&Response[any]{
		Code:    code,
		Message: errMsg.GetErrMsg(code),
		Data:    data,
	})
}
