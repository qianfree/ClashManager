import request from '@/utils/request'

export function getNodes() {
  return request({
    url: '/nodes',
    method: 'get'
  })
}

export function createNode(data) {
  return request({
    url: '/nodes',
    method: 'post',
    data
  })
}

export function updateNode(id, data) {
  return request({
    url: `/nodes/${id}`,
    method: 'put',
    data
  })
}

export function deleteNode(id) {
  return request({
    url: `/nodes/${id}`,
    method: 'delete'
  })
}

export function importNode(link) {
  return request({
    url: '/nodes/import',
    method: 'post',
    data: { link }
  })
}

export function exportNode(id) {
  return request({
    url: `/nodes/${id}/export`,
    method: 'get'
  })
}
