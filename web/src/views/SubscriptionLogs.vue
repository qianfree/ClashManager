<template>
  <div class="subscription-logs-page">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon total">
              <el-icon><svg viewBox="0 0 1024 1024" width="24" height="24"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.total_subscriptions || 0 }}</div>
              <div class="stat-label">总订阅次数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon success">
              <el-icon><svg viewBox="0 0 1024 1024" width="24" height="24"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path><path fill="currentColor" d="M384 699.6 266.4 582 190.4 658l193.6 193.6 448-448-76-76L384 699.6z"></path></svg></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.success_count || 0 }}</div>
              <div class="stat-label">成功次数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon fail">
              <el-icon><svg viewBox="0 0 1024 1024" width="24" height="24"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path><path fill="currentColor" d="M664 320 448 536 232 320l-76 76 292 292 292-292-76-76z"></path></svg></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.fail_count || 0 }}</div>
              <div class="stat-label">失败次数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon rate">
              <el-icon><svg viewBox="0 0 1024 1024" width="24" height="24"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ successRate }}%</div>
              <div class="stat-label">成功率</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Top 用户 -->
    <el-card shadow="never" class="top-users-card" v-if="stats.top_users && stats.top_users.length > 0">
      <template #header>
        <div class="card-header">
          <span>活跃用户 Top 5</span>
        </div>
      </template>
      <div class="top-users-list">
        <div v-for="(user, index) in stats.top_users" :key="user.username" class="top-user-item">
          <span class="user-rank" :class="'rank-' + (index + 1)">{{ index + 1 }}</span>
          <span class="user-name">{{ user.username }}</span>
          <span class="user-count">{{ user.count }} 次</span>
        </div>
      </div>
    </el-card>

    <!-- 日志列表 -->
    <el-card shadow="never" class="logs-card">
      <template #header>
        <div class="card-header">
          <span>订阅日志</span>
          <div class="header-actions">
            <el-select v-model="successFilter" placeholder="筛选状态" clearable style="width: 120px; margin-right: 10px" @change="handleFilterChange">
              <el-option label="全部" value="" />
              <el-option label="成功" :value="true" />
              <el-option label="失败" :value="false" />
            </el-select>
            <el-select v-model="statsDays" placeholder="统计周期" style="width: 120px; margin-right: 10px" @change="loadStats">
              <el-option label="最近7天" :value="7" />
              <el-option label="最近30天" :value="30" />
              <el-option label="最近90天" :value="90" />
            </el-select>
            <el-button type="danger" :icon="Delete" @click="handleDeleteOldLogs">清理旧日志</el-button>
          </div>
        </div>
      </template>

      <el-table :data="logs" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user.username" label="用户" width="120" />
        <el-table-column prop="token" label="Token" width="120">
          <template #default="{ row }">
            <code>{{ row.token || '-' }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP地址" width="150" />
        <el-table-column prop="userAgent" label="User Agent" min-width="200" show-overflow-tooltip />
        <el-table-column prop="success" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'" size="small">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="error" label="错误信息" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.error || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="访问时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import { getSubscriptionLogs, getSubscriptionStats, deleteOldLogs } from '@/api/subscription'

const loading = ref(false)
const logs = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const successFilter = ref('')
const statsDays = ref(30)
const stats = ref({})

const successRate = computed(() => {
  if (!stats.value.total_subscriptions) return 0
  return ((stats.value.success_count || 0) / stats.value.total_subscriptions * 100).toFixed(1)
})

// 加载统计数据
const loadStats = async () => {
  try {
    const data = await getSubscriptionStats({ days: statsDays.value })
    stats.value = data
  } catch {
    ElMessage.error('加载统计数据失败')
  }
}

// 加载日志列表
const loadLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    if (successFilter.value !== '') {
      params.success = successFilter.value === 'true' ? true : successFilter.value === 'false' ? false : successFilter.value
    }
    const data = await getSubscriptionLogs(params)
    logs.value = data.logs || []
    total.value = data.total || 0
  } catch {
    ElMessage.error('加载日志失败')
  } finally {
    loading.value = false
  }
}

// 筛选变化
const handleFilterChange = () => {
  currentPage.value = 1
  loadLogs()
}

// 页码变化
const handlePageChange = () => {
  loadLogs()
}

// 每页数量变化
const handleSizeChange = () => {
  currentPage.value = 1
  loadLogs()
}

// 清理旧日志
const handleDeleteOldLogs = () => {
  ElMessageBox.prompt('请输入要保留的天数（删除此天数之前的日志）', '清理旧日志', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /^[7-9][0-9]$|^[1-9][0-9][0-9]+$/,
    inputErrorMessage: '请输入7天以上的数字',
    inputValue: '90'
  }).then(({ value }) => {
    deleteOldLogs(parseInt(value)).then(() => {
      ElMessage.success('清理成功')
      loadLogs()
      loadStats()
    }).catch(() => {
      ElMessage.error('清理失败')
    })
  }).catch(() => {
    // 取消
  })
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

onMounted(() => {
  loadStats()
  loadLogs()
})
</script>

<style scoped>
.subscription-logs-page {
  height: 100%;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  cursor: default;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.success {
  background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
}

.stat-icon.fail {
  background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);
}

.stat-icon.rate {
  background: linear-gradient(135deg, #f6d365 0%, #fda085 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}

.top-users-card {
  margin-bottom: 20px;
}

.top-users-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.top-user-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px 15px;
  background: #f5f7fa;
  border-radius: 6px;
}

.user-rank {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 14px;
  background: #e4e7ed;
  color: #606266;
}

.user-rank.rank-1 {
  background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
  color: #fff;
}

.user-rank.rank-2 {
  background: linear-gradient(135deg, #c0c0c0 0%, #e0e0e0 100%);
  color: #fff;
}

.user-rank.rank-3 {
  background: linear-gradient(135deg, #cd7f32 0%, #e8a870 100%);
  color: #fff;
}

.user-name {
  flex: 1;
  font-weight: 500;
  color: #303133;
}

.user-count {
  font-size: 14px;
  color: #909399;
}

.logs-card {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}

code {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #e6a23c;
  background: #fdf6ec;
  padding: 2px 6px;
  border-radius: 3px;
}

:deep(.el-pagination) {
  display: flex;
  justify-content: flex-end;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style>
