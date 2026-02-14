import type {ICommonResponse} from '@/types/common.ts'

// 获取应用名称请求参数
export interface IGetAppNamesParam {}

// 获取应用名称相应
export type IGetAppNamesResp = ICommonResponse<{
  appNames: string[]
}>

// ---------------------------------------------------
export type IVerInfoItem = {
  appid: string
  appName: string
  version: string
  desc: string
  sign: string
  forceUpdate: boolean
  status: boolean
  patchUrl: string
  updateTime: string
}

// 获取应用名称相应
export type IGetVerInfoListResp = ICommonResponse<{
  infos: IVerInfoItem[]
}>

export interface IGetVerInfoListParams {
  appid?: string
}

export type IKAppIDVAppName = {
  appId: string
  appName: string
}

//--------------------------------------
export interface IVersionUpdateParams {
  appid: string
  version: string
  desc: string
  sign: string
  forceUpdate: boolean
  status: boolean
  patchUrl: string
}
