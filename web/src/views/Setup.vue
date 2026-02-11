<template>
  <div class="setup-page">
    <div class="setup-container">
      <div class="setup-box">
        <div class="setup-header">
          <el-icon size="48" color="#409eff"><svg viewBox="0 0 1024 1024" width="48" height="48"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64z m0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
          <h1>系统初始化</h1>
          <p>创建管理员账号</p>
        </div>

        <el-alert
          title="首次使用说明"
          type="info"
          :closable="false"
          show-icon
          style="margin-bottom: 20px;"
        >
          <p>欢迎使用 Clash Manager！</p>
          <p>这是您第一次使用系统，请创建管理员账号。</p>
        </el-alert>

        <el-form :model="setupForm" :rules="rules" ref="setupFormRef">
          <el-form-item prop="username">
            <el-input
              v-model="setupForm.username"
              placeholder="请输入管理员用户名"
              size="large"
              prefix-icon="User"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="setupForm.password"
              type="password"
              placeholder="请输入密码（至少6位）"
              size="large"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <el-input
              v-model="setupForm.confirmPassword"
              type="password"
              placeholder="请确认密码"
              size="large"
              prefix-icon="Lock"
              show-password
              @keyup.enter="handleSetup"
            />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              style="width: 100%"
              :loading="loading"
              @click="handleSetup"
            >
              创建管理员账号
            </el-button>
          </el-form-item>

          <el-form-item>
            <el-button
              size="large"
              style="width: 100%"
              @click="goToLogin"
            >
              返回登录
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { setup } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const setupFormRef = ref()
const loading = ref(false)

const setupForm = ref({
  username: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== setupForm.value.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleSetup = async () => {
  const valid = await setupFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await setup({
      username: setupForm.value.username,
      password: setupForm.value.password
    })
    ElMessage.success('管理员账号创建成功，请登录')
    router.push('/login')
  } catch {
    // Error already handled by request interceptor
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  router.push('/login')
}
</script>

<style scoped>
.setup-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.setup-container {
  width: 100%;
  max-width: 420px;
  padding: 20px;
}

.setup-box {
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.setup-header {
  text-align: center;
  margin-bottom: 30px;
}

.setup-header h1 {
  margin: 20px 0 10px;
  font-size: 28px;
  color: #303133;
}

.setup-header p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}

:deep(.el-alert__content) p {
  margin: 5px 0;
}
</style>
