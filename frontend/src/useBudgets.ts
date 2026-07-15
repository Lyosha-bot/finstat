import { useLocalStorage } from './useLocalStorage'
import { Budget } from '../types'

export function useBudgets(initialBudgets: Budget[]) {
  return useLocalStorage<Budget[]>('budgets', initialBudgets)
}