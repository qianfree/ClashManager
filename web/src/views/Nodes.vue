<template>
  <div class="nodes-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>节点列表</span>
          </div>
          <div class="header-right">
            <el-button type="primary" :icon="Upload" @click="showImportDialog">导入节点</el-button>
            <el-button type="success" :icon="Plus" @click="showCreateDialog">新增节点</el-button>
          </div>
        </div>
      </template>

      <el-table :data="nodes" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" min-width="60" />
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column prop="type" label="类型" min-width="100">
          <template #default="{ row }">
            <el-tag>{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="server" label="服务器" min-width="180" show-overflow-tooltip />
        <el-table-column prop="port" label="端口" min-width="80" />
        <el-table-column prop="network" label="传输" min-width="80" />
        <el-table-column label="来源" min-width="120">
          <template #default="{ row }">
            <el-tag v-if="row.source === 'subscription'" type="warning" size="small">
              {{ row.sourceName || '订阅' }}
            </el-tag>
            <el-tag v-else type="info" size="small">手动</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="80">
          <template #default="{ row }">
            <el-tag :type="row.tls ? 'success' : 'info'" size="small">
              {{ row.tls ? 'TLS' : '普通' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="success" link @click="handleExport(row)">导出</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 导入节点对话框 -->
    <el-dialog v-model="importDialogVisible" title="导入节点" width="500px">
      <el-form :model="importForm">
        <el-form-item label="分享链接">
          <el-input
            v-model="importForm.link"
            type="textarea"
            :rows="4"
            placeholder="请粘贴节点分享链接 (ss://, vmess://, trojan://, hysteria2:// 等)"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleImport">导入</el-button>
      </template>
    </el-dialog>

    <!-- 新增/编辑节点对话框 -->
    <el-dialog v-model="formDialogVisible" :title="isEdit ? '编辑节点' : '新增节点'" width="700px">
      <el-tabs v-model="activeTab" type="border-card" @tab-change="handleTabChange">
        <!-- Shadowsocks -->
        <el-tab-pane label="Shadowsocks" name="ss">
          <template #label>
            <span style="display: flex; align-items: center; gap: 5px;">
              <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64z m0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
              Shadowsocks
            </span>
          </template>
          <el-form :model="nodeForm" label-width="110px" class="node-form">
            <el-form-item label="节点名称">
              <el-input v-model="nodeForm.Name" placeholder="请输入节点名称" />
            </el-form-item>
            <el-form-item label="服务器地址">
              <el-input v-model="nodeForm.Server" placeholder="请输入服务器地址" />
            </el-form-item>
            <el-form-item label="端口">
              <el-input-number v-model="nodeForm.Port" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
            <el-form-item label="加密方式">
              <el-select v-model="nodeForm.Cipher" placeholder="请选择加密方式" style="width: 100%">
                <el-option label="aes-128-gcm" value="aes-128-gcm" />
                <el-option label="aes-192-gcm" value="aes-192-gcm" />
                <el-option label="aes-256-gcm" value="aes-256-gcm" />
                <el-option label="aes-128-cfb" value="aes-128-cfb" />
                <el-option label="aes-192-cfb" value="aes-192-cfb" />
                <el-option label="aes-256-cfb" value="aes-256-cfb" />
                <el-option label="aes-128-ctr" value="aes-128-ctr" />
                <el-option label="aes-192-ctr" value="aes-192-ctr" />
                <el-option label="aes-256-ctr" value="aes-256-ctr" />
                <el-option label="rc4-md5" value="rc4-md5" />
                <el-option label="chacha20-ietf" value="chacha20-ietf" />
                <el-option label="xchacha20" value="xchacha20" />
                <el-option label="chacha20-ietf-poly1305" value="chacha20-ietf-poly1305" />
                <el-option label="xchacha20-ietf-poly1305" value="xchacha20-ietf-poly1305" />
              </el-select>
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="nodeForm.Password" placeholder="请输入密码" show-password />
            </el-form-item>
            <el-form-item label="UDP转发">
              <el-switch v-model="nodeForm.UDP" />
              <span style="margin-left: 10px; color: #909399; font-size: 12px;">启用UDP转发</span>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- VMess -->
        <el-tab-pane label="VMess" name="vmess">
          <template #label>
            <span style="display: flex; align-items: center; gap: 5px;">
              <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64z m0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
              VMess
            </span>
          </template>
          <el-form :model="nodeForm" label-width="110px" class="node-form">
            <el-form-item label="节点名称">
              <el-input v-model="nodeForm.Name" placeholder="请输入节点名称" />
            </el-form-item>
            <el-form-item label="服务器地址">
              <el-input v-model="nodeForm.Server" placeholder="请输入服务器地址" />
            </el-form-item>
            <el-form-item label="端口">
              <el-input-number v-model="nodeForm.Port" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
            <el-form-item label="UUID">
              <el-input v-model="nodeForm.UUID" placeholder="请输入UUID" />
            </el-form-item>
            <el-form-item label="加密方式">
              <el-select v-model="nodeForm.Cipher" placeholder="请选择加密方式" style="width: 100%">
                <el-option label="auto" value="auto" />
                <el-option label="aes-128-gcm" value="aes-128-gcm" />
                <el-option label="chacha20-poly1305" value="chacha20-poly1305" />
                <el-option label="none" value="none" />
              </el-select>
            </el-form-item>
            <el-form-item label="传输协议">
              <el-select v-model="nodeForm.Network" placeholder="请选择传输协议" style="width: 100%">
                <el-option label="TCP" value="" />
                <el-option label="WebSocket" value="ws" />
                <el-option label="gRPC" value="grpc" />
              </el-select>
            </el-form-item>
            <el-form-item label="路径" v-if="nodeForm.Network === 'ws' || nodeForm.Network === 'grpc'">
              <el-input v-model="nodeForm.Path" :placeholder="nodeForm.Network === 'ws' ? '请输入WebSocket路径，如: /ws' : '请输入gRPC服务名'" />
            </el-form-item>
            <el-form-item label="Host/SNI" v-if="nodeForm.Network === 'ws' || nodeForm.Network === 'grpc' || nodeForm.TLS">
              <el-input v-model="nodeForm.Host" placeholder="请输入Host或SNI" />
            </el-form-item>
            <el-form-item label="TLS加密">
              <el-switch v-model="nodeForm.TLS" />
            </el-form-item>
            <el-form-item label="跳过证书验证" v-if="nodeForm.TLS">
              <el-switch v-model="nodeForm.SkipCert" />
            </el-form-item>
            <el-form-item label="UDP转发">
              <el-switch v-model="nodeForm.UDP" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- VLESS -->
        <el-tab-pane label="VLESS" name="vless">
          <template #label>
            <span style="display: flex; align-items: center; gap: 5px;">
              <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64z m0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
              VLESS
            </span>
          </template>
          <el-form :model="nodeForm" label-width="110px" class="node-form">
            <el-form-item label="节点名称">
              <el-input v-model="nodeForm.Name" placeholder="请输入节点名称" />
            </el-form-item>
            <el-form-item label="服务器地址">
              <el-input v-model="nodeForm.Server" placeholder="请输入服务器地址" />
            </el-form-item>
            <el-form-item label="端口">
              <el-input-number v-model="nodeForm.Port" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
            <el-form-item label="UUID">
              <el-input v-model="nodeForm.UUID" placeholder="请输入UUID" />
            </el-form-item>
            <el-form-item label="传输协议">
              <el-select v-model="nodeForm.Network" placeholder="请选择传输协议" style="width: 100%">
                <el-option label="TCP" value="" />
                <el-option label="WebSocket" value="ws" />
                <el-option label="gRPC" value="grpc" />
              </el-select>
            </el-form-item>
            <el-form-item label="路径" v-if="nodeForm.Network === 'ws' || nodeForm.Network === 'grpc'">
              <el-input v-model="nodeForm.Path" :placeholder="nodeForm.Network === 'ws' ? '请输入WebSocket路径，如: /ws' : '请输入gRPC服务名'" />
            </el-form-item>
            <el-form-item label="Host/SNI" v-if="nodeForm.Network === 'ws' || nodeForm.Network === 'grpc' || nodeForm.TLS">
              <el-input v-model="nodeForm.Host" placeholder="请输入Host或SNI" />
            </el-form-item>
            <el-form-item label="TLS加密">
              <el-switch v-model="nodeForm.TLS" />
            </el-form-item>
            <el-form-item label="跳过证书验证" v-if="nodeForm.TLS">
              <el-switch v-model="nodeForm.SkipCert" />
            </el-form-item>
            <el-form-item label="UDP转发">
              <el-switch v-model="nodeForm.UDP" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- Trojan -->
        <el-tab-pane label="Trojan" name="trojan">
          <template #label>
            <span style="display: flex; align-items: center; gap: 5px;">
              <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64z m0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
              Trojan
            </span>
          </template>
          <el-form :model="nodeForm" label-width="110px" class="node-form">
            <el-form-item label="节点名称">
              <el-input v-model="nodeForm.Name" placeholder="请输入节点名称" />
            </el-form-item>
            <el-form-item label="服务器地址">
              <el-input v-model="nodeForm.Server" placeholder="请输入服务器地址" />
            </el-form-item>
            <el-form-item label="端口">
              <el-input-number v-model="nodeForm.Port" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="nodeForm.Password" placeholder="请输入密码" show-password />
            </el-form-item>
            <el-form-item label="传输协议">
              <el-select v-model="nodeForm.Network" placeholder="请选择传输协议" style="width: 100%">
                <el-option label="TCP" value="" />
                <el-option label="WebSocket" value="ws" />
                <el-option label="gRPC" value="grpc" />
              </el-select>
            </el-form-item>
            <el-form-item label="路径" v-if="nodeForm.Network === 'ws' || nodeForm.Network === 'grpc'">
              <el-input v-model="nodeForm.Path" :placeholder="nodeForm.Network === 'ws' ? '请输入WebSocket路径，如: /ws' : '请输入gRPC服务名'" />
            </el-form-item>
            <el-form-item label="Host/SNI" v-if="nodeForm.Network === 'ws' || nodeForm.Network === 'grpc'">
              <el-input v-model="nodeForm.Host" placeholder="请输入Host或SNI" />
            </el-form-item>
            <el-form-item label="TLS加密">
              <el-switch v-model="nodeForm.TLS" />
              <span style="margin-left: 10px; color: #909399; font-size: 12px;">Trojan必须启用TLS</span>
            </el-form-item>
            <el-form-item label="跳过证书验证" v-if="nodeForm.TLS">
              <el-switch v-model="nodeForm.SkipCert" />
            </el-form-item>
            <el-form-item label="UDP转发">
              <el-switch v-model="nodeForm.UDP" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- Hysteria2 -->
        <el-tab-pane label="Hysteria2" name="hysteria2">
          <template #label>
            <span style="display: flex; align-items: center; gap: 5px;">
              <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64z m0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
              Hysteria2
            </span>
          </template>
          <el-form :model="nodeForm" label-width="110px" class="node-form">
            <el-form-item label="节点名称">
              <el-input v-model="nodeForm.Name" placeholder="请输入节点名称" />
            </el-form-item>
            <el-form-item label="服务器地址">
              <el-input v-model="nodeForm.Server" placeholder="请输入服务器地址" />
            </el-form-item>
            <el-form-item label="端口">
              <el-input-number v-model="nodeForm.Port" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="nodeForm.Password" placeholder="请输入密码" show-password />
            </el-form-item>
            <el-form-item label="SNI">
              <el-input v-model="nodeForm.Host" placeholder="请输入SNI" />
            </el-form-item>
            <el-form-item label="TLS加密">
              <el-switch v-model="nodeForm.TLS" />
              <span style="margin-left: 10px; color: #909399; font-size: 12px;">Hysteria2必须启用TLS</span>
            </el-form-item>
            <el-form-item label="跳过证书验证" v-if="nodeForm.TLS">
              <el-switch v-model="nodeForm.SkipCert" />
            </el-form-item>
            <el-form-item label="UDP转发">
              <el-switch v-model="nodeForm.UDP" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
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
import { Plus, Upload } from '@element-plus/icons-vue'
import { getNodes, createNode, updateNode, deleteNode, importNode, exportNode } from '@/api/nodes'

const nodes = ref([])
const importDialogVisible = ref(false)
const formDialogVisible = ref(false)
const activeTab = ref('vmess')
const isEdit = ref(false)
const importForm = ref({ link: '' })

const nodeForm = ref({
  Name: '',
  Type: 'vmess',
  Server: '',
  Port: 443,
  UUID: '',
  Password: '',
  Cipher: '',
  Network: '',
  Path: '',
  Host: '',
  TLS: true,
  SkipCert: false,
  UDP: true
})

const getTypeLabel = (type) => {
  const labels = {
    ss: 'Shadowsocks',
    vmess: 'VMess',
    vless: 'VLESS',
    trojan: 'Trojan',
    hysteria2: 'Hysteria2'
  }
  return labels[type] || type
}

const loadNodes = async () => {
  nodes.value = await getNodes()
}

const showImportDialog = () => {
  importForm.value.link = ''
  importDialogVisible.value = true
}

const handleImport = async () => {
  if (!importForm.value.link.trim()) {
    ElMessage.warning('请输入节点链接')
    return
  }
  await importNode(importForm.value.link)
  ElMessage.success('导入成功')
  importDialogVisible.value = false
  loadNodes()
}

const showCreateDialog = () => {
  isEdit.value = false
  activeTab.value = 'vmess'
  resetForm()
  formDialogVisible.value = true
}

const resetForm = () => {
  nodeForm.value = {
    Name: '',
    Type: activeTab.value,
    Server: '',
    Port: 443,
    UUID: '',
    Password: '',
    Cipher: activeTab.value === 'ss' ? 'aes-256-gcm' : 'auto',
    Network: '',
    Path: '',
    Host: '',
    TLS: true,
    SkipCert: false,
    UDP: true
  }
}

const handleTabChange = (tabName) => {
  nodeForm.value.Type = tabName
  // 切换tab时重置表单，保留已填写的名称和服务器
  const keepFields = { Name: nodeForm.value.Name, Server: nodeForm.value.Server, Port: nodeForm.value.Port }
  resetForm()
  Object.assign(nodeForm.value, keepFields)
}

const handleEdit = (row) => {
  isEdit.value = true
  activeTab.value = row.type
  // 将后端返回的蛇形命名字段转换为表单的驼峰命名
  nodeForm.value = {
    ID: row.id,
    Name: row.name,
    Type: row.type,
    Server: row.server,
    Port: row.port,
    UUID: row.uuid || '',
    Password: row.password || '',
    Cipher: row.cipher || '',
    Network: row.network || '',
    Path: row.path || '',
    Host: row.host || '',
    TLS: row.tls || false,
    SkipCert: row.skipCert || false,
    UDP: row.udp || false,
    ALPN: row.alpn || '',
    ExtraConfig: row.extraConfig || ''
  }
  formDialogVisible.value = true
}

const handleSave = async () => {
  if (!nodeForm.value.Name || !nodeForm.value.Server) {
    ElMessage.warning('请填写节点名称和服务器地址')
    return
  }

  // 根据类型验证必填字段
  const type = nodeForm.value.Type
  if (type === 'vmess' || type === 'vless') {
    if (!nodeForm.value.UUID) {
      ElMessage.warning('请填写UUID')
      return
    }
  } else if (type === 'ss' || type === 'trojan' || type === 'hysteria2') {
    if (!nodeForm.value.Password) {
      ElMessage.warning('请填写密码')
      return
    }
  }

  // 转换为后端期望的 snake_case 格式
  const data = {
    name: nodeForm.value.Name,
    type: nodeForm.value.Type,
    server: nodeForm.value.Server,
    port: nodeForm.value.Port,
    uuid: nodeForm.value.UUID || '',
    password: nodeForm.value.Password || '',
    cipher: nodeForm.value.Cipher || '',
    network: nodeForm.value.Network || '',
    path: nodeForm.value.Path || '',
    host: nodeForm.value.Host || '',
    tls: nodeForm.value.TLS,
    skip_cert: nodeForm.value.SkipCert,
    udp: nodeForm.value.UDP,
    alpn: nodeForm.value.ALPN || '',
    extra_config: nodeForm.value.ExtraConfig || ''
  }

  if (isEdit.value) {
    await updateNode(nodeForm.value.ID, data)
    ElMessage.success('更新成功')
  } else {
    await createNode(data)
    ElMessage.success('创建成功')
  }
  formDialogVisible.value = false
  loadNodes()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该节点吗？', '提示', { type: 'warning' })
  await deleteNode(row.id)
  ElMessage.success('删除成功')
  loadNodes()
}

const handleExport = async (row) => {
  try {
    const result = await exportNode(row.id)
    const link = result.link

    // 复制到剪贴板
    await navigator.clipboard.writeText(link)
    ElMessage.success('节点分享链接已复制到剪贴板')

    // 可选：同时显示链接让用户确认
    console.log('导出的节点链接:', link)
  } catch (error) {
    // 如果剪贴板复制失败，显示链接让用户手动复制
    try {
      const result = await exportNode(row.id)
      const link = result.link

      ElMessageBox.alert(
        link,
        '节点分享链接',
        {
          confirmButtonText: '确定',
          type: 'success',
          inputType: 'textarea',
          showInput: false,
          dangerouslyUseHTMLString: false,
          message: link,
          customClass: 'export-link-dialog'
        }
      ).catch(() => {})
    } catch (err) {
      console.error('导出失败:', err)
      ElMessage.error('导出失败：' + (err.response?.data?.error || err.message))
    }
  }
}

onMounted(() => {
  loadNodes()
})
</script>

<style scoped>
.nodes-page {
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

.header-right {
  display: flex;
  gap: 10px;
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

.node-form {
  padding: 20px;
  max-height: 500px;
  overflow-y: auto;
}

:deep(.el-tabs__header) {
  margin: 0;
}

:deep(.el-tabs__content) {
  padding: 0;
}

:deep(.el-tab-pane) {
  background: #fff;
}

/* 导出链接对话框样式 */
:deep(.export-link-dialog .el-message-box__content) {
  padding: 10px 20px;
}

:deep(.export-link-dialog .el-message-box__message) {
  word-break: break-all;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 12px;
  line-height: 1.5;
  max-height: 200px;
  overflow-y: auto;
}
</style>
