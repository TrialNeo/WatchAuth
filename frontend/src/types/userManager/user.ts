import type {ICommonResponse} from '@/types/common.ts'

export type IUserItem = {
  userID: number
  userName: string
  email: string
  privilege: number
  account: string
  balance: string
  spend: string
  createTime: string
}
//IUserItemResp 用户列表
export type IUserItemResp = ICommonResponse<{
  list: IUserItem[]
}>
//IUserBanParams 用户封禁
export interface IUserBanParams {
  userid: number
}
//IUpdateUserProfileParams 修改用户资料
export interface IUpdateUserProfileParams {
  userid: number
  userName: string
  password: string
  email: string
}

//ITopUpParams 用户充值
export interface ITopUpParams {
  userid: number
  amt: string
}

// 机器信息类型
export type IMachineInfo = {
  platform: string
  arch: string
  deviceId: string
  machineName: string
  cpu: string
  gpu: string
  ram: string
}

// 使用的应用类型
export type IUsedApp = {
  appId: number
  online: boolean
  loginIp?: string
  lastOnlineAt: string
  lastHeartbeatAt: string
  lastOfflineAt: string | null
}

// 机器项类型
export type IMachineItem = {
  machineId: number
  belong: number
  isBan?: boolean
  machine: IMachineInfo
  usedApps: IUsedApp[]
}

// 日志记录类型
export type ILogItem = {
  time: string
  type: 'info' | 'debug' | 'warning' | 'error'
  appId: number
  appVersionId: string
  trainId: number
  spanId: number
  module: string
  funcName: string
  msg: string
}

// 机器列表响应类型 - data 直接是数组
export type IMachineListResp = ICommonResponse<IMachineItem[]>

// 日志列表响应类型 - data 直接是数组
export type ILogListResp = ICommonResponse<ILogItem[]>
