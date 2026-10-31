import service from '@/utils/request'

export const getAppDataList = (resource, params) => service({ url: `/appData/${resource}/list`, method: 'get', params })
export const getAppDataDetail = (resource, params) => service({ url: `/appData/${resource}/detail`, method: 'get', params })
