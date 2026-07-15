import { useState, useEffect, useCallback, useMemo } from 'react'
import {
  getTransactions,
  createTransaction,
  updateTransaction,
  deleteTransaction,
  type CreateTransactionPayload,
  type UpdateTransactionPayload,
} from '../api/transactions'
import type { Transaction } from '../types'

interface UseTransactionsParams {
  from?: string
  to?: string
  categories?: number[]
  type?: number
  limit?: number
  page?: number
}

export function useTransactions(params: UseTransactionsParams = {}) {
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const paramsKey = useMemo(() => JSON.stringify(params), [params])

  const fetchTransactions = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getTransactions(params)
      // Преобразуем value в число (бэкенд может возвращать строку из-за decimal.Decimal)
      const parsed = data.result.map(t => ({
        ...t,
        value: typeof t.value === 'string' ? parseFloat(t.value) : t.value,
      }))
      setTransactions(parsed)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }, [paramsKey])

  useEffect(() => {
    fetchTransactions()
  }, [fetchTransactions])

  const addTransaction = async (payload: CreateTransactionPayload) => {
    await createTransaction(payload)
    await fetchTransactions()
  }

  const updateTransactionItem = async (id: number, payload: UpdateTransactionPayload) => {
    await updateTransaction(id, payload)
    await fetchTransactions()
  }

  const deleteTransactionItem = async (id: number) => {
    await deleteTransaction(id)
    await fetchTransactions()
  }

  return {
    transactions,
    loading,
    error,
    refetch: fetchTransactions,
    addTransaction,
    updateTransaction: updateTransactionItem,
    deleteTransaction: deleteTransactionItem,
  }
}