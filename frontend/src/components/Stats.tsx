import { memo } from 'react'
import { formatMoney } from '../utils/format'

interface StatsProps {
  balance: number
  totalIncome: number
  totalExpense: number
}

export const Stats = memo(({ balance, totalIncome, totalExpense }: StatsProps) => {
  return (
    <section className="stats">
      <div className="stat-card">
        <div className="stat-label">Баланс</div>
        <div className="stat-value balance">{formatMoney(balance)}</div>
      </div>
      <div className="stat-card">
        <div className="stat-label">Доходы</div>
        <div className="stat-value income">{formatMoney(totalIncome)}</div>
      </div>
      <div className="stat-card">
        <div className="stat-label">Расходы</div>
        <div className="stat-value expense">{formatMoney(totalExpense)}</div>
      </div>
    </section>
  )
})