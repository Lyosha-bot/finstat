import { useState, useEffect, useCallback } from 'react'
import {
  getBudgets,
  createBudget,
  deleteBudget,
  updateBudget,
  type Budget,
  type CreateBudgetPayload,
  type UpdateBudgetPayload,
} from '../api/budgets'

export function useBudgets(date: string) {
  const [budgets, setBudgets] = useState<Budget[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchBudgets = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getBudgets(date)
      const parsed = data.result.map(b => ({
        ...b,
        name: b.category_name,
        limit_value: typeof b.limit_value === 'string' ? parseFloat(b.limit_value) : b.limit_value,
        current_value: typeof b.current_value === 'string' ? parseFloat(b.current_value) : b.current_value,
      }))
      setBudgets(parsed)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }, [date])

  useEffect(() => {
    fetchBudgets()
  }, [fetchBudgets])

  const addBudget = useCallback(async (payload: CreateBudgetPayload) => {
    await createBudget(payload)
    await fetchBudgets()
  }, [fetchBudgets])

  const removeBudget = useCallback(async (id: number) => {
    await deleteBudget(id)
    await fetchBudgets()
  }, [fetchBudgets])

  const editBudget = useCallback(async (id: number, payload: UpdateBudgetPayload) => {
    await updateBudget(id, payload)
    await fetchBudgets()
  }, [fetchBudgets])

  return {
    budgets,
    loading,
    error,
    addBudget,
    removeBudget,
    editBudget,
    refetch: fetchBudgets,
  }
}