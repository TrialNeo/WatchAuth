package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
)

type AdminLicenseItem struct {
	ID         uint   `json:"id"`
	LicenseKey string `json:"licenseKey"`
	AppID      string `json:"appId"`
	AppName    string `json:"appName"`
	MachineID  int    `json:"machineId"`
	UserID     uint   `json:"userId"`
	Username   string `json:"username"`
	ExpireAt   string `json:"expireAt"`
	IssuedAt   string `json:"issuedAt"`
	Status     uint8  `json:"status"`
}

func (*AdminService) LicenseList() ([]AdminLicenseItem, uint) {
	var licenses []dao.License
	global.DataBase.Order("id desc").Find(&licenses)

	appNames := make(map[string]string)
	usernames := make(map[uint]string)
	items := make([]AdminLicenseItem, 0, len(licenses))
	for _, l := range licenses {
		if _, ok := appNames[l.AppID]; !ok {
			var app dao.App
			if err := global.DataBase.Where("app_id = ?", l.AppID).First(&app).Error; err == nil {
				appNames[l.AppID] = app.AppName
			}
		}
		if l.UserID > 0 {
			if _, ok := usernames[l.UserID]; !ok {
				var user dao.User
				if err := global.DataBase.First(&user, l.UserID).Error; err == nil {
					usernames[l.UserID] = user.Username
				}
			}
		}
		items = append(items, AdminLicenseItem{
			ID:         l.ID,
			LicenseKey: l.LicenseKey,
			AppID:      l.AppID,
			AppName:    appNames[l.AppID],
			MachineID:  l.MachineID,
			UserID:     l.UserID,
			Username:   usernames[l.UserID],
			ExpireAt:   l.ExpireAt.Format("2006-01-02 15:04:05"),
			IssuedAt:   l.IssuedAt.Format("2006-01-02 15:04:05"),
			Status:     l.Status,
		})
	}
	return items, errMsg.SUCCESS
}
