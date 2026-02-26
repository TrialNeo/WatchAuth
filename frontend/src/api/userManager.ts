import request from '@/utils/request.ts'
import type {IUpdateUserProfileParams, IUserListResponse} from '@/types/system/user.ts'
import type {ICommonResponse} from '@/types/common.ts'
import type {ITopUpParams} from '@/types/userManager/user.ts'

export const userManager = {
  ban(userId: number) {
    return request.post<IUserListResponse>('admin/user/ban', {
      userId: userId,
    })
  },
  userList() {
    return request.post<IUserListResponse>('admin/user/list', null)
  },
  update(data?: IUpdateUserProfileParams) {
    return request.post<ICommonResponse<null>>('/admin/user/update', data)
  },
  TopUp(data?: ITopUpParams) {
    return request.post<ICommonResponse<null>>('admin/user/topUp', data)
  }
}
