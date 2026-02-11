import request from '@/utils/request'

export function getRules(params) {
  return request({
    url: '/rules',
    method: 'get',
    params
  })
}

export function createRule(data) {
  return request({
    url: '/rules',
    method: 'post',
    data
  })
}

export function importRules(content) {
  return request({
    url: '/rules/import',
    method: 'post',
    data: { content }
  })
}

export function updateRule(id, data) {
  return request({
    url: `/rules/${id}`,
    method: 'put',
    data
  })
}

export function deleteRule(id) {
  return request({
    url: `/rules/${id}`,
    method: 'delete'
  })
}
