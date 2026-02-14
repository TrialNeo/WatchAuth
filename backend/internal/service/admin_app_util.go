package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"context"
	"encoding/json"
	"time"
)

// updRSAppListFromDB 从数据库中读取所有app信息，并且保存成一个appList格式的json，外部不可用
func (a *AdminService) updRSAppListFromDB() []dao.App {
	var (
		appList = make([]dao.App, 0)
		pipe    = global.Redis.Pipeline()
	)
	// 读取并且序列化
	global.DataBase.Model(new(dao.App)).Find(&appList)
	pipe.Del(context.Background(), "cache:apps")
	for _, app := range appList {
		data, _ := json.Marshal(app)
		pipe.HSet(context.Background(), "cache:apps", app.AppID, data)
	}
	pipe.Expire(context.Background(), "cache:apps", time.Hour*72)
	_, _ = pipe.Exec(context.Background())
	return appList
}

// getAppsFromRedis 从redis中获取apps
func (a *AdminService) getAppsFromRedis() ([]*dao.App, error) {
	data, err := global.Redis.HGetAll(context.Background(), "cache:apps").Result()
	if err != nil {
		return nil, err
	}
	var (
		apps = make([]*dao.App, 0)
	)
	for _, val := range data {
		appDao := new(dao.App)
		_ = json.Unmarshal([]byte(val), appDao)
		apps = append(apps, appDao)
	}
	return apps, nil
}

func (a *AdminService) updRSAVerList() []*dao.Version {
	verList := make([]*dao.Version, 0)
	global.DataBase.Model(new(dao.App)).Find(&verList)
	global.Redis.Del(context.Background(), "version")
	for i := range verList {
		global.Redis.Set(context.Background(), "version:"+verList[i].Appid, verList[i], time.Hour).Err()
	}
	return verList
}
