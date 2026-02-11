import request from '@/utils/request'

export function getDNS() {
  return request({
    url: '/settings/dns',
    method: 'get'
  })
}

export function saveDNS(data) {
  return request({
    url: '/settings/dns',
    method: 'post',
    data
  })
}
