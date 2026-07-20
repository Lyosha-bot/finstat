import { useState, useEffect, useCallback } from 'react'
import {
  getCategories,
  createCategory,
  updateCategory,
  deleteCategory,
  type Category,
  type CreateCategoryPayload,
  type UpdateCategoryPayload,
} from '../api/categories'

export function useCategories() {
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchCategories = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getCategories()
      setCategories(data.result)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchCategories()
  }, [fetchCategories])

  const addCategory = useCallback(async (payload: CreateCategoryPayload) => {
    await createCategory(payload)
    await fetchCategories()
  }, [fetchCategories])

  const editCategory = useCallback(async (id: number, payload: UpdateCategoryPayload) => {
    await updateCategory(id, payload)
    await fetchCategories()
  }, [fetchCategories])

  const removeCategory = useCallback(async (id: number) => {
    await deleteCategory(id)
    await fetchCategories()
  }, [fetchCategories])

  return {
    categories,
    loading,
    error,
    addCategory,
    editCategory,
    removeCategory,
    refetch: fetchCategories,
  }
}