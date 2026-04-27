import request from '@/utils/request.ts'
import type { ICommonResponse } from '@/types/common.ts'

export interface ConfigItem {
  key: string
  value: string
}

export const getSystemConfigs = () => {
  return request.get<ICommonResponse<{ configs: ConfigItem[] }>>('/admin/config/list')
}

export const updateSystemConfig = (key: string, value: string) => {
  return request.post<ICommonResponse<null>>('/admin/config/update', { key, value })
}
