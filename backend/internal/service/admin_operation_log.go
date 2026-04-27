package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
	"time"
)

type OperationLogItem struct {
	ID        uint   `json:"id"`
	AdminID   uint   `json:"adminId"`
	AdminName string `json:"adminName"`
	Action    string `json:"action"`
	Target    string `json:"target"`
	TargetID  string `json:"targetId"`
	Detail    string `json:"detail"`
	IP        string `json:"ip"`
	CreatedAt string `json:"createdAt"`
}

type OperationLogListResp struct {
	List  []OperationLogItem `json:"list"`
	Total int64              `json:"total"`
}

func (*AdminService) LogOperation(adminID uint, adminName, action, target, targetID, detail, ip string) {
	global.DataBase.Create(&dao.OperationLog{
		AdminID:   adminID,
		AdminName: adminName,
		Action:    action,
		Target:    target,
		TargetID:  targetID,
		Detail:    detail,
		IP:        ip,
	})
}

func (*AdminService) GetOperationLogs(page, pageSize int, action, target string) (*OperationLogListResp, uint) {
	var total int64
	query := global.DataBase.Model(&dao.OperationLog{})
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if target != "" {
		query = query.Where("target LIKE ?", "%"+target+"%")
	}
	query.Count(&total)

	var logs []dao.OperationLog
	query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	items := make([]OperationLogItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, OperationLogItem{
			ID:        l.ID,
			AdminID:   l.AdminID,
			AdminName: l.AdminName,
			Action:    l.Action,
			Target:    l.Target,
			TargetID:  l.TargetID,
			Detail:    l.Detail,
			IP:        l.IP,
			CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &OperationLogListResp{List: items, Total: total}, errMsg.SUCCESS
}

func (*AdminService) GetRecentLogs(limit int) ([]OperationLogItem, uint) {
	var logs []dao.OperationLog
	global.DataBase.Order("id desc").Limit(limit).Find(&logs)
	items := make([]OperationLogItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, OperationLogItem{
			ID:        l.ID,
			AdminID:   l.AdminID,
			AdminName: l.AdminName,
			Action:    l.Action,
			Target:    l.Target,
			TargetID:  l.TargetID,
			Detail:    l.Detail,
			IP:        l.IP,
			CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return items, errMsg.SUCCESS
}

// TodayStats returns counts for items created today
type TodayStats struct {
	NewMachines  int64 `json:"newMachines"`
	NewLicenses  int64 `json:"newLicenses"`
	NewUsers     int64 `json:"newUsers"`
	OnlineCount  int64 `json:"onlineCount"`
	TotalAgents  int64 `json:"totalAgents"`
}

func (*AdminService) GetTodayStats() (*TodayStats, uint) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	stats := &TodayStats{}

	global.DataBase.Model(&dao.MachineInfo{}).Where("created_at >= ?", todayStart).Count(&stats.NewMachines)
	global.DataBase.Model(&dao.License{}).Where("issued_at >= ?", todayStart).Count(&stats.NewLicenses)
	global.DataBase.Model(&dao.User{}).Where("created_at >= ?", todayStart).Count(&stats.NewUsers)
	global.DataBase.Model(&dao.UsedApp{}).Where("last_heartbeat_at > ?", now.Add(-5*time.Minute)).Count(&stats.OnlineCount)
	global.DataBase.Model(&dao.Agent{}).Count(&stats.TotalAgents)

	return stats, errMsg.SUCCESS
}
