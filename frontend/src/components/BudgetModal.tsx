interface BudgetModalProps {
  isOpen: boolean
  onClose: () => void
  budgetForm: {
    category: string
    limit: string
    period: 'monthly' | 'weekly' | 'yearly'
    startDate: string
    endDate: string
  }
  setBudgetForm: (data: any) => void
  onSubmit: (e: React.FormEvent) => void
}

export const BudgetModal = ({ isOpen, onClose, budgetForm, setBudgetForm, onSubmit }: BudgetModalProps) => {
  if (!isOpen) return null

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <h2>Новый бюджет</h2>
        <form onSubmit={onSubmit}>
          <div className="form-group">
            <label>Категория</label>
            <input type="text" value={budgetForm.category} onChange={(e) => setBudgetForm({...budgetForm, category: e.target.value})} placeholder="Например: Продукты" required />
          </div>
          <div className="form-group">
            <label>Лимит (₽)</label>
            <input type="number" step="0.01" min="0.01" value={budgetForm.limit} onChange={(e) => setBudgetForm({...budgetForm, limit: e.target.value})} placeholder="15000" required />
          </div>
          <div className="form-group">
            <label>Период</label>
            <select value={budgetForm.period} onChange={(e) => {
              const period = e.target.value as 'monthly' | 'weekly' | 'yearly'
              const now = new Date()
              let start = now, end = now
              if (period === 'monthly') {
                start = new Date(now.getFullYear(), now.getMonth(), 1)
                end = new Date(now.getFullYear(), now.getMonth() + 1, 0)
              } else if (period === 'weekly') {
                const day = now.getDay() || 7
                start = new Date(now)
                start.setDate(now.getDate() - day + 1)
                end = new Date(start)
                end.setDate(start.getDate() + 6)
              } else if (period === 'yearly') {
                start = new Date(now.getFullYear(), 0, 1)
                end = new Date(now.getFullYear(), 11, 31)
              }
              setBudgetForm({ ...budgetForm, period, startDate: start.toISOString().split('T')[0], endDate: end.toISOString().split('T')[0] })
            }}>
              <option value="monthly">Ежемесячно</option>
              <option value="weekly">Еженедельно</option>
              <option value="yearly">Ежегодно</option>
            </select>
          </div>
          <div className="form-group">
            <label>Дата начала</label>
            <input type="date" value={budgetForm.startDate} onChange={(e) => setBudgetForm({...budgetForm, startDate: e.target.value})} required />
          </div>
          <div className="form-group">
            <label>Дата окончания</label>
            <input type="date" value={budgetForm.endDate} onChange={(e) => setBudgetForm({...budgetForm, endDate: e.target.value})} required />
          </div>
          <div className="modal-actions">
            <button type="submit" className="btn btn-primary">Создать</button>
            <button type="button" className="btn btn-secondary" onClick={onClose}>Отмена</button>
          </div>
        </form>
      </div>
    </div>
  )
}