<template>
  <div class="subscription-page" v-loading="loading">
    <div class="card-container">
      <el-card shadow="never">
        <template #header>
          <div class="card-header">
            <div class="header-left">
              <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
              <span>订阅链接</span>
            </div>
            <el-button
              type="primary"
              :icon="RefreshRight"
              :loading="refreshing"
              @click="handleRefreshToken"
              size="small"
            >
              刷新Token
            </el-button>
          </div>
        </template>

        <el-form label-width="100px">
          <el-form-item label="订阅地址">
            <el-input
              :model-value="subscriptionUrl"
              readonly
            >
              <template #append>
                <el-button @click="copyUrl">复制</el-button>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="订阅Token">
            <el-input
              :model-value="token"
              readonly
            >
              <template #append>
                <el-button @click="copyToken">复制</el-button>
              </template>
            </el-input>
          </el-form-item>
        </el-form>

        <el-alert
          title="使用说明"
          type="info"
          :closable="false"
          show-icon
        >
          <p>1. 复制上方订阅链接或Token</p>
          <p>2. 在Clash客户端中添加订阅</p>
          <p>3. 客户端会自动获取最新配置</p>
          <p>4. 如需更换订阅地址，点击右上角"刷新Token"按钮</p>
        </el-alert>
      </el-card>

      <el-card shadow="never">
        <template #header>
          <div class="card-header">
            <div class="header-left">
              <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
              <span>二维码</span>
            </div>
          </div>
        </template>

        <div class="qrcode-container">
          <img :src="qrcodeUrl" alt="订阅二维码" />
          <p>使用Clash客户端扫描二维码导入</p>
        </div>
      </el-card>
    </div>

    <el-card shadow="never" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>配置预览</span>
          </div>
          <el-button type="primary" :icon="View" @click="loadConfig">加载配置</el-button>
        </div>
      </template>

      <pre v-if="configContent" class="config-content">{{ configContent }}</pre>
      <el-empty v-else description="点击上方按钮加载配置" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { View, RefreshRight } from '@element-plus/icons-vue'
import { getSubscriptionURL, refreshToken } from '@/api/subscription'
import axios from 'axios'

const subscriptionUrl = ref('')
const token = ref('')
const qrcodeUrl = ref('')
const configContent = ref('')
const loading = ref(false)
const refreshing = ref(false)

// 加载订阅信息
const loadSubscription = async () => {
  loading.value = true
  try {
    const data = await getSubscriptionURL()
    subscriptionUrl.value = data.url
    token.value = data.token
    qrcodeUrl.value = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(data.url)}`
  } catch {
    ElMessage.error('加载订阅信息失败')
  } finally {
    loading.value = false
  }
}

// 刷新Token
const handleRefreshToken = async () => {
  refreshing.value = true
  try {
    const data = await refreshToken()
    ElMessage.success('Token已刷新')
    await loadSubscription()
  } catch {
    ElMessage.error('刷新Token失败')
  } finally {
    refreshing.value = false
  }
}

// 复制URL
const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(subscriptionUrl.value)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

// 复制Token
const copyToken = async () => {
  try {
    await navigator.clipboard.writeText(token.value)
    ElMessage.success('Token已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

// 加载配置预览
const loadConfig = async () => {
  if (!token.value) {
    ElMessage.error('Token不存在，请刷新页面')
    return
  }
  try {
    // 使用相对路径通过Vite代理访问后端
    const response = await axios.get(`/sub/${token.value}`, {
      headers: { Accept: 'application/yaml' }
    })
    configContent.value = response.data
    ElMessage.success('配置加载成功')
  } catch (err) {
    console.error('加载配置失败:', err)
    ElMessage.error('加载配置失败: ' + (err.response?.data?.error || err.message))
  }
}

onMounted(() => {
  loadSubscription()
})
</script>

<style scoped>
.subscription-page {
  height: 100%;
}

.card-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
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

.qrcode-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  gap: 15px;
}

.qrcode-container img {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 10px;
  background: #fff;
}

.qrcode-container p {
  color: #909399;
  font-size: 14px;
  margin: 0;
}

.config-content {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 13px;
  line-height: 1.6;
  max-height: 500px;
  overflow-y: auto;
  color: #303133;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

:deep(.el-alert) {
  margin-top: 15px;
}

:deep(.el-alert p) {
  margin: 5px 0;
}

:deep(.el-empty) {
  padding: 40px 0;
}
</style>
