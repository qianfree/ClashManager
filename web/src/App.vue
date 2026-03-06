<template>
  <!-- 登录和初始化页面不显示主布局 -->
  <router-view v-if="isAuthPage" />

  <!-- 其他页面显示完整的主布局 -->
  <div v-else class="app-container">
    <el-container style="height: 100vh;">
      <el-aside width="240px">
        <div class="logo">
          <div class="logo-icon">
            <svg viewBox="0 0 1024 1024" width="28" height="28">
              <path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64z" opacity="0.2"/>
              <path fill="currentColor" d="M765.9 186.2c-119.4-114.7-308.9-111.7-424.4 6.9S227.5 496.7 346.9 611.4c119.4 114.7 308.9 111.7 424.4-6.9s113.9-303.6-5.5-418.3zM393.3 657.6c-91.8-88.1-94.4-233.6-5.8-324.9s235.2-95.3 327 7.1c91.8 102.4 89.2 247.9-2.6 336.2-88.5 85.2-230.4 83.5-318.6-18.4z"/>
              <path fill="currentColor" d="M512 320c-17.7 0-32 14.3-32 32v160c0 17.7 14.3 32 32 32s32-14.3 32-32V352c0-17.7-14.3-32-32-32zm0 280c-22.1 0-40 17.9-40 40s17.9 40 40 40 40-17.9 40-40-17.9-40-40-40z"/>
            </svg>
          </div>
          <span class="logo-text">Clash Manager</span>
        </div>
        <el-menu
          :default-active="currentPath"
          router
          class="sidebar-menu"
        >
          <el-menu-item index="/nodes">
            <el-icon><Connection /></el-icon>
            <span>节点管理</span>
          </el-menu-item>
          <el-menu-item index="/rules">
            <el-icon><DocumentCopy /></el-icon>
            <span>规则管理</span>
          </el-menu-item>
          <el-menu-item index="/groups">
            <el-icon><Grid /></el-icon>
            <span>代理组管理</span>
          </el-menu-item>
          <el-menu-item index="/subscription">
            <el-icon><Link /></el-icon>
            <span>订阅配置</span>
          </el-menu-item>
          <el-menu-item index="/subscription-logs">
            <el-icon><Document /></el-icon>
            <span>订阅日志</span>
          </el-menu-item>
          <el-menu-item index="/sources">
            <el-icon><Download /></el-icon>
            <span>订阅源管理</span>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
        <div class="aside-footer">
          <el-divider />
          <div class="menu-item" @click="router.push('/password')">
            <el-icon><Lock /></el-icon>
            <span>修改密码</span>
          </div>
          <div class="menu-item logout-section" @click="handleLogout">
            <el-icon><SwitchButton /></el-icon>
            <span>退出登录</span>
          </div>
        </div>
      </el-aside>
      <el-container>
        <el-header>
          <div class="header-content">
            <div class="breadcrumb">
              <el-breadcrumb separator="/">
                <el-breadcrumb-item>{{ currentPageTitle }}</el-breadcrumb-item>
              </el-breadcrumb>
            </div>
          </div>
          <div class="user-info">
            <el-dropdown @command="handleCommand">
              <div class="user-dropdown">
                <div class="user-avatar">
                  <el-icon><User /></el-icon>
                </div>
                <span class="username">{{ userStore.username }}</span>
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="password">
                    <el-icon><Lock /></el-icon>
                    修改密码
                  </el-dropdown-item>
                  <el-dropdown-item command="logout" divided>
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Connection,
  DocumentCopy,
  Grid,
  Link,
  Document,
  Setting,
  Lock,
  SwitchButton,
  User,
  ArrowDown,
  Download
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const currentPath = computed(() => route.path)

// 判断是否是登录或初始化页面
const isAuthPage = computed(() => {
  return route.path === '/login' || route.path === '/setup'
})

const pageTitleMap = {
  '/nodes': '节点管理',
  '/rules': '规则管理',
  '/groups': '代理组管理',
  '/subscription': '订阅配置',
  '/subscription-logs': '订阅日志',
  '/sources': '订阅源管理',
  '/settings': '系统设置',
  '/password': '修改密码'
}

const currentPageTitle = computed(() => {
  return pageTitleMap[route.path] || 'Clash配置管理'
})

const handleCommand = (command) => {
  if (command === 'logout') {
    handleLogout()
  } else if (command === 'password') {
    router.push('/password')
  }
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }).catch(() => {
    // User cancelled
  })
}
</script>

<style scoped>
.app-container {
  width: 100%;
}

.el-aside {
  background: linear-gradient(180deg, #1a1f3a 0%, #131729 100%);
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 64px;
  gap: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
  padding: 0 16px;
}

.logo-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.sidebar-menu {
  border-right: none;
  flex: 1;
  padding: 12px 8px;
  background: transparent;
}

:deep(.sidebar-menu.el-menu) {
  background: transparent;
}

:deep(.sidebar-menu .el-menu-item) {
  color: rgba(255, 255, 255, 0.7);
  border-radius: 8px;
  margin-bottom: 4px;
  height: 44px;
  line-height: 44px;
  transition: all 0.3s ease;
}

:deep(.sidebar-menu .el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

:deep(.sidebar-menu .el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(102, 126, 234, 0.3) 0%, rgba(118, 75, 162, 0.3) 100%);
  color: #fff;
  position: relative;
}

:deep(.sidebar-menu .el-menu-item.is-active::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  border-radius: 0 3px 3px 0;
}

:deep(.sidebar-menu .el-menu-item .el-icon) {
  margin-right: 8px;
  font-size: 18px;
}

.aside-footer {
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
  padding: 12px 8px;
}

:deep(.aside-footer .el-divider) {
  margin: 0 0 12px 0;
  border-color: rgba(255, 255, 255, 0.08);
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 44px;
  padding: 0 16px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 8px;
  font-size: 14px;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.menu-item.logout-section:hover {
  background: rgba(245, 108, 108, 0.15);
  color: #f56c6c;
}

.menu-item .el-icon {
  font-size: 16px;
}

.el-header {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.breadcrumb {
  flex: 1;
}

:deep(.breadcrumb .el-breadcrumb__item) {
  font-size: 15px;
}

:deep(.breadcrumb .el-breadcrumb__inner) {
  color: #303133;
  font-weight: 500;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.user-dropdown:hover {
  background: #f5f7fa;
}

.user-avatar {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
}

.username {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.el-main {
  background: #f5f7fa;
  padding: 24px;
  overflow-y: auto;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-dropdown-menu__item .el-icon) {
  font-size: 16px;
}
</style>
