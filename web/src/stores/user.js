import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const username = ref(localStorage.getItem('username') || '')

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUsername = (newUsername) => {
    username.value = newUsername
    localStorage.setItem('username', newUsername)
  }

  const logout = () => {
    token.value = ''
    username.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('username')
  }

  const isLoggedIn = () => {
    return !!token.value
  }

  return {
    token,
    username,
    setToken,
    setUsername,
    logout,
    isLoggedIn
  }
})
