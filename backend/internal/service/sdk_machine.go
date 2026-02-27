package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Machine struct {
	Platform    string `json:"platform"`
	Arch        string `json:"arch"`
	DeviceId    string `json:"deviceId"`
	MachineName string `json:"machineName"`
	Cpu         string `json:"cpu"`
	Gpu         string `json:"gpu"`
	Ram         string `json:"ram"`
}

func Login(appid int, machine *Machine, loginIp string) {
	var (
		queryMachineInfo = new(dao.MachineInfo)
	)
	res := global.DataBase.Where("device_id = ?", machine.DeviceId).Find(queryMachineInfo)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		//新用户，不可能出现的情况，肯定是已经绑好的。。
		return
	}
	//比对其他，不然就触发风险控制
	if queryMachineInfo.Platform == "" {
		//首次登录
		res = global.DataBase.Where("machine_id = ?", queryMachineInfo.MachineId).Updates(&dao.MachineInfo{
			MachineId:   queryMachineInfo.MachineId,
			Platform:    machine.Platform,
			Arch:        machine.Arch,
			DeviceId:    machine.DeviceId,
			MachineName: machine.MachineName,
			Cpu:         machine.Cpu,
			Gpu:         machine.Gpu,
			Ram:         machine.Ram,
		})
		now := time.Now()
		res = global.DataBase.Create(&dao.UsedApp{
			MachineId:       queryMachineInfo.MachineId,
			AppId:           appid,
			Online:          true,
			LoginIp:         loginIp,
			LastOnlineAt:    now,
			LastHeartbeatAt: &now,
			LastOfflineAt:   nil,
		})
		return
	}
	res = global.DataBase.Where("machine_id = ?", queryMachineInfo.MachineId).Updates(&dao.MachineInfo{})

}

//
//func Heartbeat() {
//
//}
