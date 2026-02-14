import request from '@/utils/request.ts'
import type {
  IAppResponse,
  IAppsParams,
  IAppsResponse,
  ICreateOrUpdateAppParams,
} from '@/types/app/app.ts'
import type {ICommonResponse} from '@/types/common.ts'

/**
 * 获取角色列表分页
 */
export const appPage = (params?: IAppsParams) => {
  return request.get<IAppsResponse>('/admin/app/list', {params})
}

/**
 * 创建角色
 */
export const createApp2 = (data: ICreateOrUpdateAppParams) => {
  return request.post<ICommonResponse<unknown>>('/admin/app/create', data)
}


export const updateApp = (data: ICreateOrUpdateAppParams) => {
  return request.post<ICommonResponse<unknown>>('/admin/app/create', data)
}

/**
 * 删除应用
 */
export const deleteApp = (names: string[]) => {
  return request.post<ICommonResponse<unknown>>('/admin/app/delete', {appids: names})
}

/**
 * 更新应用
 */
export const infoApp = (appid: string) => {
  return request.post<IAppResponse>('/admin/app/info', {appid: appid})
}
