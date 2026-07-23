import { memo } from 'react'
import type { Category } from '../api/categories'

interface AddTransactionModalProps {
  isOpen: boolean
  onClose: () => void
  formData: {
    description: string
    amount: string
    category: number
    type: 'income' | 'expense'
    date: string
  }
  setFormData: (data: any) => void
  onSubmit: (e: React.FormEvent) => void
  isEditing: boolean
  categories: Category[]
}

export const AddTransactionModal = memo(({
  isOpen,
  onClose,
  formData,
  setFormData,
  onSubmit,
  isEditing,
  categories,
}: AddTransactionModalProps) => {
  if (!isOpen) return null

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <h2>{isEditing ? 'Редактировать' : 'Новая транзакция'}</h2>
        <form onSubmit={onSubmit}>
          <div className="form-group">
            <label>Тип</label>
            <select
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value as 'income' | 'expense' })}
            >
              <option value="expense">Расход</option>
              <option value="income">Доход</option>
            </select>
          </div>
          <div className="form-group">
            <label>Описание</label>
            <input
              type="text"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Например: Продукты"
              required
            />
          </div>
          <div className="form-group">
            <label>Сумма (₽)</label>
            <input
              type="number"
              step="0.01"
              min="0.01"
              value={formData.amount}
              onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
              placeholder="1000"
              required
            />
          </div>
          <div className="form-group">
            <label>Категория</label>
            <select
              value={formData.category}
              onChange={(e) => setFormData({ ...formData, category: Number(e.target.value) })}
            >
              {categories.map(cat => (
                <option key={cat.id} value={cat.id}>{cat.name}</option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label>Дата</label>
            <input
              type="date"
              value={formData.date}
              onChange={(e) => setFormData({ ...formData, date: e.target.value })}
              required
            />
          </div>
          <div className="modal-actions">
            <button type="submit" className="btn btn-primary">
              {isEditing ? 'Сохранить' : 'Добавить'}
            </button>
            <button type="button" className="btn btn-secondary" onClick={onClose}>
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  )
})