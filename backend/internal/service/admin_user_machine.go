package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"gorm.io/gorm"
)

type AdminUserMachineService struct{}

// Auth 给机器授权
func (u *AdminUserMachineService) Auth(userId int, deviceId string) (err error) {
	db := global.DataBase.Model(new(dao.Machine)).Where("device_id = ?", deviceId)
	if db.RowsAffected > 0 {
		return gorm.ErrDuplicatedKey
	}
	db = global.DataBase.Create(&dao.Machine{
		Belong: userId,
		MachineInfo: dao.MachineInfo{
			DeviceId: deviceId,
		},
	})
	if db.Error != nil {
		return err
	}
	return nil
}
