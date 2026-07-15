import type { Transaction } from '../types'
import type { Category } from '../api/categories'
import { formatMoney, formatDate } from '../utils/format'

interface TransactionListProps {
  grouped: { date: string; total: number; transactions: Transaction[] }[]
  onEdit: (t: Transaction) => void
  onDelete: (id: number) => void
  categories: Category[]
}

export const TransactionList = ({ grouped, onEdit, onDelete, categories }: TransactionListProps) => {
  const getCategoryName = (id: number) => {
    const cat = categories.find(c => c.id === id)
    return cat ? cat.name : 'Без категории'
  }

  if (grouped.length === 0) {
    return (
      <div className="empty-state">
        <p>Нет транзакций</p>
        <p className="empty-sub">Добавьте первую транзакцию</p>
      </div>
    )
  }

  return (
    <div className="transaction-groups">
      {grouped.map(group => (
        <div key={group.date} className="transaction-group">
          <div className="group-header">
            <span className="group-date">{formatDate(group.date)}</span>
            <span className="group-total">
              <span className={group.total >= 0 ? 'income' : 'expense'}>
                {group.total >= 0 ? '+' : ''}{formatMoney(group.total)}
              </span>
            </span>
          </div>
          <div className="group-transactions">
            {group.transactions.map(t => (
              <div key={t.id} className="transaction-item">
                <div className="transaction-info">
                  <div className="transaction-description">{t.description}</div>
                  <div className="transaction-category">
                    <span className="category-badge">{getCategoryName(t.category_id)}</span>
                  </div>
                </div>
                <div className="transaction-amounts">
                  <span className={t.value >= 0 ? 'income' : 'expense'}>
                    {t.value >= 0 ? '+' : '-'}{formatMoney(Math.abs(t.value))}
                  </span>
                  <div className="transaction-actions">
                    <button className="btn-edit" onClick={() => onEdit(t)}>✏️</button>
                    <button className="btn-delete" onClick={() => onDelete(t.id)}>🗑️</button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}