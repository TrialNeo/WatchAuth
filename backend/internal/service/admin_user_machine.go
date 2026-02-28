package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
	"errors"
	"gorm.io/gorm"
	"strconv"
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
	MachineId int  `json:"machineId"`
	Belong    int  `json:"belong"`
	IsBan     bool `json:"isBan"`
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
	response = make([]*MachineItem, 0, len(machinesDao))
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
			IsBan:     machineDao.IsBan,
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

func (u *AdminUserMachineService) Ban(machineId int) (code uint) {
	//先检查一下是否存在再说
	var (
		machineDao dao.Machine
	)
	res := global.DataBase.Where("machine_id = ?", machineId).First(&machineDao)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		code = errMsg.ErrorAdminAppVerUsed
		return
	}
	if res.Error != nil {
		code = errMsg.ERRORDataBaseErr
		return
	}
	res = global.DataBase.Model(machineDao).Update("is_ban", !machineDao.IsBan)
	if res.Error != nil {
		code = errMsg.ERRORDataBaseErr
		return
	}
	return
}

type MachineReportItem struct {
	Time         time.Time `json:"time"`
	Type         string    `json:"type"`
	AppId        uint      `json:"appId"`
	AppVersionId string    `json:"appVersionId"`
	TrainId      uint      `json:"trainId"`
	SpanId       uint      `json:"spanId"`
	Module       string    `json:"module"`
	FuncName     string    `json:"funcName"`
	Msg          string    `json:"msg"`
}

// ReadLog 读取上报日志
func (u *AdminUserMachineService) ReadLog(machineId int, num int) (code uint, reports []*MachineReportItem) {
	//先检查一下是否存在再说
	var (
		machineDao    dao.Machine
		machineLogDao = make([]*dao.MachineLog, 0)
	)
	res := global.DataBase.Where("machine_id = ?", machineId).First(&machineDao)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		code = errMsg.ErrorAdminAppVerUsed
		return
	}
	if res.Error != nil {
		code = errMsg.ERRORDataBaseErr
		return
	}
	res = global.DataBase.Where("machine_id = ?", machineId).
		Limit(num).Order("time DESC").
		Find(&machineLogDao)
	if res.RowsAffected > 0 {
		reports = make([]*MachineReportItem, 0, res.RowsAffected)
		for _, report := range machineLogDao {
			reports = append(reports, &MachineReportItem{
				Time:         report.Time,
				Type:         report.Type,
				AppId:        report.AppId,
				AppVersionId: strconv.Itoa(int(report.AppVersionId)),
				TrainId:      report.TrainId,
				SpanId:       report.SpanId,
				Module:       report.Module,
				FuncName:     report.FuncName,
				Msg:          report.Msg,
			})
		}
	}
	return
}
