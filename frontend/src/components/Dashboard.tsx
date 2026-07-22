import React, { useState, useEffect, useMemo, useCallback } from 'react'
import type { Transaction } from '../types'
import { useTransactions } from '../hooks/useTransactions'
import { useCategories } from '../hooks/useCategories'
import { useBudgets } from '../hooks/useBudgets'
import { toLocalDateStr } from '../utils/format'
import {
  Header,
  Stats,
  Filters,
  Budgets,
  TransactionList,
  StatsDashboard,
  AddTransactionModal,
  BudgetModal,
  CategoryManagerModal,
} from './index'

interface DashboardProps {
  username: string
  onLogout: () => void
}

export const Dashboard = React.memo(({ username, onLogout }: DashboardProps) => {
  // ===== Категории =====
  const { categories, addCategory, editCategory, removeCategory } = useCategories()
  const [showCategoryManagerModal, setShowCategoryManagerModal] = useState(false)

  // ===== Бюджеты (получаем refetch) =====
  const currentDate = toLocalDateStr(new Date())
  const { budgets, loading: budgetsLoading, error: budgetsError, addBudget, removeBudget, editBudget, refetch: refetchBudgets } = useBudgets(currentDate)

  // ===== Модалка подтверждения =====
  const [confirmState, setConfirmState] = useState<{
    show: boolean
    message: string
    onConfirm?: () => void
    isError?: boolean
  }>({ show: false, message: '', isError: false })

  const openConfirm = useCallback((message: string, onConfirm?: () => void, isError = false) => {
    setConfirmState({ show: true, message, onConfirm, isError })
  }, [])

  const handleConfirm = useCallback(() => {
    if (confirmState.onConfirm) confirmState.onConfirm()
    setConfirmState({ show: false, message: '', isError: false })
  }, [confirmState])

  // ===== Модалка редактирования бюджета =====
  const [editBudgetState, setEditBudgetState] = useState<{ show: boolean; id: number | null; limit: string }>({
    show: false,
    id: null,
    limit: '',
  })

  const openEditBudget = useCallback((id: number, currentLimit: number) => {
    setEditBudgetState({ show: true, id, limit: String(currentLimit) })
  }, [])

  const handleEditBudgetSubmit = useCallback(async () => {
    const newLimit = parseFloat(editBudgetState.limit)
    if (isNaN(newLimit) || newLimit <= 0) {
      openConfirm('Введите корректное число', undefined, true)
      return
    }
    if (editBudgetState.id !== null) {
      try {
        await editBudget(editBudgetState.id, { limit: newLimit })
        setEditBudgetState({ show: false, id: null, limit: '' })
      } catch (err: any) {
        openConfirm(err.message || 'Ошибка обновления', undefined, true)
      }
    }
  }, [editBudgetState, editBudget, openConfirm])

  // ===== Фильтры =====
  const [filterType, setFilterType] = useState<'all' | 'income' | 'expense'>('all')
  const [filterCategory, setFilterCategory] = useState<number | 'all'>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [periodFilter, setPeriodFilter] = useState<'all' | 'today' | 'week' | 'month'>('all')
  const [statsDateFrom, setStatsDateFrom] = useState('')
  const [statsDateTo, setStatsDateTo] = useState('')

  const apiParams = useMemo(() => {
    const params: any = { limit: 1000, page: 1 }
    const now = new Date()
    if (periodFilter === 'today') {
      params.from = toLocalDateStr(new Date(now.getFullYear(), now.getMonth(), now.getDate()))
      params.to = toLocalDateStr(new Date(now.getFullYear(), now.getMonth(), now.getDate()))
    } else if (periodFilter === 'week') {
      const weekAgo = new Date(now)
      weekAgo.setDate(now.getDate() - 7)
      params.from = toLocalDateStr(weekAgo)
      params.to = toLocalDateStr(new Date(now.getFullYear(), now.getMonth(), now.getDate()))
    } else if (periodFilter === 'month') {
      const monthAgo = new Date(now)
      monthAgo.setMonth(now.getMonth() - 1)
      params.from = toLocalDateStr(monthAgo)
      params.to = toLocalDateStr(new Date(now.getFullYear(), now.getMonth(), now.getDate()))
    }
    if (filterType === 'income') params.type = 1
    else if (filterType === 'expense') params.type = -1
    if (filterCategory !== 'all') params.categories = [filterCategory]
    return params
  }, [periodFilter, filterType, filterCategory])

  // ===== Транзакции =====
  const {
    transactions,
    loading: transactionsLoading,
    error: transactionsError,
    addTransaction: addTransactionApi,
    updateTransaction: updateTransactionApi,
    deleteTransaction: deleteTransactionApi,
  } = useTransactions(apiParams)

  // ===== UI состояния =====
  const [activeTab, setActiveTab] = useState<'transactions' | 'stats'>('transactions')
  const [showAddModal, setShowAddModal] = useState(false)
  const [showBudgetModal, setShowBudgetModal] = useState(false)
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null)

  // ===== Форма транзакции =====
  const getDefaultCategoryId = useCallback(() => {
    if (categories.length > 0) return categories[0].id
    return 0
  }, [categories])

  const [formData, setFormData] = useState({
    description: '',
    amount: '',
    category: getDefaultCategoryId(),
    type: 'expense' as 'income' | 'expense',
    date: toLocalDateStr(new Date()),
  })

  useEffect(() => {
    if (categories.length > 0 && formData.category === 0) {
      setFormData(prev => ({ ...prev, category: categories[0].id }))
    }
  }, [categories, formData.category])

  // ===== Вычисления =====
  const filteredBySearch = useMemo(() => {
    if (!searchQuery.trim()) return transactions
    return transactions.filter(t => t.description.toLowerCase().includes(searchQuery.toLowerCase()))
  }, [transactions, searchQuery])

  const grouped = useMemo(() => {
    const groups: { [date: string]: Transaction[] } = {}
    filteredBySearch.forEach(t => {
      if (!groups[t.date]) groups[t.date] = []
      groups[t.date].push(t)
    })
    return Object.entries(groups)
      .sort(([a], [b]) => new Date(b).getTime() - new Date(a).getTime())
      .map(([date, items]) => ({
        date,
        total: items.reduce((sum, t) => sum + t.value, 0),
        transactions: items,
      }))
  }, [filteredBySearch])

  const statsData = useMemo(() => {
    let filtered = [...transactions]
    if (statsDateFrom) filtered = filtered.filter(t => t.date >= statsDateFrom)
    if (statsDateTo) filtered = filtered.filter(t => t.date <= statsDateTo)
    return filtered
  }, [transactions, statsDateFrom, statsDateTo])

  const totalIncome = statsData.filter(t => t.value > 0).reduce((s, t) => s + t.value, 0)
  const totalExpense = statsData.filter(t => t.value < 0).reduce((s, t) => s + t.value, 0)
  const balance = totalIncome + totalExpense

  const incomeStats = useMemo(() => {
    const map: Record<number, number> = {}
    statsData.filter(t => t.value > 0).forEach(t => {
      map[t.category_id] = (map[t.category_id] || 0) + t.value
    })
    return Object.entries(map)
      .map(([id, amount]) => {
        const cat = categories.find(c => c.id === Number(id))
        return [cat ? cat.name : 'Без категории', amount] as [string, number]
      })
      .sort((a, b) => b[1] - a[1])
      .slice(0, 10)
  }, [statsData, categories])

  const expenseStats = useMemo(() => {
    const map: Record<number, number> = {}
    statsData.filter(t => t.value < 0).forEach(t => {
      map[t.category_id] = (map[t.category_id] || 0) + t.value
    })
    return Object.entries(map)
      .map(([id, amount]) => {
        const cat = categories.find(c => c.id === Number(id))
        return [cat ? cat.name : 'Без категории', Math.abs(amount)] as [string, number]
      })
      .sort((a, b) => b[1] - a[1])
      .slice(0, 10)
  }, [statsData, categories])

  const monthlyStats = useMemo(() => {
    const months: Record<string, { income: number; expense: number }> = {}
    statsData.forEach(t => {
      const m = new Date(t.date).toLocaleString('ru-RU', { month: 'short', year: 'numeric' })
      if (!months[m]) months[m] = { income: 0, expense: 0 }
      if (t.value > 0) months[m].income += t.value
      else months[m].expense += Math.abs(t.value)
    })
    return Object.entries(months).sort((a, b) => new Date(a[0]).getTime() - new Date(b[0]).getTime()).slice(-10)
  }, [statsData])

  const cumulative = useMemo(() => {
    const dailyMap: Record<string, number> = {}
    statsData.forEach(t => {
      dailyMap[t.date] = (dailyMap[t.date] || 0) + t.value
    })
    const sortedDates = Object.keys(dailyMap).sort((a, b) => new Date(a).getTime() - new Date(b).getTime())
    let cum = 0
    const result: { label: string; value: number }[] = []
    sortedDates.forEach(date => {
      cum += dailyMap[date]
      result.push({ label: date, value: cum })
    })
    return result.slice(-10)
  }, [statsData])

  const avgDaily = useMemo(() => {
    const days = new Set(statsData.map(t => t.date)).size || 1
    const totalInc = statsData.filter(t => t.value > 0).reduce((s, t) => s + t.value, 0)
    const totalExp = statsData.filter(t => t.value < 0).reduce((s, t) => s + t.value, 0)
    return { avgIncome: totalInc / days, avgExpense: Math.abs(totalExp) / days }
  }, [statsData])

  const weekdayStats = useMemo(() => {
    const weekdays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
    const map: Record<string, number> = { Пн: 0, Вт: 0, Ср: 0, Чт: 0, Пт: 0, Сб: 0, Вс: 0 }
    statsData.filter(t => t.value < 0).forEach(t => {
      const d = new Date(t.date)
      let wd = d.getDay()
      if (wd === 0) wd = 7
      const name = weekdays[wd - 1]
      map[name] = (map[name] || 0) + Math.abs(t.value)
    })
    return Object.entries(map).map(([name, value]) => ({ name, value }))
  }, [statsData])

  const maxIncome = incomeStats.length ? incomeStats[0][1] : 1
  const maxExpense = expenseStats.length ? expenseStats[0][1] : 1
  const maxMonthly = monthlyStats.reduce((max, [_, { income, expense }]) => Math.max(max, income, expense), 1)
  const maxWeekday = weekdayStats.reduce((max, { value }) => Math.max(max, value), 1)

  // ===== Обработчики транзакций =====
  const handleAddTransaction = useCallback(async (e: React.FormEvent) => {
    e.preventDefault()
    const amount = parseFloat(formData.amount)
    if (!formData.description || isNaN(amount) || amount <= 0) {
      openConfirm('Заполните поля корректно', undefined, true)
      return
    }
    if (!formData.category) {
      openConfirm('Выберите категорию', undefined, true)
      return
    }
    const payload = {
      description: formData.description,
      amount: formData.type === 'income' ? amount : -amount,
      category_id: formData.category,
      date: formData.date,
    }
    try {
      await addTransactionApi(payload)
      await refetchBudgets()
      setFormData({
        description: '',
        amount: '',
        category: categories.length > 0 ? categories[0].id : 0,
        type: 'expense',
        date: toLocalDateStr(new Date()),
      })
      setShowAddModal(false)
    } catch (err: any) {
      openConfirm(err.message || 'Ошибка создания', undefined, true)
    }
  }, [formData, addTransactionApi, categories, openConfirm, refetchBudgets])

  const handleEditTransaction = useCallback((transaction: Transaction) => {
    setEditingTransaction(transaction)
    setFormData({
      description: transaction.description,
      amount: Math.abs(transaction.value).toString(),
      category: transaction.category_id,
      type: transaction.value >= 0 ? 'income' : 'expense',
      date: transaction.date,
    })
    setShowAddModal(true)
  }, [])

  const handleUpdateTransaction = useCallback(async (e: React.FormEvent) => {
    e.preventDefault()
    const amount = parseFloat(formData.amount)
    if (!formData.description || isNaN(amount) || amount <= 0) {
      openConfirm('Заполните поля корректно', undefined, true)
      return
    }
    if (!formData.category) {
      openConfirm('Выберите категорию', undefined, true)
      return
    }
    if (!editingTransaction) return
    const payload = {
      description: formData.description,
      amount: formData.type === 'income' ? amount : -amount,
      category_id: formData.category,
      date: formData.date,
    }
    try {
      await updateTransactionApi(editingTransaction.id, payload)
      await refetchBudgets()
      setEditingTransaction(null)
      setFormData({
        description: '',
        amount: '',
        category: categories.length > 0 ? categories[0].id : 0,
        type: 'expense',
        date: toLocalDateStr(new Date()),
      })
      setShowAddModal(false)
    } catch (err: any) {
      openConfirm(err.message || 'Ошибка обновления', undefined, true)
    }
  }, [formData, editingTransaction, updateTransactionApi, categories, openConfirm, refetchBudgets])

  const handleDeleteTransaction = useCallback((id: number) => {
    openConfirm('Удалить транзакцию?', async () => {
      try {
        await deleteTransactionApi(id)
        await refetchBudgets()
      } catch (err: any) {
        openConfirm(err.message || 'Ошибка удаления', undefined, true)
      }
    })
  }, [deleteTransactionApi, openConfirm, refetchBudgets])

  // ===== Обработчики бюджетов =====
  const handleCreateBudget = useCallback(async (payload: { category_id: number; limit: number }) => {
    try {
      await addBudget(payload)
    } catch (err: any) {
      openConfirm(err.message || 'Ошибка создания бюджета', undefined, true)
    }
  }, [addBudget, openConfirm])

  const handleDeleteBudget = useCallback((id: number) => {
    openConfirm('Удалить бюджет?', async () => {
      try {
        await removeBudget(id)
      } catch (err: any) {
        openConfirm(err.message || 'Ошибка удаления бюджета', undefined, true)
      }
    })
  }, [removeBudget, openConfirm])

  const handleEditBudget = useCallback((id: number, currentLimit: number) => {
    openEditBudget(id, currentLimit)
  }, [openEditBudget])

  // ===== Обработчики категорий =====
  const handleAddCategory = useCallback(async (payload: { name: string }) => {
    await addCategory(payload)
  }, [addCategory])

  const handleEditCategory = useCallback(async (id: number, payload: { name: string }) => {
    await editCategory(id, payload)
  }, [editCategory])

  const handleDeleteCategory = useCallback(async (id: number) => {
    await removeCategory(id)
    await refetchBudgets() // обновляем бюджеты после удаления категории
  }, [removeCategory, refetchBudgets])

  // ===== Рендер =====
  return (
    <div className="app">
      <Header
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        onAddClick={() => {
          setEditingTransaction(null)
          setFormData({
            description: '',
            amount: '',
            category: categories.length > 0 ? categories[0].id : 0,
            type: 'expense',
            date: toLocalDateStr(new Date()),
          })
          setShowAddModal(true)
        }}
        onBudgetClick={() => setShowBudgetModal(true)}
        onCategoryManagerClick={() => setShowCategoryManagerModal(true)}
        onLogout={onLogout}
        username={username}
      />

      <Stats balance={balance} totalIncome={totalIncome} totalExpense={Math.abs(totalExpense)} />

      {activeTab === 'transactions' && (
        <Filters
          periodFilter={periodFilter}
          setPeriodFilter={setPeriodFilter}
          filterType={filterType}
          setFilterType={setFilterType}
          filterCategory={filterCategory}
          setFilterCategory={setFilterCategory}
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          categories={categories}
        />
      )}

      {activeTab === 'transactions' ? (
        <>
          <Budgets
            budgets={budgets}
            loading={budgetsLoading}
            error={budgetsError}
            onDeleteBudget={handleDeleteBudget}
            onEditBudget={handleEditBudget}
          />
          <section className="transactions">
            <h2>Транзакции</h2>
            {transactionsLoading && <p>Загрузка...</p>}
            {transactionsError && <p className="error-text">{transactionsError}</p>}
            {!transactionsLoading && !transactionsError && (
              <TransactionList
                grouped={grouped}
                onEdit={handleEditTransaction}
                onDelete={handleDeleteTransaction}
                categories={categories}
              />
            )}
          </section>
        </>
      ) : (
        <StatsDashboard
          incomeStats={incomeStats}
          expenseStats={expenseStats}
          budgetStats={budgets.map(b => ({
            ...b,
            category: b.name,
            spent: Math.abs(b.current_value),
            percent: Math.min((Math.abs(b.current_value) / b.limit_value) * 100, 100),
          }))}
          cumulative={cumulative}
          avgDaily={avgDaily}
          weekdayStats={weekdayStats}
          monthlyStats={monthlyStats}
          maxIncome={maxIncome}
          maxExpense={maxExpense}
          maxMonthly={maxMonthly}
          maxWeekday={maxWeekday}
          statsDateFrom={statsDateFrom}
          statsDateTo={statsDateTo}
          setStatsDateFrom={setStatsDateFrom}
          setStatsDateTo={setStatsDateTo}
        />
      )}

      <AddTransactionModal
        isOpen={showAddModal}
        onClose={() => {
          setShowAddModal(false)
          setEditingTransaction(null)
        }}
        formData={formData}
        setFormData={setFormData}
        onSubmit={editingTransaction ? handleUpdateTransaction : handleAddTransaction}
        isEditing={!!editingTransaction}
        categories={categories}
      />

      <BudgetModal
        isOpen={showBudgetModal}
        onClose={() => setShowBudgetModal(false)}
        categories={categories}
        onCreateBudget={handleCreateBudget}
      />

      <CategoryManagerModal
        isOpen={showCategoryManagerModal}
        onClose={() => setShowCategoryManagerModal(false)}
        categories={categories}
        onAddCategory={handleAddCategory}
        onEditCategory={handleEditCategory}
        onDeleteCategory={handleDeleteCategory}
        onConfirm={openConfirm}
      />

      {confirmState.show && (
        <div className="modal-overlay" onClick={() => setConfirmState({ show: false, message: '', isError: false })}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2>{confirmState.isError ? 'Ошибка' : 'Подтверждение'}</h2>
            <p style={{ marginBottom: '1.5rem' }}>{confirmState.message}</p>
            <div className="modal-actions">
              <button className="btn btn-primary" onClick={handleConfirm}>
                {confirmState.isError ? 'ОК' : 'Да'}
              </button>
              {!confirmState.isError && (
                <button className="btn btn-secondary" onClick={() => setConfirmState({ show: false, message: '', isError: false })}>
                  Отмена
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {editBudgetState.show && (
        <div className="modal-overlay" onClick={() => setEditBudgetState({ show: false, id: null, limit: '' })}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2>Редактировать лимит</h2>
            <div className="form-group">
              <label>Новый лимит (₽)</label>
              <input
                type="number"
                step="0.01"
                min="0.01"
                value={editBudgetState.limit}
                onChange={(e) => setEditBudgetState(prev => ({ ...prev, limit: e.target.value }))}
                placeholder="15000"
                autoFocus
              />
            </div>
            <div className="modal-actions">
              <button className="btn btn-primary" onClick={handleEditBudgetSubmit}>Сохранить</button>
              <button className="btn btn-secondary" onClick={() => setEditBudgetState({ show: false, id: null, limit: '' })}>Отмена</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
})