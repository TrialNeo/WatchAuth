import request from '@/utils/request.ts'
import type { ICommonResponse } from '@/types/common.ts'

export interface OperationLogItem {
  id: number
  adminId: number
  adminName: string
  action: string
  target: string
  targetId: string
  detail: string
  ip: string
  createdAt: string
}

export interface OperationLogListResp {
  list: OperationLogItem[]
  total: number
}

export const getOperationLogs = (params: { page?: number; pageSize?: number; action?: string; target?: string }) => {
  return request.get<ICommonResponse<OperationLogListResp>>('/admin/log/list', { params })
}

export const getRecentLogs = () => {
  return request.get<ICommonResponse<{ list: OperationLogItem[] }>>('/admin/log/recent')
}

export interface TodayStats {
  newMachines: number
  newLicenses: number
  newUsers: number
  onlineCount: number
  totalAgents: number
}

export const getTodayStats = () => {
  return request.get<ICommonResponse<TodayStats>>('/admin/stats/today')
}
