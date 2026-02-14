package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
	"time"
)

type AppNameListResp struct {
	Code     uint              `json:"code"`
	ErrMsg   string            `json:"errMsg"`
	AppNames []IKAppIDVAppName `json:"appNames"`
}

type IKAppIDVAppName struct {
	AppId   string `json:"appId"`
	AppName string `json:"appName"`
}

// GetAppNameList 获取应用列表
func (a *AdminService) GetAppNameList() (resp AppNameListResp) {
	//先尝试从redis中获取 提高性能
	appDao, err := a.getAppsFromRedis()
	if err != nil {
		resp.Code = errMsg.ERROR
		resp.ErrMsg = err.Error()
		return resp
	}
	appNames := make([]IKAppIDVAppName, 0)
	for _, app := range appDao {
		appNames = append(appNames, IKAppIDVAppName{
			AppId:   app.AppID,
			AppName: app.AppName,
		})
	}
	resp = AppNameListResp{
		Code:     errMsg.SUCCESS,
		ErrMsg:   "",
		AppNames: appNames,
	}
	return
}

type VerListResp struct {
	Code    uint       `json:"code"`
	ErrMsg  string     `json:"errMsg"`
	VerList []*VerInfo `json:"verList"`
}

type VerInfo struct {
	Appid       string    `json:"appid"`
	AppName     string    `json:"appName"`
	Version     string    `json:"version"`
	Desc        string    `json:"desc"`
	Sign        string    `json:"sign"`
	ForceUpdate bool      `json:"forceUpdate"`
	Status      bool      `json:"status"`
	PatchUrl    string    `json:"patch_url"`
	CreatedAt   time.Time `json:"updatedTime"`
}

// GetVerList 获取版本信息列表
func (a *AdminService) GetVerList(appid string) (resp VerListResp) {
	daoVerList := make([]*dao.Version, 0)
	if appid == "" {
		global.DataBase.Model(new(dao.Version)).Find(&daoVerList)
	} else {
		global.DataBase.Model(new(dao.Version)).Where("appid=?", appid).Find(&daoVerList)
	}
	//映射
	verList := make([]*VerInfo, 0)
	for _, ver := range daoVerList {
		verList = append(verList, &VerInfo{
			Appid:       ver.Appid,
			Version:     ver.Version,
			Desc:        ver.Desc,
			Sign:        ver.Sign,
			ForceUpdate: ver.ForceUpdate,
			Status:      ver.Status,
			PatchUrl:    ver.PatchUrl,
			CreatedAt:   ver.CreatedAt,
		})
	}
	resp.Code = errMsg.SUCCESS
	resp.VerList = verList
	return
}

type NewVerResp struct {
	Code   uint   `json:"code"`
	ErrMsg string `json:"errMsg"`
}

// NewVer 进行更新
func (a *AdminService) NewVer(Appid, Version, Desc, Sign, PatchUrl string, ForceUpdate, Status bool) (resp NewVerResp) {
	//检查应用是否存在
	if a.AppInfo(Appid).Code == errMsg.ErrorAdminAppNotFound {
		resp.Code = errMsg.ErrorAdminAppNotFound
		resp.ErrMsg = errMsg.GetErrMsg(errMsg.ErrorAdminAppNotFound)
		return
	}
	if global.DataBase.Model(new(dao.Version)).Where("appid = ? and version = ?", Appid, Version).First(nil).RowsAffected != 0 {
		resp.Code = errMsg.ErrorAdminAppVerUsed
		resp.ErrMsg = errMsg.GetErrMsg(errMsg.ErrorAdminAppVerUsed)
		return
	}
	if err := global.DataBase.Create(&dao.Version{
		Appid:    Appid,
		Version:  Version,
		Desc:     Desc,
		Sign:     Sign,
		PatchUrl: PatchUrl,
	}).Error; err != nil {
		resp.Code = errMsg.ERRORDataBaseErr
		resp.ErrMsg = errMsg.GetErrMsg(errMsg.ERRORDataBaseErr)
		return
	}
	resp.Code = errMsg.SUCCESS
	return
}
