import request from '@/utils/request.ts'
import type { ICommonResponse } from '@/types/common.ts'

export interface AnnouncementItem {
  id: number
  title: string
  content: string
  status: number
  createdAt: string
  updatedAt: string
}

export const getAnnouncementList = () => {
  return request.get<ICommonResponse<{ list: AnnouncementItem[] }>>('/admin/announcement/list')
}

export const getActiveAnnouncements = () => {
  return request.get<ICommonResponse<{ list: AnnouncementItem[] }>>('/admin/announcement/active')
}

export const createAnnouncement = (data: { title: string; content: string; status: number }) => {
  return request.post<ICommonResponse<{ id: number }>>('/admin/announcement/create', data)
}

export const updateAnnouncement = (data: { id: number; title: string; content: string; status: number }) => {
  return request.post<ICommonResponse<null>>('/admin/announcement/update', data)
}

export const deleteAnnouncement = (id: number) => {
  return request.post<ICommonResponse<null>>('/admin/announcement/delete', { id })
}
