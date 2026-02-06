// 角色管理类型文件
import type {ICommonResponse} from '@/types/common.ts'

// 角色列表项
export interface IAppItem {
  appId: string
  appName: string
  version: string
  feeType: number
  fee: number
  secretKeys: string
  description: string
  enctype: number
  status: number
  createdAt?: string
  updatedAt?: string
}

// 应用列表查询参数
export interface IAppsParams {
  page: number
  pageSize: number
  name?: string
  code?: string
  status?: 'active' | 'inactive'
  sortOrder?: 'asc' | 'desc'
}

// 创建/更新应用参数
export interface ICreateOrUpdateAppParams {
  appid:string
  appName: string
  description: string
  status: number
  feeType: number
  fee: number
  encType: number
}


// 应用信息响应
export type IAppsResponse = ICommonResponse<{
  apps: IAppItem[]
  total: number
  page: number
  pageSize: number
}>


// 角色列表响应
export type IAppResponse = ICommonResponse<{
  app: IAppItem
}>
