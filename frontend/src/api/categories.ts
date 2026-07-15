import { apiClient } from './apiClient'

export interface Category {
  id: number
  name: string
}

export const getCategories = async (): Promise<{ result: Category[] }> => {
  const response = await apiClient('/categories', { method: 'GET' })
  if (!response.ok) {
    let errorMessage = 'Ошибка получения категорий'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const getSystemCategories = async (): Promise<{ result: Category[] }> => {
  const response = await apiClient('/system-categories', { method: 'GET' })
  if (!response.ok) {
    let errorMessage = 'Ошибка получения системных категорий'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const getUserCategories = async (): Promise<{ result: Category[] }> => {
  const response = await apiClient('/user-categories', { method: 'GET' })
  if (!response.ok) {
    let errorMessage = 'Ошибка получения пользовательских категорий'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}