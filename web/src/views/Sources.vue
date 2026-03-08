<template>
  <div class="sources-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20">
              <path d="M853.333333 512c0 235.648-191.018667 426.666667-426.666666 426.666667S0 747.648 0 512 191.018667 85.333333 426.666667 85.333333 853.333333 276.352 853.333333 512z" fill="#E5E8EF" opacity="0.3"/>
              <path fill="currentColor" d="M426.666667 170.666667C235.392 170.666667 85.333333 324.053333 85.333333 512s150.058667 341.333333 341.333334 341.333333 341.333333-153.386667 341.333333-341.333333S617.941333 170.666667 426.666667 170.666667z m0 640c-165.205333 0-298.666667-137.216-298.666667-298.666667S261.461333 213.333333 426.666667 213.333333s298.666667 137.216 298.666666 298.666667-133.461333 298.666667-298.666666 298.666667z"/>
              <path fill="currentColor" d="M554.666667 426.666667h-85.333334V298.666667c0-23.573333-19.413333-42.666667-42.666666-42.666667s-42.666667 19.093333-42.666667 42.666667v170.666666c0 23.573333 19.413333 42.666667 42.666667 42.666667h128c23.573333 0 42.666667-19.093333 42.666666-42.666667s-19.093333-42.666667-42.666666-42.666666z"/>
            </svg></el-icon>
            <span>订阅源管理</span>
          </div>
          <el-button type="primary" :icon="Plus" @click="showCreateDialog">添加订阅源</el-button>
        </div>
      </template>

      <el-table :data="sources" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="url" label="订阅链接" min-width="250" show-overflow-tooltip />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" @change="handleToggleEnabled(row)" />
          </template>
        </el-table-column>
        <el-table-column label="同步模式" width="120">
          <template #default="{ row }">
            <el-tag :type="getSyncModeType(row.sync_mode)">
              {{ getSyncModeText(row.sync_mode) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="更新间隔" width="100">
          <template #default="{ row }">
            {{ row.update_interval }}小时
          </template>
        </el-table-column>
        <el-table-column label="最后同步" width="160">
          <template #default="{ row }">
            <span v-if="row.last_sync">{{ formatTime(row.last_sync) }}</span>
            <span v-else style="color: #909399;">未同步</span>
          </template>
        </el-table-column>
        <el-table-column label="节点标签" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.node_tag" size="small" type="info">{{ row.node_tag }}</el-tag>
            <span v-else style="color: #909399;">-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleSync(row)" :loading="syncingId === row.id">
              同步
            </el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑订阅源对话框 -->
    <el-dialog v-model="formDialogVisible" :title="isEdit ? '编辑订阅源' : '添加订阅源'" width="500px">
      <el-form :model="sourceForm" label-width="100px" :rules="formRules" ref="sourceFormRef">
        <el-form-item label="名称" prop="name">
          <el-input v-model="sourceForm.name" placeholder="请输入订阅源名称" />
        </el-form-item>
        <el-form-item label="订阅链接" prop="url">
          <el-input v-model="sourceForm.url" placeholder="请输入订阅链接" />
        </el-form-item>
        <el-form-item label="节点标签" prop="node_tag">
          <el-input v-model="sourceForm.node_tag" placeholder="可选，用于标识节点来源" />
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">
            留空则使用订阅源名称作为标签
          </div>
        </el-form-item>
        <el-form-item label="同步模式" prop="sync_mode">
          <el-select v-model="sourceForm.sync_mode" placeholder="请选择同步模式">
            <el-option label="追加模式" value="append">
              <div>
                <div>追加模式</div>
                <div style="color: #909399; font-size: 12px;">保留现有节点，添加新节点</div>
              </div>
            </el-option>
            <el-option label="替换模式" value="replace">
              <div>
                <div>替换模式</div>
                <div style="color: #909399; font-size: 12px;">清空现有节点，只保留订阅节点</div>
              </div>
            </el-option>
            <el-option label="智能合并" value="smart">
              <div>
                <div>智能合并</div>
                <div style="color: #909399; font-size: 12px;">按名称合并，更新已存在的节点</div>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="更新间隔">
          <el-input-number v-model="sourceForm.update_interval" :min="1" :max="168" />
          <span style="margin-left: 8px;">小时</span>
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="sourceForm.enabled" />
        </el-form-item>
        <el-form-item>
          <el-button @click="handleTest" :loading="testing">测试连接</el-button>
          <span v-if="testResult" style="margin-left: 12px; font-size: 13px;">
            <span v-if="testResult.success" style="color: #67c23a;">
              成功！发现 {{ testResult.nodesCount }} 个节点
            </span>
            <span v-else style="color: #f56c6c;">
              {{ testResult.error }}
            </span>
          </span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getSources, createSource, updateSource, deleteSource, syncSource, testSource } from '@/api/sources'

const sources = ref([])
const formDialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const syncingId = ref(null)
const saving = ref(false)
const testing = ref(false)
const testResult = ref(null)
const sourceFormRef = ref(null)

const sourceForm = ref({
  name: '',
  url: '',
  enabled: true,
  update_interval: 24,
  node_tag: '',
  sync_mode: 'append'
})

const formRules = {
  name: [{ required: true, message: '请输入订阅源名称', trigger: 'blur' }],
  url: [
    { required: true, message: '请输入订阅链接', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ],
  sync_mode: [{ required: true, message: '请选择同步模式', trigger: 'change' }]
}

const loadSources = async () => {
  sources.value = await getSources()
}

const getSyncModeType = (mode) => {
  const types = {
    append: '',
    replace: 'danger',
    smart: 'success'
  }
  return types[mode] || ''
}

const getSyncModeText = (mode) => {
  const texts = {
    append: '追加',
    replace: '替换',
    smart: '智能'
  }
  return texts[mode] || mode
}

const formatTime = (timeStr) => {
  if (!timeStr) return '-'
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const showCreateDialog = () => {
  isEdit.value = false
  editId.value = null
  testResult.value = null
  sourceForm.value = {
    name: '',
    url: '',
    enabled: true,
    update_interval: 24,
    node_tag: '',
    sync_mode: 'append'
  }
  formDialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  editId.value = row.id
  testResult.value = null
  sourceForm.value = {
    name: row.name,
    url: row.url,
    enabled: row.enabled,
    update_interval: row.update_interval,
    node_tag: row.node_tag || '',
    sync_mode: row.sync_mode || 'append'
  }
  formDialogVisible.value = true
}

const handleTest = async () => {
  if (!sourceForm.value.url) {
    ElMessage.warning('请先输入订阅链接')
    return
  }

  testing.value = true
  testResult.value = null
  try {
    const result = await testSource(sourceForm.value.url)
    testResult.value = result
    if (result.success) {
      ElMessage.success(`测试成功！发现 ${result.nodesCount} 个节点`)
      if (result.preview && result.preview.length > 0) {
        console.log('节点预览:', result.preview)
      }
    }
  } catch (error) {
    testResult.value = { success: false, error: error.message || '测试失败' }
  } finally {
    testing.value = false
  }
}

const handleSave = async () => {
  const valid = await sourceFormRef.value?.validate().catch(() => false)
  if (!valid) {
    return
  }

  saving.value = true
  try {
    const data = {
      name: sourceForm.value.name,
      url: sourceForm.value.url,
      enabled: sourceForm.value.enabled,
      update_interval: sourceForm.value.update_interval,
      node_tag: sourceForm.value.node_tag,
      sync_mode: sourceForm.value.sync_mode
    }

    if (isEdit.value) {
      await updateSource(editId.value, data)
      ElMessage.success('更新成功')
    } else {
      await createSource(data)
      ElMessage.success('创建成功')
    }
    formDialogVisible.value = false
    loadSources()
  } finally {
    saving.value = false
  }
}

const handleToggleEnabled = async (row) => {
  try {
    await updateSource(row.id, {
      name: row.name,
      url: row.url,
      enabled: row.enabled,
      update_interval: row.update_interval,
      node_tag: row.node_tag,
      sync_mode: row.sync_mode
    })
    ElMessage.success(row.enabled ? '已启用' : '已禁用')
  } catch (error) {
    row.enabled = !row.enabled
    ElMessage.error('操作失败')
  }
}

const handleSync = async (row) => {
  await ElMessageBox.confirm(
    `确定要同步订阅源「${row.name}」吗？\n\n同步模式：${getSyncModeText(row.sync_mode)}`,
    '确认同步',
    { type: 'warning' }
  )

  syncingId.value = row.id
  try {
    const result = await syncSource(row.id)
    ElMessage.success(`同步成功！共 ${result.nodesCount} 个节点`)
    loadSources()
  } catch (error) {
    ElMessage.error('同步失败：' + (error.message || '未知错误'))
  } finally {
    syncingId.value = null
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm(`确定删除订阅源「${row.name}」吗？`, '提示', { type: 'warning' })
  await deleteSource(row.id)
  ElMessage.success('删除成功')
  loadSources()
}

onMounted(() => {
  loadSources()
})
</script>

<style scoped>
.sources-page {
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

:deep(.el-select-dropdown__item) {
  height: auto;
  padding: 8px 12px;
}
</style>
