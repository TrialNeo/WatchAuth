package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
)

type AnnouncementItem struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    uint8  `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func announcementToItem(a dao.Announcement) AnnouncementItem {
	return AnnouncementItem{
		ID:        a.ID,
		Title:     a.Title,
		Content:   a.Content,
		Status:    a.Status,
		CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (*AdminService) CreateAnnouncement(title, content string, status uint8) (uint, uint) {
	a := dao.Announcement{
		Title:   title,
		Content: content,
		Status:  status,
	}
	if err := global.DataBase.Create(&a).Error; err != nil {
		return 0, errMsg.ERRORDataBaseErr
	}
	return a.ID, errMsg.SUCCESS
}

func (*AdminService) UpdateAnnouncement(id uint, title, content string, status uint8) uint {
	var a dao.Announcement
	if err := global.DataBase.First(&a, id).Error; err != nil {
		return errMsg.ERRORAnnouncementNotFound
	}
	global.DataBase.Model(&a).Updates(map[string]interface{}{
		"title":   title,
		"content": content,
		"status":  status,
	})
	return errMsg.SUCCESS
}

func (*AdminService) DeleteAnnouncement(id uint) uint {
	var a dao.Announcement
	if err := global.DataBase.First(&a, id).Error; err != nil {
		return errMsg.ERRORAnnouncementNotFound
	}
	global.DataBase.Delete(&a)
	return errMsg.SUCCESS
}

func (*AdminService) GetAnnouncementList() ([]AnnouncementItem, uint) {
	var list []dao.Announcement
	global.DataBase.Order("id desc").Find(&list)
	items := make([]AnnouncementItem, 0, len(list))
	for _, a := range list {
		items = append(items, announcementToItem(a))
	}
	return items, errMsg.SUCCESS
}

// GetActiveAnnouncements returns published announcements (for user-facing endpoints)
func (*AdminService) GetActiveAnnouncements() ([]AnnouncementItem, uint) {
	var list []dao.Announcement
	global.DataBase.Where("status = 0").Order("id desc").Find(&list)
	items := make([]AnnouncementItem, 0, len(list))
	for _, a := range list {
		items = append(items, announcementToItem(a))
	}
	return items, errMsg.SUCCESS
}
