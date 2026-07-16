import { apiClient } from './apiClient'

export interface Budget {
  id: number
  name: string       
  limit_value: number
  current_value: number
}

export interface CreateBudgetPayload {
  category_id: number
  limit: number
}

export interface UpdateBudgetPayload {
  limit: number
}

export const getBudgets = async (date: string): Promise<{ result: Budget[] }> => {
  const response = await apiClient(`/budgets?date=${date}`, { method: 'GET' })
  if (!response.ok) {
    let errorMessage = 'Ошибка получения бюджетов'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const createBudget = async (payload: CreateBudgetPayload): Promise<{ message: string }> => {
  const response = await apiClient('/budgets', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    let errorMessage = 'Ошибка создания бюджета'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const deleteBudget = async (id: number): Promise<{ message: string }> => {
  const response = await apiClient(`/budgets/${id}`, { method: 'DELETE' })
  if (!response.ok) {
    let errorMessage = 'Ошибка удаления бюджета'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const updateBudget = async (id: number, payload: UpdateBudgetPayload): Promise<{ message: string }> => {
  const response = await apiClient(`/budgets/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    let errorMessage = 'Ошибка обновления бюджета'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}