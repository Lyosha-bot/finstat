import { type Category } from '../api/categories'
import { type CreateBudgetPayload } from '../api/budgets'
import { useState } from 'react'

interface BudgetModalProps {
  isOpen: boolean
  onClose: () => void
  categories: Category[]
  onCreateBudget: (payload: CreateBudgetPayload) => Promise<void>
}

export const BudgetModal = ({ isOpen, onClose, categories, onCreateBudget }: BudgetModalProps) => {
  const [selectedCategoryId, setSelectedCategoryId] = useState<number | ''>('')
  const [limit, setLimit] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  if (!isOpen) return null

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    if (selectedCategoryId === '' || !limit || parseFloat(limit) <= 0) {
      setError('Выберите категорию и укажите корректный лимит')
      setLoading(false)
      return
    }

    try {
      await onCreateBudget({
        category_id: selectedCategoryId,
        limit: parseFloat(limit),
      })
      onClose()
      setSelectedCategoryId('')
      setLimit('')
    } catch (err: any) {
      setError(err.message || 'Ошибка создания бюджета')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <h2>Новый бюджет</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Категория</label>
            <select
              value={selectedCategoryId}
              onChange={(e) => setSelectedCategoryId(Number(e.target.value))}
              required
            >
              <option value="">Выберите категорию</option>
              {categories.map(cat => (
                <option key={cat.id} value={cat.id}>{cat.name}</option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label>Лимит (₽)</label>
            <input
              type="number"
              step="0.01"
              min="0.01"
              value={limit}
              onChange={(e) => setLimit(e.target.value)}
              placeholder="15000"
              required
            />
          </div>
          {error && <div className="error-text">{error}</div>}
          <div className="modal-actions">
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Создание...' : 'Создать'}
            </button>
            <button type="button" className="btn btn-secondary" onClick={onClose}>
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}