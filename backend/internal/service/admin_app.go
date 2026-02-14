package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
)

type AppListResp struct {
	Code   uint   `json:"code"`
	ErrMsg string `json:"errMsg"`
	Apps   []*App `json:"apps"`
}

type App struct {
	AppName     string  `json:"appName"`
	AppID       string  `json:"appId"`
	Version     string  `json:"version"`
	FeeType     uint    `json:"feeType"`
	Fee         float64 `json:"fee"`
	Enctype     uint    `json:"enctype"`
	SecretKeys  string  `json:"secretKeys"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Description string  `json:"description"`
	Status      uint8   `json:"status"`
}

// GetAppList 获取应用列表
func (a *AdminService) GetAppList() (resp AppListResp) {
	var (
		appsDao = make([]*dao.App, 0)
		apps    = make([]*App, 0)
	)
	appsDao, _ = a.getAppsFromRedis()
	//	一个简单的对应转化
	for _, app := range appsDao {
		apps = append(apps, &App{
			AppID:       app.AppID,
			AppName:     app.AppName,
			Version:     app.Version,
			FeeType:     uint(app.FeeType),
			Fee:         app.Fee,
			Enctype:     uint(app.EncType),
			SecretKeys:  app.SecretKeys,
			CreatedAt:   app.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   app.UpdatedAt.Format("2006-01-02 15:04:05"),
			Description: app.Desc,
			Status:      app.Status,
		})
	}
	resp = AppListResp{
		Code:   errMsg.SUCCESS,
		ErrMsg: "",
		Apps:   apps,
	}
	return
}

// CreateApp 创建或者更新应用
func (a *AdminService) CreateApp(appid, appName, description string, EncType, feeType, status uint8, fee float64) (resp AppListResp) {
	//现在数据库中查找是否以及存在该App
	if errors.Is(global.DataBase.Where("app_name = ?", appName).First(new(dao.App)).Error, gorm.ErrRecordNotFound) {
		//	不存在，执行创建工作
		//根据算法生成随机secretKeys
		global.DataBase.
			Create(&dao.App{
				AppName:    appName,
				AppID:      strconv.Itoa(int(rand.Uint32())),
				Desc:       description,
				Status:     status,
				Version:    "1.0.0",
				FeeType:    feeType,
				Fee:        fee,
				EncType:    EncType,
				SecretKeys: "",
			})
	} else {
		global.DataBase.Where("app_id = ?", appid).
			Updates(&dao.App{
				Desc:    description,
				FeeType: feeType,
				Fee:     fee,
				EncType: EncType,
				Status:  status,
			})
	}
	a.updRSAppListFromDB()
	return
}

// DelApp 删除应用
func (a *AdminService) DelApp(appIDs []string) (resp AppListResp) {
	fmt.Println(appIDs)
	result := global.DataBase.Delete(&dao.App{}, "app_id IN (?)", appIDs)

	if len(appIDs) == 0 {
		resp.Code = errMsg.SUCCESS
		resp.ErrMsg = errMsg.GetErrMsg(errMsg.SUCCESS)
		return
	}

	if result.Error != nil {
		resp.Code = errMsg.ErrorAdminAppDelDBFail
		resp.ErrMsg = errMsg.GetErrMsg(errMsg.ErrorAdminAppDelDBFail)
		return
	}
	if result.RowsAffected == 0 {
		resp.Code = errMsg.ErrorAdminAppNotFound
		resp.ErrMsg = errMsg.GetErrMsg(errMsg.ErrorAdminAppNotFound)
		return
	}
	resp.Code = errMsg.SUCCESS
	resp.ErrMsg = errMsg.GetErrMsg(errMsg.SUCCESS)
	//刷新redis
	a.updRSAppListFromDB()
	return
}

type AppInfoResp struct {
	Code   uint   `json:"code"`
	ErrMsg string `json:"errMsg"`
	App    *App   `json:"app"`
}

// AppInfo 删除应用
func (a *AdminService) AppInfo(appid string) (resp AppInfoResp) {
	//直接从缓存获取，如果不走后台而是来抓包的话，其实也没问题，报错概率很小没必要加一个判断了
	var (
		appDao dao.App
	)
	isExists, _ := global.Redis.HExists(context.Background(), "cache:apps", appid).Result()
	if !isExists {
		resp.Code = errMsg.ErrorAdminAppNotFound
		resp.ErrMsg = errMsg.GetErrMsg(errMsg.ErrorAdminAppNotFound)
		return
	}
	appListJsonData, _ := global.Redis.HGet(context.Background(), "cache:apps", appid).Bytes()
	err := json.Unmarshal(appListJsonData, &appDao)
	if err != nil {
		resp.Code = errMsg.ERROR
		resp.ErrMsg = err.Error()
		return
	}
	//做一个查找 O(n)
	resp.Code = errMsg.SUCCESS
	resp.App = &App{
		AppName:     appDao.AppName,
		AppID:       appDao.AppID,
		Version:     appDao.Version,
		FeeType:     uint(appDao.FeeType),
		Fee:         appDao.Fee,
		Enctype:     uint(appDao.EncType),
		SecretKeys:  appDao.SecretKeys,
		CreatedAt:   appDao.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   appDao.UpdatedAt.Format("2006-01-02 15:04:05"),
		Description: appDao.Desc,
		Status:      appDao.Status,
	}
	return
}
