package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
	"gorm.io/gorm"
	"time"
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

type UsedApps struct {
	AppId           int        `json:"appId"`
	Online          bool       `json:"online"`
	LoginIp         string     `json:"loginIp,omitempty"`
	LastOnlineAt    time.Time  `json:"lastOnlineAt"`
	LastHeartbeatAt *time.Time `json:"lastHeartbeatAt"`
	LastOfflineAt   *time.Time `json:"lastOfflineAt"`
}
type MachineItem struct {
	MachineId int `json:"machineId"`
	Belong    int `json:"belong"`
	Machine   struct {
		Platform    string `json:"platform"`
		Arch        string `json:"arch"`
		DeviceId    string `json:"deviceId"`
		MachineName string `json:"machineName"`
		Cpu         string `json:"cpu"`
		Gpu         string `json:"gpu"`
		Ram         string `json:"ram"`
	} `json:"machine"`
	UsedApps []UsedApps `json:"usedApps"`
}

// List 机器列表
func (u *AdminUserMachineService) List() (code uint, response []*MachineItem) {
	var (
		machinesDao = make([]*dao.Machine, 0)
	)
	res := global.DataBase.
		Preload("MachineInfo").
		Preload("UsedApps").
		Find(&machinesDao)
	if res.Error != nil {
		code = errMsg.ERRORDataBaseErr
		return
	}
	response = make([]*MachineItem, 0)
	for _, machineDao := range machinesDao {
		usedApps := make([]UsedApps, 0)
		for _, usedApp := range machineDao.UsedApps {
			usedApps = append(usedApps, UsedApps{
				AppId:           usedApp.AppId,
				Online:          usedApp.Online,
				LoginIp:         usedApp.LoginIp,
				LastOnlineAt:    usedApp.LastOnlineAt,
				LastHeartbeatAt: usedApp.LastHeartbeatAt,
				LastOfflineAt:   usedApp.LastOfflineAt,
			})
		}
		response = append(response, &MachineItem{
			MachineId: machineDao.MachineId,
			Belong:    machineDao.Belong,
			Machine: Machine{
				Platform:    machineDao.MachineInfo.Platform,
				Arch:        machineDao.MachineInfo.Arch,
				DeviceId:    machineDao.MachineInfo.DeviceId,
				MachineName: machineDao.MachineInfo.MachineName,
				Cpu:         machineDao.MachineInfo.Cpu,
				Gpu:         machineDao.MachineInfo.Gpu,
				Ram:         machineDao.MachineInfo.Ram,
			},
			UsedApps: usedApps,
		})
	}
	code = errMsg.SUCCESS
	return
}
