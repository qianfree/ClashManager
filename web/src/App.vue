<template>
  <!-- 登录和初始化页面不显示主布局 -->
  <router-view v-if="isAuthPage" />

  <!-- 其他页面显示完整的主布局 -->
  <div v-else class="app-container">
    <el-container style="height: 100vh;">
      <el-aside width="220px">
        <div class="logo">
          <el-icon><svg viewBox="0 0 1024 1024" width="32" height="32"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path><path fill="currentColor" d="M464 336a48 48 0 1 0 96 0 48 48 0 1 0-96 0zm72 112h-48c-4.4 0-8 3.6-8 8v272c0 4.4 3.6 8 8 8h48c4.4 0 8-3.6 8-8V456c0-4.4-3.6-8-8-8z"></path></svg></el-icon>
          <span>Clash Manager</span>
        </div>
        <el-menu
          :default-active="currentPath"
          router
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409eff"
        >
          <el-menu-item index="/nodes">
            <el-icon><svg viewBox="0 0 1024 1024" width="18" height="18"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>节点管理</span>
          </el-menu-item>
          <el-menu-item index="/rules">
            <el-icon><svg viewBox="0 0 1024 1024" width="18" height="18"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>规则管理</span>
          </el-menu-item>
          <el-menu-item index="/groups">
            <el-icon><svg viewBox="0 0 1024 1024" width="18" height="18"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>代理组管理</span>
          </el-menu-item>
          <el-menu-item index="/subscription">
            <el-icon><svg viewBox="0 0 1024 1024" width="18" height="18"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>订阅配置</span>
          </el-menu-item>
          <el-menu-item index="/subscription-logs">
            <el-icon><svg viewBox="0 0 1024 1024" width="18" height="18"><path fill="currentColor" d="M128 192h768v128H192v640h640v-64h64v128H128V192z"></path><path fill="currentColor" d="M384 384h384v64H384v320h320v-64h64v128H384V384z"></path></svg></el-icon>
            <span>订阅日志</span>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><svg viewBox="0 0 1024 1024" width="18" height="18"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372 166.6 372-372 372z"></path><path fill="currentColor" d="M464 336a48 48 0 1 0 96 0 48 48 0 1 0-96 0zm72 112h-48c-4.4 0-8 3.6-8 8v272c0 4.4 3.6 8 8 8h48c4.4 0 8-3.6 8-8V456c0-4.4-3.6-8-8-8z"></path></svg></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
        <div class="aside-footer">
          <el-divider />
          <div class="menu-item" @click="router.push('/password')">
            <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
            <span>修改密码</span>
          </div>
          <div class="menu-item logout-section" @click="handleLogout">
            <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M876 148H524V104c0-22.1-17.9-40-40-40H256c-22.1 0-40 17.9-40 40v44H148c-22.1 0-40 17.9-40 40v272c0 22.1 17.9 40 40 40h68v272c0 22.1 17.9 40 40 40h272c22.1 0 40-17.9 40-40V500h68c22.1 0 40-17.9 40-40V188c0-22.1-17.9-40-40-40z m-44 272H540V272h292v148z m-340 0h-52V148h52v272z m404 208H300V188h436v440z"></path></svg></el-icon>
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
                <el-icon size="20"><svg viewBox="0 0 1024 1024" width="20" height="20"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
                <span>{{ userStore.username }}</span>
                <el-icon class="el-icon--right"><svg viewBox="0 0 1024 1024" width="12" height="12"><path fill="currentColor" d="M831.872 340.864 512 652.672 192.128 340.864 340.864 192.128 652.672 512 652.672c0 159.168-129.024 288.192-288.192 288.192S-64.192 671.744 512 671.744c159.168 0 288.192-129.024 288.192-288.192S671.744 340.864 512 340.864zM640 448c0 70.656-57.344 128-128 128s-128-57.344-128-128 57.344-128 128-128 128 57.344 128 128z"></path></svg></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="password">
                    <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"></path></svg></el-icon>
                    修改密码
                  </el-dropdown-item>
                  <el-dropdown-item command="logout" divided>
                    <el-icon><svg viewBox="0 0 1024 1024" width="16" height="16"><path fill="currentColor" d="M876 148H524V104c0-22.1-17.9-40-40-40H256c-22.1 0-40 17.9-40 40v44H148c-22.1 0-40 17.9-40 40v272c0 22.1 17.9 40 40 40h68v272c0 22.1 17.9 40 40 40h272c22.1 0 40-17.9 40-40V500h68c22.1 0 40-17.9 40-40V188c0-22.1-17.9-40-40-40z m-44 272H540V272h292v148z m-340 0h-52V148h52v272z m404 208H300V188h436v440z"></path></svg></el-icon>
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
  background-color: #304156;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60px;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  gap: 10px;
  border-bottom: 1px solid #1f2d3d;
  flex-shrink: 0;
}

.logo .el-icon {
  color: #409eff;
}

.el-menu {
  border-right: none;
  flex: 1;
}

.aside-footer {
  border-top: 1px solid #1f2d3d;
  flex-shrink: 0;
}

:deep(.aside-footer .el-divider) {
  margin: 0;
  border-color: #1f2d3d;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height:  48px;
  padding-left: 20px;
  color: #bfcbd9;
  cursor: pointer;
  transition: all 0.3s;
}

.menu-item:hover {
  background-color: #263445;
}

.menu-item.logout-section:hover {
  background-color: #263445;
  color: #f56c6c;
}

.el-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
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

.user-info {
  display: flex;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-dropdown:hover {
  background-color: #f5f7fa;
}

.el-main {
  background: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

:deep(.el-menu-item) {
  display: flex;
  align-items: center;
}

:deep(.el-menu-item .el-icon) {
  margin-right: 8px;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
