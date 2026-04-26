import request from '@/utils/request.ts'
import type { ICommonResponse } from '@/types/common.ts'

export interface LicenseItem {
  id: number
  licenseKey: string
  appId: string
  appName: string
  machineId: number
  userId: number
  username: string
  expireAt: string
  issuedAt: string
  status: number
}

export const getLicenseList = () => {
  return request.get<ICommonResponse<LicenseItem[]>>('/admin/license/list')
}
