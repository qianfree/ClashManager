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
        <el-table-column prop="ID" label="ID" min-width="60" />
        <el-table-column prop="Name" label="名称" min-width="120" />
        <el-table-column prop="Type" label="类型" min-width="120">
          <template #default="{ row }">
            <el-tag :type="row.Type === 'select' ? 'primary' : 'success'">
              {{ row.Type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ProxyNodes" label="包含节点" min-width="250">
          <template #default="{ row }">
            <el-tag v-for="proxy in displayNodes(row)" :key="proxy.ID || proxy" size="small" style="margin-right: 5px;">
              {{ proxy.Name || proxy }}
            </el-tag>
            <span v-if="displayNodes(row).length === 0" style="color: #909399;">无</span>
          </template>
        </el-table-column>
        <el-table-column prop="URL" label="测试URL" min-width="200" show-overflow-tooltip />
        <el-table-column prop="Interval" label="间隔(秒)" min-width="100" />
        <el-table-column label="操作" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑代理组对话框 -->
    <el-dialog v-model="formDialogVisible" :title="isEdit ? '编辑代理组' : '新增代理组'" width="500px">
      <el-form :model="groupForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="groupForm.Name" placeholder="请输入代理组名称" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="groupForm.Type" placeholder="请选择类型">
            <el-option label="手动选择" value="select" />
            <el-option label="自动测速" value="url-test" />
            <el-option label="负载均衡" value="load-balance" />
          </el-select>
        </el-form-item>
        <el-form-item label="包含节点">
          <el-select
            v-model="groupForm.ProxyIDs"
            multiple
            placeholder="请选择节点"
            style="width: 100%"
          >
            <el-option
              v-for="node in nodes"
              :key="node.ID"
              :label="node.Name"
              :value="node.ID"
            />
          </el-select>
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">
            已选择 {{ groupForm.ProxyIDs.length }} 个节点
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
import { ref, onMounted } from 'vue'
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

const loadGroups = async () => {
  groups.value = await getGroups()
}

// 兼容新旧数据格式显示节点
const displayNodes = (row) => {
  // 新格式：使用 ProxyNodes 数组
  if (row.ProxyNodes && row.ProxyNodes.length > 0) {
    return row.ProxyNodes
  }
  // 旧格式：解析 Proxies JSON 字符串
  if (row.Proxies) {
    try {
      const proxies = JSON.parse(row.Proxies)
      if (Array.isArray(proxies)) {
        return proxies.map(p => ({ ID: null, Name: p }))
      }
    } catch {
      return []
    }
  }
  return []
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
  editId.value = row.ID
  // Parse ProxyIDs from JSON if needed
  let proxyIDs = []
  if (row.ProxyIDs) {
    try {
      proxyIDs = JSON.parse(row.ProxyIDs)
    } catch {
      proxyIDs = []
    }
  }
  groupForm.value = {
    Name: row.Name,
    Type: row.Type,
    ProxyIDs: proxyIDs,
    URL: row.URL || 'http://www.gstatic.com/generate_204',
    Interval: row.Interval || 300
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
  await deleteGroup(row.ID)
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
</style>
