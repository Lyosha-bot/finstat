import type { Budget } from '../api/budgets'
import { formatMoney } from '../utils/format'

interface BudgetsProps {
  budgets: Budget[]
  loading: boolean
  error: string | null
}

export const Budgets = ({ budgets, loading, error }: BudgetsProps) => {
  if (loading) {
    return (
      <section className="budgets">
        <h2>Бюджеты</h2>
        <p>Загрузка...</p>
      </section>
    )
  }

  if (error) {
    return (
      <section className="budgets">
        <h2>Бюджеты</h2>
        <p className="error-text">{error}</p>
      </section>
    )
  }

  if (budgets.length === 0) {
    return (
      <section className="budgets">
        <h2>Бюджеты</h2>
        <p>Нет бюджетов</p>
      </section>
    )
  }

  return (
    <section className="budgets">
      <h2>Бюджеты</h2>
      <div className="budget-grid">
        {budgets.map(budget => {
          const current = Number(budget.current_value) || 0
          const limit = Number(budget.limit_value) || 0
          const spent = Math.abs(current)
          const percent = limit > 0 ? Math.min((spent / limit) * 100, 100) : 0
          const color = percent > 90 ? '#ef4444' : percent > 70 ? '#f59e0b' : '#4ade80'
          return (
            <div key={budget.id} className="budget-card">
              <div className="budget-header">
                <span className="budget-category">{budget.name}</span>
              </div>
              <div className="budget-amounts">
                <span>{formatMoney(spent)}</span>
                <span>/ {formatMoney(limit)}</span>
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