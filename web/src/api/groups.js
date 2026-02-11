import request from '@/utils/request'

export function getGroups() {
  return request({
    url: '/groups',
    method: 'get'
  })
}

export function createGroup(data) {
  return request({
    url: '/groups',
    method: 'post',
    data
  })
}

export function updateGroup(id, data) {
  return request({
    url: `/groups/${id}`,
    method: 'put',
    data
  })
}

export function deleteGroup(id) {
  return request({
    url: `/groups/${id}`,
    method: 'delete'
  })
}
