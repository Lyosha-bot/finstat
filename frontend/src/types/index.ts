export interface Transaction {
  id: number
  userID: number
  value: number          
  category_id: number
  description: string
  date: string
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