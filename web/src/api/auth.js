import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data,
    noAuth: true
  })
}

export function setup(data) {
  return request({
    url: '/auth/setup',
    method: 'post',
    data,
    noAuth: true
  })
}

export function changePassword(data) {
  return request({
    url: '/auth/password',
    method: 'post',
    data
  })
}

export function checkSetup() {
  return request({
    url: '/auth/setup',
    method: 'get',
    noAuth: true
  })
}
