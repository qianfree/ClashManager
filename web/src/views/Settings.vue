<template>
  <div class="settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>DNS 设置</span>
          </div>
          <el-button type="primary" :icon="Check" @click="handleSave" :loading="saving">保存配置</el-button>
        </div>
      </template>

      <el-form :model="dnsForm" label-width="140px" v-loading="loading">
        <el-divider content-position="left">基础设置</el-divider>

        <el-form-item label="启用 DNS">
          <el-switch v-model="dnsForm.enable" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">关闭后将使用系统DNS</span>
        </el-form-item>

        <el-form-item label="监听地址">
          <el-input v-model="dnsForm.listen" placeholder="0.0.0.0:53" style="width: 300px;" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">DNS服务器监听地址</span>
        </el-form-item>

        <el-form-item label="增强模式">
          <el-select v-model="dnsForm.enhancedMode" style="width: 300px;">
            <el-option label="Fake-IP" value="fake-ip" />
            <el-option label="Redir-Host" value="redir-host" />
            <el-option label="Mapping" value="mapping" />
          </el-select>
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">DNS处理模式</span>
        </el-form-item>

        <el-divider content-position="left">DNS 服务器</el-divider>

        <!-- 主DNS服务器 -->
        <el-form-item label="主DNS服务器">
          <div class="dns-input-container">
            <div class="input-wrapper">
              <el-input
                v-model="nameserverInput"
                placeholder="输入DNS地址（如：223.5.5.5）"
                @keyup.enter="addNameserver"
                clearable
                style="flex: 1;"
              >
                <template #append>
                  <el-button @click="addNameserver">添加</el-button>
                </template>
              </el-input>
            </div>
            <div class="quick-add">
              <span class="quick-label">快速添加：</span>
              <el-button size="small" @click="quickAddNameserver('223.5.5.5')">阿里</el-button>
              <el-button size="small" @click="quickAddNameserver('119.29.29.29')">腾讯</el-button>
              <el-button size="small" @click="quickAddNameserver('180.76.76.76')">百度</el-button>
            </div>
          </div>
          <div class="dns-tags" v-if="dnsForm.nameserver.length > 0">
            <el-tag
              v-for="(dns, index) in dnsForm.nameserver"
              :key="index"
              closable
              @close="removeNameserver(index)"
              class="dns-tag"
            >
              {{ dns }}
              <span v-if="isCommonDNS(dns)" class="dns-name">({{ getDNSName(dns) }})</span>
            </el-tag>
          </div>
          <div v-else class="empty-tip">暂无配置，请添加主DNS服务器</div>
        </el-form-item>

        <!-- 备用DNS服务器 -->
        <el-form-item label="备用DNS服务器">
          <div class="dns-input-container">
            <div class="input-wrapper">
              <el-input
                v-model="fallbackInput"
                placeholder="输入备用DNS地址"
                @keyup.enter="addFallback"
                clearable
                style="flex: 1;"
              >
                <template #append>
                  <el-button @click="addFallback">添加</el-button>
                </template>
              </el-input>
            </div>
            <div class="quick-add">
              <span class="quick-label">快速添加：</span>
              <el-button size="small" @click="quickAddFallback('8.8.8.8')">Google</el-button>
              <el-button size="small" @click="quickAddFallback('1.1.1.1')">Cloudflare</el-button>
              <el-button size="small" @click="quickAddFallback('208.67.222.222')">OpenDNS</el-button>
            </div>
          </div>
          <div class="dns-tags" v-if="dnsForm.fallback.length > 0">
            <el-tag
              v-for="(dns, index) in dnsForm.fallback"
              :key="index"
              closable
              @close="removeFallback(index)"
              class="dns-tag"
            >
              {{ dns }}
              <span v-if="isCommonDNS(dns)" class="dns-name">({{ getDNSName(dns) }})</span>
            </el-tag>
          </div>
          <div v-else class="empty-tip">暂无配置</div>
        </el-form-item>

        <!-- 默认DNS服务器 -->
        <el-form-item label="默认DNS服务器">
          <div class="dns-input-container">
            <div class="input-wrapper">
              <el-input
                v-model="defaultNameserverInput"
                placeholder="用于解析DNS服务器的DNS"
                @keyup.enter="addDefaultNameserver"
                clearable
                style="flex: 1;"
              >
                <template #append>
                  <el-button @click="addDefaultNameserver">添加</el-button>
                </template>
              </el-input>
            </div>
            <div class="quick-add">
              <span class="quick-label">快速添加：</span>
              <el-button size="small" @click="quickAddDefaultNameserver('114.114.114.114')">114DNS</el-button>
            </div>
          </div>
          <div class="dns-tags" v-if="dnsForm.defaultNameserver.length > 0">
            <el-tag
              v-for="(dns, index) in dnsForm.defaultNameserver"
              :key="index"
              closable
              @close="removeDefaultNameserver(index)"
              class="dns-tag"
            >
              {{ dns }}
            </el-tag>
          </div>
          <div v-else class="empty-tip">暂无配置</div>
        </el-form-item>

        <el-divider content-position="left">Fake-IP 过滤</el-divider>

        <el-form-item label="过滤列表">
          <div class="dns-input-container">
            <el-input
              v-model="fakeIPFilterInput"
              placeholder="输入域名（如：*.lan）"
              @keyup.enter="addFakeIPFilter"
              clearable
              style="flex: 1;"
            >
              <template #append>
                <el-button @click="addFakeIPFilter">添加</el-button>
              </template>
            </el-input>
          </div>
          <div class="dns-tags" v-if="dnsForm.fakeIPFilter.length > 0">
            <el-tag
              v-for="(filter, index) in dnsForm.fakeIPFilter"
              :key="index"
              closable
              @close="removeFakeIPFilter(index)"
              class="dns-tag"
            >
              {{ filter }}
            </el-tag>
          </div>
          <div v-else class="empty-tip">暂无配置</div>
        </el-form-item>
      </el-form>
    </el-card>

    <el-alert
      title="配置说明"
      type="info"
      :closable="false"
      show-icon
      style="margin-top: 20px;"
    >
      <p><strong>DNS 服务器说明：</strong></p>
      <ul style="margin: 5px 0; padding-left: 20px;">
        <li><strong>主DNS服务器</strong>：用于常规DNS查询，推荐使用国内DNS如阿里、腾讯</li>
        <li><strong>备用DNS服务器</strong>：当主DNS查询失败时使用，推荐使用国外DNS</li>
        <li><strong>默认DNS服务器</strong>：用于解析其他DNS服务器的地址，通常使用运营商DNS</li>
      </ul>
      <p style="margin: 10px 0 5px 0;"><strong>增强模式说明：</strong></p>
      <ul style="margin: 5px 0; padding-left: 20px;">
        <li><strong>Fake-IP</strong>：返回假的IP地址，速度快，但部分应用可能不兼容</li>
        <li><strong>Redir-Host</strong>：返回真实IP地址，兼容性好但稍慢</li>
        <li><strong>Mapping</strong>：使用映射表，介于两者之间</li>
      </ul>
    </el-alert>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import { getDNS, saveDNS } from '@/api/settings'

const loading = ref(false)
const saving = ref(false)

// 输入框
const nameserverInput = ref('')
const fallbackInput = ref('')
const defaultNameserverInput = ref('')
const fakeIPFilterInput = ref('')

const dnsForm = ref({
  enable: true,
  listen: '0.0.0.0:53',
  enhancedMode: 'fake-ip',
  nameserver: [],
  fallback: [],
  defaultNameserver: [],
  fakeIPFilter: []
})

// 常用DNS映射
const dnsMap = {
  '223.5.5.5': '阿里 DNS',
  '119.29.29.29': '腾讯 DNS',
  '180.76.76.76': '百度 DNS',
  '8.8.8.8': 'Google DNS',
  '1.1.1.1': 'Cloudflare DNS',
  '208.67.222.222': 'OpenDNS',
  '114.114.114.114': '114DNS'
}

const isCommonDNS = (ip) => dnsMap[ip] !== undefined

const getDNSName = (ip) => dnsMap[ip] || ip

// 主DNS操作
const addNameserver = () => {
  const dns = nameserverInput.value.trim()
  if (!dns) {
    ElMessage.warning('请输入DNS地址')
    return
  }
  if (dnsForm.value.nameserver.includes(dns)) {
    ElMessage.warning('该DNS已存在')
    return
  }
  dnsForm.value.nameserver.push(dns)
  nameserverInput.value = ''
}

const removeNameserver = (index) => {
  dnsForm.value.nameserver.splice(index, 1)
}

const quickAddNameserver = (dns) => {
  if (dnsForm.value.nameserver.includes(dns)) {
    ElMessage.warning('该DNS已存在')
    return
  }
  dnsForm.value.nameserver.push(dns)
}

// 备用DNS操作
const addFallback = () => {
  const dns = fallbackInput.value.trim()
  if (!dns) {
    ElMessage.warning('请输入DNS地址')
    return
  }
  if (dnsForm.value.fallback.includes(dns)) {
    ElMessage.warning('该DNS已存在')
    return
  }
  dnsForm.value.fallback.push(dns)
  fallbackInput.value = ''
}

const removeFallback = (index) => {
  dnsForm.value.fallback.splice(index, 1)
}

const quickAddFallback = (dns) => {
  if (dnsForm.value.fallback.includes(dns)) {
    ElMessage.warning('该DNS已存在')
    return
  }
  dnsForm.value.fallback.push(dns)
}

// 默认DNS操作
const addDefaultNameserver = () => {
  const dns = defaultNameserverInput.value.trim()
  if (!dns) {
    ElMessage.warning('请输入DNS地址')
    return
  }
  if (dnsForm.value.defaultNameserver.includes(dns)) {
    ElMessage.warning('该DNS已存在')
    return
  }
  dnsForm.value.defaultNameserver.push(dns)
  defaultNameserverInput.value = ''
}

const removeDefaultNameserver = (index) => {
  dnsForm.value.defaultNameserver.splice(index, 1)
}

const quickAddDefaultNameserver = (dns) => {
  if (dnsForm.value.defaultNameserver.includes(dns)) {
    ElMessage.warning('该DNS已存在')
    return
  }
  dnsForm.value.defaultNameserver.push(dns)
}

// FakeIP过滤操作
const addFakeIPFilter = () => {
  const filter = fakeIPFilterInput.value.trim()
  if (!filter) {
    ElMessage.warning('请输入域名')
    return
  }
  if (dnsForm.value.fakeIPFilter.includes(filter)) {
    ElMessage.warning('该过滤规则已存在')
    return
  }
  dnsForm.value.fakeIPFilter.push(filter)
  fakeIPFilterInput.value = ''
}

const removeFakeIPFilter = (index) => {
  dnsForm.value.fakeIPFilter.splice(index, 1)
}

const loadDNS = async () => {
  loading.value = true
  try {
    const data = await getDNS()
    dnsForm.value = {
      enable: data.enable ?? true,
      listen: data.listen || '0.0.0.0:53',
      enhancedMode: data.enhancedMode || 'fake-ip',
      nameserver: data.nameserver || [],
      fallback: data.fallback || [],
      defaultNameserver: data.defaultNameserver || [],
      fakeIPFilter: data.fakeIPFilter || []
    }
  } catch {
    ElMessage.error('加载DNS配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (dnsForm.value.nameserver.length === 0) {
    ElMessage.warning('请至少配置一个主DNS服务器')
    return
  }

  saving.value = true
  try {
    const data = {
      enable: dnsForm.value.enable,
      listen: dnsForm.value.listen,
      enhancedMode: dnsForm.value.enhancedMode,
      nameserver: dnsForm.value.nameserver,
      fallback: dnsForm.value.fallback,
      defaultNameserver: dnsForm.value.defaultNameserver,
      fakeIPFilter: dnsForm.value.fakeIPFilter
    }
    await saveDNS(data)
    ElMessage.success('保存成功')
  } catch {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadDNS()
})
</script>

<style scoped>
.settings-page {
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

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-divider__text) {
  font-weight: 500;
  color: #303133;
}

:deep(.el-alert) p {
  margin: 5px 0;
}

:deep(.el-alert) ul {
  margin: 0;
}

/* DNS输入容器样式 */
.dns-input-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.input-wrapper {
  display: flex;
  gap: 10px;
}

.quick-add {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.quick-label {
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
}

/* DNS标签样式 */
.dns-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px 0;
}

.dns-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  background: #ecf5ff;
  border: 1px solid #d9ecff;
  color: #409eff;
  font-size: 13px;
}

.dns-name {
  font-size: 11px;
  color: #909399;
  font-style: italic;
}

.empty-tip {
  color: #909399;
  font-size: 13px;
  padding: 8px 0;
}
</style>
