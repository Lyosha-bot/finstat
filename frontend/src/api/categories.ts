import { apiClient } from './apiClient'

export interface Category {
  id: number
  name: string
}

export interface CreateCategoryPayload {
  name: string
}

export interface UpdateCategoryPayload {
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

export const createCategory = async (payload: CreateCategoryPayload): Promise<{ id: number }> => {
  const response = await apiClient('/categories', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    let errorMessage = 'Ошибка создания категории'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const updateCategory = async (id: number, payload: UpdateCategoryPayload): Promise<void> => {
  const response = await apiClient(`/categories/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    let errorMessage = 'Ошибка обновления категории'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
}

export const deleteCategory = async (id: number): Promise<void> => {
  const response = await apiClient(`/categories/${id}`, { method: 'DELETE' })
  if (!response.ok) {
    let errorMessage = 'Ошибка удаления категории'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
}