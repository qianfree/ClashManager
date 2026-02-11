<template>
  <div class="rules-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>规则列表</span>
          </div>
          <div class="header-actions">
            <el-button :icon="Upload" @click="showImportDialog">导入规则</el-button>
            <el-button type="primary" :icon="Plus" @click="showCreateDialog">新增规则</el-button>
          </div>
        </div>
      </template>

      <!-- 搜索和过滤区域 -->
      <div class="filter-section">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索匹配内容、目标或备注"
          :prefix-icon="Search"
          clearable
          style="width: 280px; margin-right: 10px;"
          @change="handleSearchChange"
        />
        <el-select
          v-model="filterType"
          placeholder="全部规则类型"
          clearable
          style="width: 180px; margin-right: 10px;"
          @change="handleSearchChange"
        >
          <el-option label="DOMAIN-SUFFIX" value="DOMAIN-SUFFIX" />
          <el-option label="DOMAIN" value="DOMAIN" />
          <el-option label="DOMAIN-KEYWORD" value="DOMAIN-KEYWORD" />
          <el-option label="IP-CIDR" value="IP-CIDR" />
          <el-option label="GEOIP" value="GEOIP" />
          <el-option label="MATCH" value="MATCH" />
        </el-select>
        <el-button @click="resetFilter">重置</el-button>
        <div style="margin-left: auto; color: #909399; font-size: 14px;">
          共 {{ total }} 条规则
        </div>
      </div>

      <el-table :data="rules" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="ID" label="ID" min-width="60" />
        <el-table-column prop="Priority" label="序号" min-width="70">
          <template #default="{ row }">
            <span>{{ row.Priority ?? 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="Type" label="规则类型" min-width="130">
          <template #default="{ row }">
            <el-tag size="small">{{ row.Type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="Payload" label="匹配内容" min-width="180" show-overflow-tooltip />
        <el-table-column prop="Target" label="目标" min-width="130">
          <template #default="{ row }">
            <el-tag v-if="row.Target === 'PROXY'" type="primary" size="small">PROXY</el-tag>
            <el-tag v-else-if="row.Target === 'DIRECT'" type="success" size="small">DIRECT</el-tag>
            <el-tag v-else-if="row.Target === 'REJECT'" type="danger" size="small">REJECT</el-tag>
            <el-tag v-else-if="row.TargetType === 'node'" type="info" size="small">{{ row.Target }}</el-tag>
            <el-tag v-else type="warning" size="small">{{ row.Target }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="Remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column prop="CreatedAt" label="创建时间" min-width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100, 200]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 新增/编辑规则对话框 -->
    <el-dialog v-model="formDialogVisible" :title="isEdit ? '编辑规则' : '新增规则'" width="500px">
      <el-form :model="ruleForm" label-width="100px">
        <el-form-item label="序号">
          <el-input-number v-model="ruleForm.Priority" :min="0" :max="9999" placeholder="数字越小优先级越高" style="width: 100%" />
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">数字越小优先级越高，0为最高优先级</div>
        </el-form-item>
        <el-form-item label="规则类型">
          <el-select v-model="ruleForm.Type" placeholder="请选择规则类型" style="width: 100%">
            <el-option label="DOMAIN-SUFFIX - 域名后缀匹配" value="DOMAIN-SUFFIX" />
            <el-option label="DOMAIN - 完整域名匹配" value="DOMAIN" />
            <el-option label="DOMAIN-KEYWORD - 域名关键字匹配" value="DOMAIN-KEYWORD" />
            <el-option label="IP-CIDR - IP段匹配" value="IP-CIDR" />
            <el-option label="GEOIP - 地理位置匹配" value="GEOIP" />
            <el-option label="MATCH - 全匹配（默认规则）" value="MATCH" />
          </el-select>
        </el-form-item>
        <el-form-item label="匹配内容">
          <el-input v-model="ruleForm.Payload" placeholder="如: google.com 或 192.168.1.0/24" />
        </el-form-item>
        <el-form-item label="目标">
          <el-select v-model="ruleForm.Target" placeholder="请选择目标" style="width: 100%" @change="handleTargetChange" filterable>
            <el-option-group label="内置目标">
              <el-option label="PROXY - 代理" value="PROXY">
                <span>PROXY</span>
                <span style="color: #8492a6; font-size: 12px; margin-left: 8px;">代理</span>
              </el-option>
              <el-option label="DIRECT - 直连" value="DIRECT">
                <span>DIRECT</span>
                <span style="color: #8492a6; font-size: 12px; margin-left: 8px;">直连</span>
              </el-option>
              <el-option label="REJECT - 拒绝" value="REJECT">
                <span>REJECT</span>
                <span style="color: #8492a6; font-size: 12px; margin-left: 8px;">拒绝</span>
              </el-option>
            </el-option-group>
            <el-option-group label="代理节点">
              <el-option
                v-for="node in nodes"
                :key="node.ID"
                :label="node.Name"
                :value="`node:${node.ID}:${node.Name}`"
              >
                <span>{{ node.Name }}</span>
                <span style="color: #8492a6; font-size: 12px; margin-left: 8px;">{{ node.Type }}</span>
              </el-option>
            </el-option-group>
            <el-option-group label="代理组">
              <el-option
                v-for="group in groups"
                :key="group.ID"
                :label="group.Name"
                :value="`group:${group.ID}:${group.Name}`"
              >
                <span>{{ group.Name }}</span>
                <span style="color: #8492a6; font-size: 12px; margin-left: 8px;">{{ group.Type }}</span>
              </el-option>
            </el-option-group>
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="ruleForm.Remark" placeholder="可选，用于记录规则用途" maxlength="200" show-word-limit />
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">备注仅用于展示，不会生成到配置文件中</div>
        </el-form-item>
        <el-form-item label="No Resolve">
          <el-switch v-model="ruleForm.NoResolve" />
          <span style="margin-left: 10px; color: #999; font-size: 12px;">不反查域名</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- 导入规则对话框 -->
    <el-dialog v-model="importDialogVisible" title="导入规则" width="600px">
      <el-form label-width="100px">
        <el-form-item label="选择文件">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :show-file-list="true"
            :limit="1"
            accept=".yaml,.yml"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
          >
            <el-button type="primary">选择 YAML 文件</el-button>
            <template #tip>
              <div style="color: #909399; font-size: 12px; margin-top: 5px;">
                支持 .yaml 或 .yml 格式文件，将自动解析 rules 节点
              </div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item label="文件内容" v-if="importContent">
          <el-input
            v-model="importContent"
            type="textarea"
            :rows="10"
            placeholder="文件内容将显示在这里"
            readonly
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleImport" :loading="importing">导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Upload } from '@element-plus/icons-vue'
import { getRules, createRule, updateRule, deleteRule, importRules } from '@/api/rules'
import { getGroups } from '@/api/groups'
import { getNodes } from '@/api/nodes'

const rules = ref([])
const groups = ref([])
const nodes = ref([])
const formDialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const loading = ref(false)

// 导入相关
const importDialogVisible = ref(false)
const importContent = ref('')
const importing = ref(false)

const ruleForm = ref({
  Type: 'DOMAIN-SUFFIX',
  Payload: '',
  Target: 'PROXY',
  TargetID: 0,
  TargetType: '',
  Priority: 0,
  NoResolve: false,
  Remark: ''
})

// 搜索和过滤
const searchKeyword = ref('')
const filterType = ref('')

// 分页相关
const currentPage = ref(1)
const pageSize = ref(50)
const total = ref(0)

// 格式化日期时间
const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadRules = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    if (filterType.value) params.type = filterType.value
    if (searchKeyword.value) params.keyword = searchKeyword.value

    const result = await getRules(params)
    rules.value = result.rules || []
    total.value = result.total || 0
  } catch (error) {
    ElMessage.error('加载规则失败')
  } finally {
    loading.value = false
  }
}

const loadGroups = async () => {
  groups.value = await getGroups()
}

const loadNodes = async () => {
  nodes.value = await getNodes()
}

const handleTargetChange = (value) => {
  // Check if it's a node reference (format: "node:ID:Name")
  if (value && value.startsWith('node:')) {
    const parts = value.split(':')
    if (parts.length >= 3) {
      ruleForm.value.TargetID = parseInt(parts[1])
      ruleForm.value.TargetType = 'node'
      ruleForm.value.Target = parts[2] // Store name for display
    }
  } else if (value && value.startsWith('group:')) {
    // Check if it's a group reference (format: "group:ID:Name")
    const parts = value.split(':')
    if (parts.length >= 3) {
      ruleForm.value.TargetID = parseInt(parts[1])
      ruleForm.value.TargetType = 'group'
      ruleForm.value.Target = parts[2] // Store name for display
    }
  } else {
    // Built-in target (PROXY, DIRECT, REJECT)
    ruleForm.value.TargetID = 0
    ruleForm.value.TargetType = ''
    ruleForm.value.Target = value
  }
}

const showCreateDialog = async () => {
  isEdit.value = false
  editId.value = null
  ruleForm.value = {
    Type: 'DOMAIN-SUFFIX',
    Payload: '',
    Target: 'PROXY',
    TargetID: 0,
    TargetType: '',
    Priority: 0,
    NoResolve: false,
    Remark: ''
  }
  await Promise.all([loadGroups(), loadNodes()])
  formDialogVisible.value = true
}

const handleEdit = async (row) => {
  isEdit.value = true
  editId.value = row.ID
  // Build Target value from TargetID and TargetType
  let targetValue = row.Target || 'PROXY'
  if (row.TargetID > 0 && row.TargetType === 'node') {
    targetValue = `node:${row.TargetID}:${row.Target}`
  } else if (row.TargetID > 0 && row.TargetType === 'group') {
    targetValue = `group:${row.TargetID}:${row.Target}`
  }
  ruleForm.value = {
    Type: row.Type,
    Payload: row.Payload,
    Target: targetValue,
    TargetID: row.TargetID || 0,
    TargetType: row.TargetType || '',
    Priority: row.Priority ?? 0,
    NoResolve: row.NoResolve,
    Remark: row.Remark || ''
  }
  await Promise.all([loadGroups(), loadNodes()])
  formDialogVisible.value = true
}

const handleSave = async () => {
  if (!ruleForm.value.Payload || !ruleForm.value.Target) {
    ElMessage.warning('请填写完整信息')
    return
  }
  const data = {
    Type: ruleForm.value.Type,
    Payload: ruleForm.value.Payload,
    Target: ruleForm.value.Target,
    TargetID: ruleForm.value.TargetID || 0,
    TargetType: ruleForm.value.TargetType || '',
    Priority: ruleForm.value.Priority ?? 0,
    NoResolve: ruleForm.value.NoResolve,
    Remark: ruleForm.value.Remark || ''
  }
  if (isEdit.value) {
    await updateRule(editId.value, data)
    ElMessage.success('更新成功')
  } else {
    await createRule(data)
    ElMessage.success('创建成功')
  }
  formDialogVisible.value = false
  loadRules()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该规则吗？', '提示', { type: 'warning' })
  await deleteRule(row.ID)
  ElMessage.success('删除成功')
  loadRules()
}

const handleSearchChange = () => {
  currentPage.value = 1
  loadRules()
}

const resetFilter = () => {
  searchKeyword.value = ''
  filterType.value = ''
  currentPage.value = 1
  loadRules()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadRules()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadRules()
}

// 导入相关函数
const showImportDialog = () => {
  importContent.value = ''
  importDialogVisible.value = true
}

const handleFileChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    importContent.value = e.target.result
  }
  reader.readAsText(file.raw)
}

const handleFileRemove = () => {
  importContent.value = ''
}

const handleImport = async () => {
  if (!importContent.value) {
    ElMessage.warning('请先选择文件')
    return
  }

  importing.value = true
  try {
    const result = await importRules(importContent.value)
    ElMessage.success(`成功导入 ${result.count} 条规则`)
    importDialogVisible.value = false
    importContent.value = ''
    loadRules()
  } catch (error) {
    ElMessage.error('导入失败: ' + (error.message || '未知错误'))
  } finally {
    importing.value = false
  }
}

onMounted(() => {
  loadRules()
})
</script>

<style scoped>
.rules-page {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
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

.filter-section {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  background: #f9f9f9;
  border-bottom: 1px solid #f0f0f0;
}

.pagination-section {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid #f0f0f0;
}
</style>
