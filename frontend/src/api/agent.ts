import request from '@/utils/request.ts'
import type { ICommonResponse } from '@/types/common.ts'

export interface AgentItem {
  id: number
  name: string
  contact: string
  parentId: number
  level: number
  discount: number
  balance: number
  status: number
  createdAt: string
  updatedAt: string
  children?: AgentItem[]
}

export const getAgentList = () => {
  return request.get<ICommonResponse<{ agents: AgentItem[] }>>('/admin/agent/list')
}

export const getAgent = (id: number) => {
  return request.get<ICommonResponse<{ agent: AgentItem }>>(`/admin/agent/${id}`)
}

export const createAgent = (data: { name: string; contact: string; parentId: number; discount: number }) => {
  return request.post<ICommonResponse<{ id: number }>>('/admin/agent/create', data)
}

export const updateAgent = (data: { id: number; name: string; contact: string; discount: number; status: number }) => {
  return request.post<ICommonResponse<null>>('/admin/agent/update', data)
}

export const deleteAgent = (id: number) => {
  return request.post<ICommonResponse<null>>('/admin/agent/delete', { id })
}
