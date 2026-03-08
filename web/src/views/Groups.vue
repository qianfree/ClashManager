<template>
  <div class="groups-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>代理组列表</span>
          </div>
          <el-button type="primary" :icon="Plus" @click="showCreateDialog">新增代理组</el-button>
        </div>
      </template>

      <el-table :data="groups" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" min-width="60" />
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="type" label="类型" min-width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'select' ? 'primary' : 'success'">
              {{ row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" min-width="100">
          <template #default="{ row }">
            <el-tag :type="getSourceType(row.source)" size="small">
              {{ getSourceText(row.source) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="proxy_nodes" label="包含节点" min-width="280">
          <template #default="{ row }">
            <el-tag v-for="proxy in displayNodes(row).slice(0, 5)" :key="proxy.id || proxy" size="small" style="margin-right: 5px;">
              {{ proxy.name || proxy }}
            </el-tag>
            <span v-if="displayNodes(row).length > 5" style="color: #909399; font-size: 12px;">
              ...等{{ displayNodes(row).length }}个
            </span>
            <span v-if="displayNodes(row).length === 0" style="color: #909399;">无</span>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="测试URL" min-width="200" show-overflow-tooltip />
        <el-table-column prop="interval" label="间隔(秒)" min-width="100" />
        <el-table-column label="操作" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑代理组对话框 -->
    <el-dialog v-model="formDialogVisible" :title="isEdit ? '编辑代理组' : '新增代理组'" width="700px">
      <el-form :model="groupForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="groupForm.Name" placeholder="请输入代理组名称" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="groupForm.Type" placeholder="请选择类型" style="width: 100%;">
            <el-option label="手动选择" value="select" />
            <el-option label="自动测速" value="url-test" />
            <el-option label="负载均衡" value="load-balance" />
          </el-select>
        </el-form-item>
        <el-form-item label="包含节点">
          <el-transfer
            v-model="groupForm.ProxyIDs"
            :data="transferData"
            :props="{
              key: 'id',
              label: 'label'
            }"
            filterable
            filter-placeholder="搜索节点"
            :titles="['可选节点', '已选节点']"
            style="text-align: left; justify-content: flex-start;"
          />
          <div style="color: #909399; font-size: 12px; margin-top: 8px;">
            已选择 {{ groupForm.ProxyIDs.length }} / {{ nodes.length }} 个节点
          </div>
        </el-form-item>
        <el-form-item label="测试URL">
          <el-input v-model="groupForm.URL" placeholder="http://www.gstatic.com/generate_204" />
        </el-form-item>
        <el-form-item label="测试间隔(秒)">
          <el-input-number v-model="groupForm.Interval" :min="10" :max="3600" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getGroups, createGroup, updateGroup, deleteGroup } from '@/api/groups'
import { getNodes } from '@/api/nodes'

const groups = ref([])
const nodes = ref([])
const formDialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const groupForm = ref({
  Name: '',
  Type: 'select',
  ProxyIDs: [],
  URL: 'http://www.gstatic.com/generate_204',
  Interval: 300
})

// 穿梭框数据源
const transferData = computed(() => {
  return nodes.value.map(node => ({
    id: node.id,
    label: node.name,
    disabled: false
  }))
})

const loadGroups = async () => {
  groups.value = await getGroups()
}

// 兼容新旧数据格式显示节点
const displayNodes = (row) => {
  // 新格式：使用 proxy_nodes 数组
  if (row.proxy_nodes && row.proxy_nodes.length > 0) {
    return row.proxy_nodes
  }
  // 旧格式：解析 Proxies JSON 字符串
  if (row.proxies) {
    try {
      const proxies = JSON.parse(row.proxies)
      if (Array.isArray(proxies)) {
        return proxies.map(p => ({ id: null, name: p }))
      }
    } catch {
      return []
    }
  }
  return []
}

const getSourceText = (source) => {
  const map = {
    'local': '本地',
    'remote': '远程'
  }
  return map[source] || source
}

const getSourceType = (source) => {
  const map = {
    'local': 'info',
    'remote': 'warning'
  }
  return map[source] || ''
}

const loadNodes = async () => {
  nodes.value = await getNodes()
}

const showCreateDialog = async () => {
  isEdit.value = false
  editId.value = null
  groupForm.value = {
    Name: '',
    Type: 'select',
    ProxyIDs: [],
    URL: 'http://www.gstatic.com/generate_204',
    Interval: 300
  }
  await loadNodes()
  formDialogVisible.value = true
}

const handleEdit = async (row) => {
  isEdit.value = true
  editId.value = row.id
  // Parse proxy_ids from JSON if needed
  let proxyIDs = []
  if (row.proxy_ids) {
    try {
      const parsed = JSON.parse(row.proxy_ids)
      proxyIDs = Array.isArray(parsed) ? parsed : []
    } catch {
      proxyIDs = []
    }
  }
  groupForm.value = {
    Name: row.name,
    Type: row.type,
    ProxyIDs: proxyIDs,
    URL: row.url || 'http://www.gstatic.com/generate_204',
    Interval: row.interval || 300
  }
  await loadNodes()
  formDialogVisible.value = true
}

const handleSave = async () => {
  if (!groupForm.value.Name) {
    ElMessage.warning('请输入代理组名称')
    return
  }

  if (groupForm.value.ProxyIDs.length === 0) {
    ElMessage.warning('请至少选择一个节点')
    return
  }

  const data = {
    Name: groupForm.value.Name,
    Type: groupForm.value.Type,
    ProxyIDs: groupForm.value.ProxyIDs,
    URL: groupForm.value.URL,
    Interval: groupForm.value.Interval
  }

  if (isEdit.value) {
    await updateGroup(editId.value, data)
    ElMessage.success('更新成功')
  } else {
    await createGroup(data)
    ElMessage.success('创建成功')
  }
  formDialogVisible.value = false
  loadGroups()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该代理组吗？', '提示', { type: 'warning' })
  await deleteGroup(row.id)
  ElMessage.success('删除成功')
  loadGroups()
}

onMounted(() => {
  loadGroups()
})
</script>

<style scoped>
.groups-page {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  padding: 0;
}

:deep(.el-table) {
  border: none;
}

:deep(.el-table__header-wrapper) {
  background: #fafafa;
}

:deep(.el-table th) {
  background: #fafafa;
  color: #606266;
  font-weight: 500;
}

/* 穿梭框样式优化 */
:deep(.el-transfer) {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  max-width: 580px;
}

:deep(.el-transfer-panel) {
  width: 240px;
  flex: 0 0 240px;
}

:deep(.el-transfer__buttons) {
  padding: 0 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex-shrink: 0;
}

:deep(.el-transfer__button) {
  display: block;
  margin: 0;
  padding: 10px 14px;
  border-radius: 6px;
  transition: all 0.3s;
}

:deep(.el-transfer__button:hover) {
  transform: scale(1.05);
}

:deep(.el-transfer__button.is-with-texts) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

:deep(.el-transfer__button .el-icon) {
  font-size: 16px;
}
</style>
