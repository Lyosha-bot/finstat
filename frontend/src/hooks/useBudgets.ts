import { useState, useEffect } from 'react'
import { getBudgets, createBudget, type Budget, type CreateBudgetPayload } from '../api/budgets'

export function useBudgets(date: string) {
  const [budgets, setBudgets] = useState<Budget[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchBudgets = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getBudgets(date)
      const parsed = data.result.map(b => ({
        ...b,
        limit_value: typeof b.limit_value === 'string' ? parseFloat(b.limit_value) : b.limit_value,
        current_value: typeof b.current_value === 'string' ? parseFloat(b.current_value) : b.current_value,
      }))
      setBudgets(parsed)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchBudgets()
  }, [date])

  const addBudget = async (payload: CreateBudgetPayload) => {
    await createBudget(payload)
    await fetchBudgets()
  }

  return { budgets, loading, error, addBudget, refetch: fetchBudgets }
}