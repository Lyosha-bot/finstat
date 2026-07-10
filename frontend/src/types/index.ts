export interface Transaction {
  id: number
  date: string
  description: string
  amount: number
  category: string
  type: 'income' | 'expense'
}

export interface Budget {
  id: number
  category: string
  limit: number
  spent: number
  period: 'monthly' | 'weekly' | 'yearly'
  startDate: string
  endDate: string
}