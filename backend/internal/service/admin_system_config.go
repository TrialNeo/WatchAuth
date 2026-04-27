package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
)

type ConfigItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// InitSystemConfigs seeds default config values
func InitSystemConfigs() {
	defaults := map[string]string{
		"site_name":             "WatchAuth",
		"site_logo":             "",
		"allow_register":        "true",
		"default_license_days":  "365",
		"license_check_interval": "300",
	}
	for k, v := range defaults {
		var count int64
		global.DataBase.Model(&dao.SystemConfig{}).Where("key = ?", k).Count(&count)
		if count == 0 {
			global.DataBase.Create(&dao.SystemConfig{Key: k, Value: v})
		}
	}
}

func (*AdminService) GetSystemConfigs() ([]ConfigItem, uint) {
	var configs []dao.SystemConfig
	global.DataBase.Find(&configs)
	items := make([]ConfigItem, 0, len(configs))
	for _, c := range configs {
		items = append(items, ConfigItem{Key: c.Key, Value: c.Value})
	}
	return items, errMsg.SUCCESS
}

func (*AdminService) UpdateSystemConfig(key, value string) uint {
	var cfg dao.SystemConfig
	if err := global.DataBase.Where("key = ?", key).First(&cfg).Error; err != nil {
		return errMsg.ERRORConfigNotFound
	}
	global.DataBase.Model(&cfg).Update("value", value)
	return errMsg.SUCCESS
}
