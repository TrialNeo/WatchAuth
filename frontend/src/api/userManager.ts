import request from '@/utils/request.ts'
import type {IUpdateUserProfileParams} from '@/types/system/user.ts'
import type {ICommonResponse} from '@/types/common.ts'
import type {IMachineListResp, ITopUpParams, IUserItemResp} from '@/types/userManager/user.ts'

export const userManager = {
  ban(userId: number) {
    return request.post<IUserItemResp>('admin/user/ban', {
      userId: userId,
    })
  },
  userList() {
    return request.post<IUserItemResp>('admin/user/list', null)
  },
  update(data?: IUpdateUserProfileParams) {
    return request.post<ICommonResponse<null>>('/admin/user/update', data)
  },
  TopUp(data?: ITopUpParams) {
    return request.post<ICommonResponse<null>>('admin/user/topUp', data)
  }
}

export const MachineManager = {
  list(params?: { deviceId: string; belong: string }) {
    return request.get<ICommonResponse<IMachineListResp>>('admin/machine/list', params)
  },
  ban(machineId: number) {
    return request.post<ICommonResponse<null>>('admin/machine/ban', {
      machineId: machineId,
    })
  },
  offline(machineId: number) {
    return request.post<ICommonResponse<null>>('admin/machine/offline', {
      machineId: machineId,
    })
  },
  readLog(machineId: number) {
    return request.post<ILogListResp>('admin/machine/readLog', {
      machineId: machineId,
    })
  }
}
