import { memo, useState } from 'react'
import { SYSTEM_CATEGORIES } from '../constants'
import type { Category, CreateCategoryPayload, UpdateCategoryPayload } from '../api/categories'

interface CategoryManagerModalProps {
  isOpen: boolean
  onClose: () => void
  categories: Category[]
  onAddCategory: (payload: CreateCategoryPayload) => Promise<void>
  onEditCategory: (id: number, payload: UpdateCategoryPayload) => Promise<void>
  onDeleteCategory: (id: number) => Promise<void>
  onConfirm: (message: string, onConfirm?: () => void) => void
}

export const CategoryManagerModal = memo(({
  isOpen,
  onClose,
  categories = [],
  onAddCategory,
  onEditCategory,
  onDeleteCategory,
  onConfirm,
}: CategoryManagerModalProps) => {
  const [newCategoryName, setNewCategoryName] = useState('')
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editingName, setEditingName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const userCategories = categories.filter(c => !SYSTEM_CATEGORIES.includes(c.name))

  if (!isOpen) return null

  const handleAdd = async () => {
    if (newCategoryName.trim().length < 3) {
      setError('Название должно быть минимум 3 символа')
      return
    }
    setLoading(true)
    setError('')
    try {
      await onAddCategory({ name: newCategoryName.trim() })
      setNewCategoryName('')
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleSaveEdit = async () => {
    if (editingName.trim().length < 3) {
      setError('Название должно быть минимум 3 символа')
      return
    }
    setLoading(true)
    setError('')
    try {
      await onEditCategory(editingId!, { name: editingName.trim() })
      setEditingId(null)
      setEditingName('')
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = (id: number) => {
    onConfirm('Удалить категорию?', async () => {
      setLoading(true)
      try {
        await onDeleteCategory(id)
      } catch (err: any) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    })
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '500px' }}>
        <h2>Управление категориями</h2>
        <div className="form-group">
          <label>Новая категория</label>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <input
              type="text"
              value={newCategoryName}
              onChange={(e) => setNewCategoryName(e.target.value)}
              placeholder="Название (мин. 3 символа)"
              style={{ flex: 1 }}
            />
            <button className="btn btn-primary" onClick={handleAdd} disabled={loading}>
              Добавить
            </button>
          </div>
        </div>
        {error && <p className="error-text">{error}</p>}
        <div style={{ marginTop: '1.5rem', maxHeight: '300px', overflowY: 'auto' }}>
          <h3 style={{ fontSize: '0.9rem', color: 'var(--c-text-secondary)', marginBottom: '0.5rem' }}>Мои категории</h3>
          {userCategories.length === 0 ? (
            <p style={{ color: 'var(--c-text-muted)', fontSize: '0.9rem' }}>Нет пользовательских категорий</p>
          ) : (
            userCategories.map(cat => (
              <div
                key={cat.id}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '0.5rem 0',
                  borderBottom: '1px solid var(--c-border)',
                }}
              >
                {editingId === cat.id ? (
                  <input
                    type="text"
                    value={editingName}
                    onChange={(e) => setEditingName(e.target.value)}
                    autoFocus
                    style={{ flex: 1, marginRight: '0.5rem' }}
                  />
                ) : (
                  <span style={{ color: 'var(--c-text)' }}>{cat.name}</span>
                )}
                <div style={{ display: 'flex', gap: '0.3rem' }}>
                  {editingId === cat.id ? (
                    <>
                      <button className="btn btn-primary" onClick={handleSaveEdit} disabled={loading}>Сохранить</button>
                      <button className="btn btn-secondary" onClick={() => { setEditingId(null); setEditingName('') }}>Отмена</button>
                    </>
                  ) : (
                    <>
                      <button className="btn-edit" onClick={() => { setEditingId(cat.id); setEditingName(cat.name) }} disabled={loading}>✏️</button>
                      <button className="btn-delete" onClick={() => handleDelete(cat.id)} disabled={loading}>🗑️</button>
                    </>
                  )}
                </div>
              </div>
            ))
          )}
        </div>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={onClose}>Закрыть</button>
        </div>
      </div>
    </div>
  )
})