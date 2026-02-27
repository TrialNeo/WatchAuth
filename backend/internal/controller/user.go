package controller

import "Diggpher/internal/service"

type AdminUserController struct {
	UserMachineController
}

type UserMachineController struct {
	Service service.AdminUserMachineService
}
