import request from '@/utils/request'

export function getToken() {
  return request({
    url: '/subscription/token',
    method: 'get'
  })
}

export function refreshToken() {
  return request({
    url: '/subscription/token/refresh',
    method: 'post'
  })
}

export function getSubscriptionURL() {
  return request({
    url: '/subscription/url',
    method: 'get'
  })
}

export function getSubscriptionLogs(params) {
  return request({
    url: '/subscription/logs',
    method: 'get',
    params
  })
}

export function getSubscriptionStats(params) {
  return request({
    url: '/subscription/stats',
    method: 'get',
    params
  })
}

export function deleteOldLogs(days) {
  return request({
    url: '/subscription/logs/old',
    method: 'delete',
    params: { days }
  })
}
