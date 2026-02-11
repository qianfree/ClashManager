package repository

import (
	"clash-manager/internal/model"
	"encoding/json"
)

type NodeRepository struct{}

func (r *NodeRepository) Create(node *model.Node) error {
	return DB.Create(node).Error
}

func (r *NodeRepository) FindAll() ([]model.Node, error) {
	var nodes []model.Node
	err := DB.Find(&nodes).Error
	return nodes, err
}

func (r *NodeRepository) FindByID(id uint) (*model.Node, error) {
	var node model.Node
	err := DB.First(&node, id).Error
	return &node, err
}

func (r *NodeRepository) Update(node *model.Node) error {
	return DB.Save(node).Error
}

// Delete 物理删除节点，同时清理关联的规则和策略组引用
func (r *NodeRepository) Delete(id uint) error {
	// 开启事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 删除引用该节点的规则（软删除）
	if err := tx.Where("target_id = ? AND target_type = ?", id, "node").Delete(&model.Rule{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 清理策略组中对该节点的引用
	var groups []model.ProxyGroupModel
	if err := tx.Find(&groups).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, group := range groups {
		if group.ProxyIDs != "" {
			// 解析JSON并移除被删除的节点ID
			updated := removeNodeFromProxyIDs(group.ProxyIDs, id)
			if updated != group.ProxyIDs {
				if err := tx.Model(&group).Update("proxy_ids", updated).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	// 3. 物理删除节点
	if err := tx.Unscoped().Delete(&model.Node{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// removeNodeFromProxyIDs 从ProxyIDs JSON字符串中移除指定节点ID
func removeNodeFromProxyIDs(proxyIDs string, nodeID uint) string {
	if proxyIDs == "" || proxyIDs == "[]" {
		return proxyIDs
	}

	var ids []uint
	if err := json.Unmarshal([]byte(proxyIDs), &ids); err != nil {
		// 解析失败，返回原值
		return proxyIDs
	}

	// 移除指定ID
	var newIDs []uint
	for _, id := range ids {
		if id != nodeID {
			newIDs = append(newIDs, id)
		}
	}

	// 转换回JSON字符串
	result, err := json.Marshal(newIDs)
	if err != nil {
		return proxyIDs
	}

	return string(result)
}
