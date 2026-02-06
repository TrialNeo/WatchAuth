package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"context"
	"encoding/json"
	"time"
)

// updRSAppListFromDB 从数据库中读取所有app信息，并且保存成一个appList格式的json，外部不可用
func (a *AdminService) updRSAppListFromDB() {
	appList := make([]dao.App, 0)
	global.DataBase.Model(new(dao.App)).Find(&appList)
	marshal, _ := json.Marshal(&appList)
	global.Redis.Set(context.Background(), "adminApps", marshal, time.Hour)
}
