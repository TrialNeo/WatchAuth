import type {ICommonResponse} from '@/types/common.ts'

export type IUserItem = {
  userid: number
  userName: string
  email: string
  account: string
  balance: string
  spend: string
  createTime: string
}
//IUserItemResp 用户列表
export type IUserItemResp = ICommonResponse<{
  IUserItems: IUserItem[]
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
