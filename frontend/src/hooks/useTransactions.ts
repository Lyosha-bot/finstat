import { useLocalStorage } from './useLocalStorage'
import type { Transaction } from '../types'

export function useTransactions() {
  const [transactions, setTransactions] = useLocalStorage<Transaction[]>('transactions', [])

  const addTransaction = (transaction: Transaction) => {
    setTransactions([...transactions, transaction])
  }

  const deleteTransaction = (id: number) => {
    setTransactions(transactions.filter(t => t.id !== id))
  }

  const updateTransaction = (updated: Transaction) => {
    setTransactions(transactions.map(t => (t.id === updated.id ? updated : t)))
  }

  return { transactions, addTransaction, deleteTransaction, updateTransaction }
}