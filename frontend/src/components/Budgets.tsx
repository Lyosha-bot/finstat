import type { Budget } from '../api/budgets'
import { formatMoney } from '../utils/format'

interface BudgetsProps {
  budgets: Budget[]
  loading: boolean
  error: string | null
  onDeleteBudget: (id: number) => void
  onEditBudget: (id: number, currentLimit: number) => void
}

export const Budgets = ({ budgets, loading, error, onDeleteBudget, onEditBudget }: BudgetsProps) => {
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
          const spent = Math.abs(budget.current_value)
          const percent = Math.min((spent / budget.limit_value) * 100, 100)
          const color = percent > 90 ? '#ef4444' : percent > 70 ? '#f59e0b' : '#4ade80'
          return (
            <div key={budget.id} className="budget-card">
              <div className="budget-header">
                <span className="budget-category">{budget.name}</span>
                <div className="budget-actions">
                  <button 
                    className="btn-edit" 
                    onClick={() => onEditBudget(budget.id, budget.limit_value)}
                  >
                    Редактировать
                  </button>
                  <button 
                    className="btn-delete" 
                    onClick={() => onDeleteBudget(budget.id)}
                  >
                    Удалить
                  </button>
                </div>
              </div>
              <div className="budget-amounts">
                <span>{formatMoney(spent)}</span>
                <span>/ {formatMoney(budget.limit_value)}</span>
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