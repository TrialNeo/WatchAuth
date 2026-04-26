import request from '@/utils/request.ts';
import type { ICommonResponse } from '@/types/common.ts';

export interface OverviewStats {
  totalMachines: number
  totalApps: number
  onlineMachines: number
}

export interface MachineStats {
  todayNew: number
  monthNew: number
  active7d: number
  active30d: number
}

export interface DistributionItem {
  name: string
  count: number
}

export interface Distributions {
  platform: DistributionItem[]
  arch: DistributionItem[]
}

export interface LicenseStats {
  valid: number
  revoked: number
  expired: number
  todayNew: number
  monthNew: number
}

export interface AppRankingItem {
  appName: string
  appId: string
  machines: number
  licenses: number
  online: number
}

export interface DailyTrendItem {
  date: string
  count: number
}

export interface StatisticsData {
  overview: OverviewStats
  machineStats: MachineStats
  distributions: Distributions
  licenseStats: LicenseStats
  appRankings: AppRankingItem[]
  machineTrend: DailyTrendItem[]
  licenseTrend: DailyTrendItem[]
}

export const getStatistics = () => {
  return request.post<ICommonResponse<StatisticsData>>('/admin/statistics')
}
