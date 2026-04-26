package service

import (
	"Diggpher/global"
	"Diggpher/internal/dao"
	"Diggpher/internal/service/errMsg"
)

type AgentItem struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Contact   string  `json:"contact"`
	ParentID  uint    `json:"parentId"`
	Level     uint8   `json:"level"`
	Discount  float64 `json:"discount"`
	Balance   float64 `json:"balance"`
	Status    uint8   `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Children  []*AgentItem `json:"children,omitempty"`
}

func agentToItem(a dao.Agent) AgentItem {
	return AgentItem{
		ID:        a.ID,
		Name:      a.Name,
		Contact:   a.Contact,
		ParentID:  a.ParentID,
		Level:     a.Level,
		Discount:  a.Discount,
		Balance:   a.Balance,
		Status:    a.Status,
		CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (*AdminService) CreateAgent(name, contact string, parentID uint, discount float64) (uint, uint) {
	var parent dao.Agent
	if parentID > 0 {
		if err := global.DataBase.First(&parent, parentID).Error; err != nil {
			return 0, errMsg.ERRORAgentParentNotFound
		}
		if parent.Status != 0 {
			return 0, errMsg.ERRORAgentParentFrozen
		}
	}

	level := uint8(1)
	if parentID == 0 {
		level = 0
	} else {
		level = parent.Level + 1
		if level > 3 {
			return 0, errMsg.ERRORAgentLevelLimit
		}
	}

	agent := dao.Agent{
		Name:     name,
		Contact:  contact,
		ParentID: parentID,
		Level:    level,
		Discount: discount,
		Status:   0,
	}
	if err := global.DataBase.Create(&agent).Error; err != nil {
		return 0, errMsg.ERRORDataBaseErr
	}
	return agent.ID, errMsg.SUCCESS
}

func (*AdminService) UpdateAgent(id uint, name, contact string, discount float64, status uint8) uint {
	var agent dao.Agent
	if err := global.DataBase.First(&agent, id).Error; err != nil {
		return errMsg.ERRORAgentNotFound
	}
	updates := map[string]interface{}{
		"name":     name,
		"contact":  contact,
		"discount": discount,
		"status":   status,
	}
	if err := global.DataBase.Model(&agent).Updates(updates).Error; err != nil {
		return errMsg.ERRORDataBaseErr
	}
	return errMsg.SUCCESS
}

func (*AdminService) DeleteAgent(id uint) uint {
	var agent dao.Agent
	if err := global.DataBase.First(&agent, id).Error; err != nil {
		return errMsg.ERRORAgentNotFound
	}
	// Check for children
	var childCount int64
	global.DataBase.Model(&dao.Agent{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return errMsg.ERRORAgentHasChildren
	}
	if err := global.DataBase.Delete(&agent).Error; err != nil {
		return errMsg.ERRORDataBaseErr
	}
	return errMsg.SUCCESS
}

// buildAgentTree recursively builds the agent tree
func buildAgentTree(all []dao.Agent, parentID uint) []*AgentItem {
	var tree []*AgentItem
	for _, a := range all {
		if a.ParentID == parentID {
			item := agentToItem(a)
			children := buildAgentTree(all, a.ID)
			if len(children) > 0 {
				item.Children = children
			}
			tree = append(tree, &item)
		}
	}
	return tree
}

func (*AdminService) GetAgentList() ([]*AgentItem, uint) {
	var agents []dao.Agent
	global.DataBase.Order("id asc").Find(&agents)
	tree := buildAgentTree(agents, 0)
	return tree, errMsg.SUCCESS
}

func (*AdminService) GetAgent(id uint) (any, uint) {
	var agent dao.Agent
	if err := global.DataBase.First(&agent, id).Error; err != nil {
		return nil, errMsg.ERRORAgentNotFound
	}
	item := agentToItem(agent)

	// Get children
	var children []dao.Agent
	global.DataBase.Where("parent_id = ?", id).Order("id asc").Find(&children)
	for _, c := range children {
		child := agentToItem(c)
		item.Children = append(item.Children, &child)
	}
	return item, errMsg.SUCCESS
}
