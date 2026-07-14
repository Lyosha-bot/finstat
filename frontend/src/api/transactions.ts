import { apiClient } from './apiClient'
import type { Transaction } from '../types'

export interface TransactionResponse {
  result: Transaction[]
}

export interface CreateTransactionPayload {
  amount: number
  category: number
  date: string
  description: string
}

export interface UpdateTransactionPayload {
  amount: number
  category: number
  date: string
  description: string
}

export const getTransactions = async (params: {
  from?: string
  to?: string
  categories?: number[]
  type?: number
  limit?: number
  page?: number
}): Promise<TransactionResponse> => {
  const searchParams = new URLSearchParams()
  if (params.from) searchParams.append('from', params.from)
  if (params.to) searchParams.append('to', params.to)
  if (params.categories && params.categories.length) {
    params.categories.forEach(id => searchParams.append('categories', id.toString()))
  }
  if (params.type !== undefined) searchParams.append('type', params.type.toString())
  if (params.limit) searchParams.append('limit', params.limit.toString())
  if (params.page) searchParams.append('page', params.page.toString())
  const url = `/transactions?${searchParams.toString()}`
  const response = await apiClient(url, { method: 'GET' })
  if (!response.ok) {
    let errorMessage = 'Ошибка получения транзакций'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const createTransaction = async (payload: CreateTransactionPayload): Promise<{ message: string }> => {
  const response = await apiClient('/transactions', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    let errorMessage = 'Ошибка создания транзакции'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const updateTransaction = async (id: number, payload: UpdateTransactionPayload): Promise<{ message: string }> => {
  const response = await apiClient(`/transactions/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    let errorMessage = 'Ошибка обновления транзакции'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}

export const deleteTransaction = async (id: number): Promise<{ message: string }> => {
  const response = await apiClient(`/transactions/${id}`, { method: 'DELETE' })
  if (!response.ok) {
    let errorMessage = 'Ошибка удаления транзакции'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}