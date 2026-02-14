import request from "@/utils/request.ts";
import type {
  IGetAppNamesResp,
  IGetVerInfoListParams,
  IGetVerInfoListResp,
  IVersionUpdateParams
} from '@/types/app/version.ts'


const AppVerService = {
  GetAppNames() {
    return request.get<IGetAppNamesResp>("/admin/app/version/appNameList", {})
  },
  GetVerInfoList(data?: IGetVerInfoListParams) {
    return request.post<IGetVerInfoListResp>("/admin/app/version/list", data)
  },
  Update(data? : IVersionUpdateParams){
    return request.post<IGetVerInfoListResp>("/admin/app/version/create", data)
  }
}

export default AppVerService
