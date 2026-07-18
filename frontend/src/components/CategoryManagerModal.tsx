import { useState } from 'react'
import { useCategories } from '../hooks/useCategories'

interface CategoryManagerModalProps {
  isOpen: boolean
  onClose: () => void
}

export const CategoryManagerModal = ({ isOpen, onClose }: CategoryManagerModalProps) => {
  const { categories, addCategory, editCategory, removeCategory } = useCategories()
  const [newCategoryName, setNewCategoryName] = useState('')
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editingName, setEditingName] = useState('')

  const handleAdd = async () => {
    if (newCategoryName.trim().length < 3) {
      alert('Название категории должно содержать минимум 3 символа')
      return
    }
    try {
      await addCategory({ name: newCategoryName.trim() })
      setNewCategoryName('')
    } catch (err: any) {
      alert(err.message)
    }
  }

  const handleEdit = (id: number, name: string) => {
    setEditingId(id)
    setEditingName(name)
  }

  const handleSaveEdit = async () => {
    if (editingName.trim().length < 3) {
      alert('Название категории должно содержать минимум 3 символа')
      return
    }
    try {
      await editCategory(editingId!, { name: editingName.trim() })
      setEditingId(null)
      setEditingName('')
    } catch (err: any) {
      alert(err.message)
    }
  }

  const handleDelete = async (id: number) => {
    if (window.confirm('Удалить категорию?')) {
      try {
        await removeCategory(id)
      } catch (err: any) {
        alert(err.message)
      }
    }
  }

  if (!isOpen) return null

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '500px' }}>
        <h2>Управление категориями</h2>
        <div className="form-group" style={{ display: 'flex', gap: '0.5rem' }}>
          <input
            type="text"
            value={newCategoryName}
            onChange={(e) => setNewCategoryName(e.target.value)}
            placeholder="Название новой категории (мин. 3 символа)"
            style={{ flex: 1 }}
          />
          <button className="btn btn-primary" onClick={handleAdd}>Добавить</button>
        </div>
        <div style={{ marginTop: '1rem', maxHeight: '300px', overflowY: 'auto' }}>
          {categories.map(cat => (
            <div key={cat.id} className="category-item" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0.4rem 0', borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
              {editingId === cat.id ? (
                <input
                  type="text"
                  value={editingName}
                  onChange={(e) => setEditingName(e.target.value)}
                  autoFocus
                  style={{ flex: 1, marginRight: '0.5rem' }}
                />
              ) : (
                <span>{cat.name}</span>
              )}
              <div style={{ display: 'flex', gap: '0.3rem' }}>
                {editingId === cat.id ? (
                  <>
                    <button className="btn btn-primary" onClick={handleSaveEdit}>Сохранить</button>
                    <button className="btn btn-secondary" onClick={() => { setEditingId(null); setEditingName('') }}>Отмена</button>
                  </>
                ) : (
                  <>
                    <button className="btn-edit" onClick={() => handleEdit(cat.id, cat.name)}>✏️</button>
                    <button className="btn-delete" onClick={() => handleDelete(cat.id)}>🗑️</button>
                  </>
                )}
              </div>
            </div>
          ))}
        </div>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={onClose}>Закрыть</button>
        </div>
      </div>
    </div>
  )
}