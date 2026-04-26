package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
	"time"
)

type (
	OverviewStats struct {
		TotalMachines  int64 `json:"totalMachines"`
		TotalApps      int64 `json:"totalApps"`
		OnlineMachines int64 `json:"onlineMachines"`
	}

	MachineStats struct {
		TodayNew int64 `json:"todayNew"`
		MonthNew int64 `json:"monthNew"`
		Active7d int64 `json:"active7d"`
		Active30d int64 `json:"active30d"`
	}

	DistributionItem struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}

	Distributions struct {
		Platform []DistributionItem `json:"platform"`
		Arch     []DistributionItem `json:"arch"`
	}

	LicenseStats struct {
		Valid    int64 `json:"valid"`
		Revoked  int64 `json:"revoked"`
		Expired  int64 `json:"expired"`
		TodayNew int64 `json:"todayNew"`
		MonthNew int64 `json:"monthNew"`
	}

	AppRankingItem struct {
		AppName  string `json:"appName"`
		AppID    string `json:"appId"`
		Machines int64  `json:"machines"`
		Licenses int64  `json:"licenses"`
		Online   int64  `json:"online"`
	}

	DailyTrendItem struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	StatisticsResponse struct {
		Overview      OverviewStats    `json:"overview"`
		MachineStats  MachineStats     `json:"machineStats"`
		Distributions Distributions    `json:"distributions"`
		LicenseStats  LicenseStats     `json:"licenseStats"`
		AppRankings   []AppRankingItem `json:"appRankings"`
		MachineTrend  []DailyTrendItem `json:"machineTrend"`
		LicenseTrend  []DailyTrendItem `json:"licenseTrend"`
	}
)

func (*AdminService) Statistics() (*StatisticsResponse, uint) {
	resp := &StatisticsResponse{}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	sevenDaysAgo := now.AddDate(0, 0, -7)
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	// --- Overview ---
	global.DataBase.Model(&dao.Machine{}).Count(&resp.Overview.TotalMachines)
	global.DataBase.Model(&dao.App{}).Count(&resp.Overview.TotalApps)
	global.DataBase.Model(&dao.UsedApp{}).
		Where("last_heartbeat_at > ?", now.Add(-5*time.Minute)).
		Count(&resp.Overview.OnlineMachines)

	// --- Machine Stats ---
	global.DataBase.Model(&dao.MachineInfo{}).
		Where("created_at >= ?", todayStart).
		Count(&resp.MachineStats.TodayNew)
	global.DataBase.Model(&dao.MachineInfo{}).
		Where("created_at >= ?", monthStart).
		Count(&resp.MachineStats.MonthNew)

	global.DataBase.Model(&dao.UsedApp{}).
		Where("last_heartbeat_at > ?", sevenDaysAgo).
		Count(&resp.MachineStats.Active7d)
	global.DataBase.Model(&dao.UsedApp{}).
		Where("last_heartbeat_at > ?", thirtyDaysAgo).
		Count(&resp.MachineStats.Active30d)

	// --- Platform Distribution ---
	global.DataBase.Model(&dao.MachineInfo{}).
		Select("platform as name, count(*) as count").
		Group("platform").
		Find(&resp.Distributions.Platform)

	// --- Arch Distribution ---
	global.DataBase.Model(&dao.MachineInfo{}).
		Select("arch as name, count(*) as count").
		Group("arch").
		Find(&resp.Distributions.Arch)

	// --- License Stats ---
	global.DataBase.Model(&dao.License{}).
		Where("status = 0").
		Count(&resp.LicenseStats.Valid)
	global.DataBase.Model(&dao.License{}).
		Where("status = 1").
		Count(&resp.LicenseStats.Revoked)
	global.DataBase.Model(&dao.License{}).
		Where("status = 2").
		Count(&resp.LicenseStats.Expired)
	global.DataBase.Model(&dao.License{}).
		Where("issued_at >= ?", todayStart).
		Count(&resp.LicenseStats.TodayNew)
	global.DataBase.Model(&dao.License{}).
		Where("issued_at >= ?", monthStart).
		Count(&resp.LicenseStats.MonthNew)

	// --- App Rankings ---
	var apps []dao.App
	global.DataBase.Find(&apps)
	for _, app := range apps {
		item := AppRankingItem{
			AppName: app.AppName,
			AppID:   app.AppID,
		}
		global.DataBase.Model(&dao.UsedApp{}).
			Where("app_id = ?", app.ID).
			Count(&item.Machines)
		global.DataBase.Model(&dao.UsedApp{}).
			Where("app_id = ? AND last_heartbeat_at > ?", app.ID, now.Add(-5*time.Minute)).
			Count(&item.Online)
		global.DataBase.Model(&dao.License{}).
			Where("app_id = ?", app.AppID).
			Count(&item.Licenses)
		resp.AppRankings = append(resp.AppRankings, item)
	}

	// --- Machine Registration Trend (30 days) ---
	var rawMachineTrend []DailyTrendItem
	global.DataBase.Model(&dao.MachineInfo{}).
		Select("to_char(created_at, 'YYYY-MM-DD') as date, count(*) as count").
		Where("created_at >= ?", thirtyDaysAgo).
		Group("date").
		Order("date asc").
		Find(&rawMachineTrend)
	resp.MachineTrend = fillDailyTrend(rawMachineTrend, thirtyDaysAgo, now)

	// --- License Issuance Trend (30 days) ---
	var rawLicenseTrend []DailyTrendItem
	global.DataBase.Model(&dao.License{}).
		Select("to_char(issued_at, 'YYYY-MM-DD') as date, count(*) as count").
		Where("issued_at >= ?", thirtyDaysAgo).
		Group("date").
		Order("date asc").
		Find(&rawLicenseTrend)
	resp.LicenseTrend = fillDailyTrend(rawLicenseTrend, thirtyDaysAgo, now)

	return resp, errMsg.SUCCESS
}

// fillDailyTrend fills in zero-count days for complete 30-day series
func fillDailyTrend(raw []DailyTrendItem, start, end time.Time) []DailyTrendItem {
	lookup := make(map[string]int64, len(raw))
	for _, r := range raw {
		lookup[r.Date] = r.Count
	}
	var result []DailyTrendItem
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		c := lookup[dateStr]
		result = append(result, DailyTrendItem{Date: dateStr, Count: c})
	}
	return result
}
