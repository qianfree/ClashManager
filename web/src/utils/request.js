import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useUserStore } from '@/stores/user'

const service = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// Request interceptor
service.interceptors.request.use(
  config => {
    // Add token if not a noAuth request
    if (!config.noAuth) {
      const userStore = useUserStore()
      const token = userStore.token
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Response interceptor
service.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      router.push('/login')
      ElMessage.error('登录已过期，请重新登录')
    } else {
      ElMessage.error(error.response?.data?.error || '请求失败')
    }
    return Promise.reject(error)
  }
)

export default service
