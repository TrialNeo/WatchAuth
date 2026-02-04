import request from '@/utils/request'
import type { IUserDetailResponse } from '@/types/system/user'
import type {
  ILoginParams,
  ILoginResponse,
  IUserPermissionsResponse,
  ILoginLogParams,
} from '@/types/login'

export const login = (params: ILoginParams) => {
  return request.post<ILoginResponse>('admin/login', params)
}

/**
 * 获取用户权限（菜单权限和按钮权限）
 */
export const userPermissions = () => {
  return request.get<IUserPermissionsResponse>('/admin/permissions')
}

/**
 * 获取用户信息
 */
export const userInfoRequest = () => {
  return request.get<IUserDetailResponse>('/admin/info')
}
