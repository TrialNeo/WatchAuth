// 通用类型文件

// 通用响应
export interface ICommonResponse<T> {
  code: number
  msg: string
  data: T
}
