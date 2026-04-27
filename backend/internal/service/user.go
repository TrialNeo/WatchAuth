package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
	"Diggpher/pkg/middleware/auth"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

type (
	UserLoginResp struct {
		Token    string `json:"token"`
		UserID   uint   `json:"userId"`
		Username string `json:"username"`
	}

	UserLicenseItem struct {
		ID         uint   `json:"id"`
		LicenseKey string `json:"licenseKey"`
		AppID      string `json:"appId"`
		AppName    string `json:"appName"`
		MachineID  int    `json:"machineId"`
		ExpireAt   string `json:"expireAt"`
		IssuedAt   string `json:"issuedAt"`
		Status     uint8  `json:"status"`
	}
)

func (*UserService) Register(username, password, email string) (*UserLoginResp, uint) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errMsg.ERROR
	}
	user := dao.User{
		Username: username,
		Password: string(hash),
		Email:    email,
		Status:   0,
	}
	if err := global.DataBase.Create(&user).Error; err != nil {
		return nil, errMsg.ERRORDataBaseErr
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, errMsg.ERROR
	}
	return &UserLoginResp{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}, errMsg.SUCCESS
}

func (*UserService) Login(username, password string) (*UserLoginResp, uint) {
	var user dao.User
	if err := global.DataBase.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errMsg.ErrorAdminUserNotFound
		}
		return nil, errMsg.ERRORDataBaseErr
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errMsg.ErrorAdminPswError
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, errMsg.ERROR
	}
	return &UserLoginResp{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}, errMsg.SUCCESS
}

func (*UserService) PurchaseLicense(userID uint, appID string) (*UserLicenseItem, uint) {
	var app dao.App
	if err := global.DataBase.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return nil, errMsg.ErrorAdminAppNotFound
	}
	key := make([]byte, 16)
	rand.Read(key)
	license := dao.License{
		LicenseKey: hex.EncodeToString(key),
		AppID:      appID,
		UserID:     userID,
		ExpireAt:   time.Now().AddDate(1, 0, 0),
		IssuedAt:   time.Now(),
		Status:     0,
	}
	if err := global.DataBase.Create(&license).Error; err != nil {
		return nil, errMsg.ERRORDataBaseErr
	}
	return &UserLicenseItem{
		ID:         license.ID,
		LicenseKey: license.LicenseKey,
		AppID:      license.AppID,
		AppName:    app.AppName,
		ExpireAt:   license.ExpireAt.Format("2006-01-02 15:04:05"),
		IssuedAt:   license.IssuedAt.Format("2006-01-02 15:04:05"),
		Status:     license.Status,
	}, errMsg.SUCCESS
}

func (*UserService) ListLicenses(userID uint) ([]UserLicenseItem, uint) {
	var licenses []dao.License
	global.DataBase.Where("user_id = ?", userID).Order("id desc").Find(&licenses)

	appNames := make(map[string]string)
	items := make([]UserLicenseItem, 0, len(licenses))
	for _, l := range licenses {
		if _, ok := appNames[l.AppID]; !ok {
			var app dao.App
			if err := global.DataBase.Where("app_id = ?", l.AppID).First(&app).Error; err == nil {
				appNames[l.AppID] = app.AppName
			}
		}
		items = append(items, UserLicenseItem{
			ID:         l.ID,
			LicenseKey: l.LicenseKey,
			AppID:      l.AppID,
			AppName:    appNames[l.AppID],
			MachineID:  l.MachineID,
			ExpireAt:   l.ExpireAt.Format("2006-01-02 15:04:05"),
			IssuedAt:   l.IssuedAt.Format("2006-01-02 15:04:05"),
			Status:     l.Status,
		})
	}
	return items, errMsg.SUCCESS
}

func (*UserService) RevokeLicense(userID uint, licenseID uint) uint {
	result := global.DataBase.Model(&dao.License{}).
		Where("id = ? AND user_id = ?", licenseID, userID).
		Update("status", 1)
	if result.RowsAffected == 0 {
		return errMsg.ERRORInvalidParams
	}
	return errMsg.SUCCESS
}

func (*UserService) GetActiveAnnouncements() ([]AnnouncementItem, uint) {
	var list []dao.Announcement
	global.DataBase.Where("status = 0").Order("id desc").Find(&list)
	items := make([]AnnouncementItem, 0, len(list))
	for _, a := range list {
		items = append(items, announcementToItem(a))
	}
	return items, errMsg.SUCCESS
}

func (*UserService) GetProfile(userID uint) (*dao.User, uint) {
	var user dao.User
	if err := global.DataBase.First(&user, userID).Error; err != nil {
		return nil, errMsg.ErrorAdminUserNotFound
	}
	return &user, errMsg.SUCCESS
}
