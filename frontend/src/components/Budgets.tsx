import type { Budget } from '../types'
import { formatMoney, formatDateRange } from '../utils/format'

interface BudgetsProps {
  budgets: Budget[]
  onDeleteBudget: (id: number) => void
  getCategorySpent: (category: string) => number
}

export const Budgets = ({ budgets, onDeleteBudget, getCategorySpent }: BudgetsProps) => {
  if (budgets.length === 0) return null

  return (
    <section className="budgets">
      <h2>Бюджеты</h2>
      <div className="budget-grid">
        {budgets.map(budget => {
          const spent = getCategorySpent(budget.category)
          const percent = Math.min((spent / budget.limit) * 100, 100)
          const color = percent > 90 ? '#ef4444' : percent > 70 ? '#f59e0b' : '#4ade80'
          return (
            <div key={budget.id} className="budget-card">
              <div className="budget-header">
                <span className="budget-category">{budget.category}</span>
                <button className="btn-delete-small" onClick={() => onDeleteBudget(budget.id)}>✕</button>
              </div>
              <div className="budget-dates">{formatDateRange(budget.startDate, budget.endDate)}</div>
              <div className="budget-amounts">
                <span>{formatMoney(spent)}</span>
                <span>/ {formatMoney(budget.limit)}</span>
              </div>
              <div className="budget-bar">
                <div className="budget-bar-fill" style={{ width: `${percent}%`, backgroundColor: color }} />
              </div>
              <div className="budget-percentage">{Math.round(percent)}%</div>
            </div>
          )
        })}
      </div>
    </section>
  )
}