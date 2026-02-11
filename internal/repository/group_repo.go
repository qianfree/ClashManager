package repository

import (
	"clash-manager/internal/model"
	"encoding/json"
)

// GroupWithNodes represents a proxy group with resolved node names
type GroupWithNodes struct {
	model.ProxyGroupModel
	ProxyNodes []model.ProxyNode `json:"ProxyNodes"`
}

type GroupRepository struct{}

func (r *GroupRepository) Create(group *model.ProxyGroupModel) error {
	return DB.Create(group).Error
}

func (r *GroupRepository) FindAll() ([]model.ProxyGroupModel, error) {
	var groups []model.ProxyGroupModel
	err := DB.Find(&groups).Error
	return groups, err
}

// FindAllWithNodes returns groups with resolved node information
func (r *GroupRepository) FindAllWithNodes() ([]GroupWithNodes, error) {
	var groups []model.ProxyGroupModel
	err := DB.Find(&groups).Error
	if err != nil {
		return nil, err
	}

	// Get all nodes for name resolution
	var nodes []model.Node
	err = DB.Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	// Build node ID to name map
	nodeMap := make(map[uint]string)
	for _, n := range nodes {
		nodeMap[n.ID] = n.Name
	}

	// Build result with resolved node names
	result := make([]GroupWithNodes, 0, len(groups))
	for _, g := range groups {
		gwn := GroupWithNodes{
			ProxyGroupModel: g,
			ProxyNodes:      []model.ProxyNode{},
		}

		// Parse ProxyIDs and resolve to names
		if g.ProxyIDs != "" {
			var nodeIDs []uint
			if err := json.Unmarshal([]byte(g.ProxyIDs), &nodeIDs); err == nil {
				for _, id := range nodeIDs {
					if name, ok := nodeMap[id]; ok {
						gwn.ProxyNodes = append(gwn.ProxyNodes, model.ProxyNode{
							ID:   id,
							Type: "node",
							Name: name,
						})
					}
				}
			}
		}

		result = append(result, gwn)
	}

	return result, nil
}

// Delete 软删除策略组，同时删除引用该策略组的规则
func (r *GroupRepository) Delete(id uint) error {
	// 开启事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 删除引用该策略组的规则（软删除）
	if err := tx.Where("target_id = ? AND target_type = ?", id, "group").Delete(&model.Rule{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 软删除策略组
	if err := tx.Delete(&model.ProxyGroupModel{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *GroupRepository) Update(group *model.ProxyGroupModel) error {
	return DB.Save(group).Error
}

func (r *GroupRepository) FindByID(id uint) (*model.ProxyGroupModel, error) {
	var group model.ProxyGroupModel
	err := DB.First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}
