package dao

import (
	"time"
)

// Machine 主表（包含 belong 和关联）
type Machine struct {
	MachineId   int          `json:"machineId" gorm:"primaryKey;column:machine_id;not null;unique"`
	Belong      int          `json:"belong" gorm:"column:belong;not null"` // 在顶层
	IsBan       bool         `json:"isBan" gorm:"column:is_ban;not null;default:false"`
	MachineInfo MachineInfo  `json:"machine" gorm:"foreignKey:MachineId;constraint:OnDelete:CASCADE"`
	MachineLogs []MachineLog `json:"machineLogs" gorm:"foreignKey:MachineId"`
	UsedApps    []UsedApp    `json:"usedApps" gorm:"foreignKey:MachineId;constraint:OnDelete:CASCADE"`
}

// MachineInfo 一对一详情表
type MachineInfo struct {
	MachineId   int    `gorm:"primaryKey;autoIncrement;column:machine_id"` // 主键 = 外键
	Platform    string `json:"platform" gorm:"column:platform"`
	Arch        string `json:"arch" gorm:"column:arch"`
	DeviceId    string `json:"deviceId" gorm:"column:device_id;uniqueIndex"`
	MachineName string `json:"machineName" gorm:"column:machine_name"`
	Cpu         string `json:"cpu" gorm:"column:cpu"`
	Gpu         string `json:"gpu" gorm:"column:gpu"`
	Ram         string `json:"ram" gorm:"column:ram"`
}

// UsedApp 一对多子表
type UsedApp struct {
	ID              int        `json:"-" gorm:"primaryKey;column:id"`
	MachineId       int        `json:"-" gorm:"column:machine_id;not null;index"`
	AppId           int        `json:"appId" gorm:"column:app_id;not null"`
	Online          bool       `json:"online" gorm:"column:online;not null"`
	LoginIp         string     `json:"loginIp,omitempty" gorm:"column:login_ip"`
	LastOnlineAt    time.Time  `json:"lastOnlineAt" gorm:"column:last_online_at;"`
	LastHeartbeatAt *time.Time `json:"lastHeartbeatAt" gorm:"column:last_heartbeat_at;"`
	LastOfflineAt   *time.Time `json:"lastOfflineAt" gorm:"column:last_offline_at"`
}

type MachineLog struct {
	ID           int       `json:"-" gorm:"autoIncrement;primaryKey;column:id"`
	MachineId    int       `json:"MachineId" gorm:"column:machine_id;not null"`
	Type         string    `json:"type" gorm:"column:type;not null"`
	AppId        uint      `json:"appId" gorm:"column:appId;not null"`
	AppVersionId uint      `json:"appVersionId" gorm:"column:appVersionId;not null"`
	TrainId      uint      `json:"trainId" gorm:"column:trainId;not null"`
	SpanId       uint      `json:"spanId" gorm:"column:spanId;not null"`
	Module       string    `json:"module" gorm:"column:module"`
	FuncName     string    `json:"funcName" gorm:"column:func_name"`
	Msg          string    `json:"msg" gorm:"column:msg;not null"`
	Time         time.Time `json:"time" gorm:"column:time;"`
}

func (Machine) TableName() string     { return "machines" }
func (MachineInfo) TableName() string { return "machine_infos" }
func (UsedApp) TableName() string     { return "used_apps" }
