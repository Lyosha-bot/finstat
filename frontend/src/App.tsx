/// <reference types="vite/client" />
import { useState, useEffect, useMemo } from 'react'
import './App.css'

import type { Transaction } from './types'
import { useTransactions } from './hooks/useTransactions'
import { useCategories } from './hooks/useCategories'
import { useBudgets } from './hooks/useBudgets'
import {
  Login,
  Header,
  Stats,
  Filters,
  Budgets,
  TransactionList,
  StatsDashboard,
  AddTransactionModal,
  BudgetModal,
} from './components'

function App() {
  // ===== Авторизация =====
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    return localStorage.getItem('auth') === 'true'
  })
  const [username, setUsername] = useState(() => {
    return localStorage.getItem('username') || 'Пользователь'
  })

  const handleLogin = (name: string) => {
    setIsAuthenticated(true)
    setUsername(name)
  }

  const handleLogout = () => {
    localStorage.removeItem('auth')
    localStorage.removeItem('username')
    setIsAuthenticated(false)
  }

  // ===== Категории =====
  const { categories, loading: categoriesLoading } = useCategories()

  // ===== Бюджеты =====
  const currentDate = new Date().toISOString().split('T')[0]
  const { budgets, loading: budgetsLoading, error: budgetsError, addBudget } = useBudgets(currentDate)

  // ===== Фильтры =====
  const [filterType, setFilterType] = useState<'all' | 'income' | 'expense'>('all')
  const [filterCategory, setFilterCategory] = useState<number | 'all'>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [periodFilter, setPeriodFilter] = useState<'all' | 'today' | 'week' | 'month'>('all')
  const [statsDateFrom, setStatsDateFrom] = useState('')
  const [statsDateTo, setStatsDateTo] = useState('')

  // ===== Параметры для API =====
  const apiParams = useMemo(() => {
    const params: any = {
      limit: 1000,
      page: 1,
    }
    const now = new Date()
    if (periodFilter === 'today') {
      params.from = new Date(now.getFullYear(), now.getMonth(), now.getDate()).toISOString().split('T')[0]
      params.to = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1).toISOString().split('T')[0]
    } else if (periodFilter === 'week') {
      const weekAgo = new Date(now)
      weekAgo.setDate(now.getDate() - 7)
      params.from = weekAgo.toISOString().split('T')[0]
      params.to = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1).toISOString().split('T')[0]
    } else if (periodFilter === 'month') {
      const monthAgo = new Date(now)
      monthAgo.setMonth(now.getMonth() - 1)
      params.from = monthAgo.toISOString().split('T')[0]
      params.to = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1).toISOString().split('T')[0]
    }
    if (filterType === 'income') params.type = 1
    else if (filterType === 'expense') params.type = -1
    if (filterCategory !== 'all') {
      params.categories = [filterCategory]
    }
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
  const getDefaultCategoryId = () => {
    if (categories.length > 0) return categories[0].id
    return 0
  }

  const [formData, setFormData] = useState({
    description: '',
    amount: '',
    category: getDefaultCategoryId(),
    type: 'expense' as 'income' | 'expense',
    date: new Date().toISOString().split('T')[0],
  })

  useEffect(() => {
    if (categories.length > 0 && formData.category === 0) {
      setFormData(prev => ({ ...prev, category: categories[0].id }))
    }
  }, [categories])

  // ===== Вычисляемые данные =====
  const filteredBySearch = useMemo(() => {
    if (!searchQuery.trim()) return transactions
    return transactions.filter(t =>
      t.description.toLowerCase().includes(searchQuery.toLowerCase())
    )
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

  // ===== Статистика =====
  const statsData = useMemo(() => {
    let filtered = [...transactions]
    if (statsDateFrom) filtered = filtered.filter(t => t.date >= statsDateFrom)
    if (statsDateTo) filtered = filtered.filter(t => t.date <= statsDateTo)
    return filtered
  }, [transactions, statsDateFrom, statsDateTo])

  const totalIncome = statsData.filter(t => t.value > 0).reduce((s, t) => s + t.value, 0)
  const totalExpense = statsData.filter(t => t.value < 0).reduce((s, t) => s + t.value, 0)
  const balance = totalIncome + totalExpense

  const getIncomeCategoryStats = () => {
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
  }

  const getExpenseCategoryStats = () => {
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
  }

  const getMonthlyStats = () => {
    const months: Record<string, { income: number; expense: number }> = {}
    statsData.forEach(t => {
      const m = new Date(t.date).toLocaleString('ru-RU', { month: 'short', year: 'numeric' })
      if (!months[m]) months[m] = { income: 0, expense: 0 }
      if (t.value > 0) months[m].income += t.value
      else months[m].expense += Math.abs(t.value)
    })
    return Object.entries(months).sort((a, b) => new Date(a[0]).getTime() - new Date(b[0]).getTime())
  }

  const getCumulativeBalance = () => {
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
    return result
  }

  const getAverageDaily = () => {
    const days = new Set(statsData.map(t => t.date)).size || 1
    const totalInc = statsData.filter(t => t.value > 0).reduce((s, t) => s + t.value, 0)
    const totalExp = statsData.filter(t => t.value < 0).reduce((s, t) => s + t.value, 0)
    return { avgIncome: totalInc / days, avgExpense: Math.abs(totalExp) / days }
  }

  const getWeekdayStats = () => {
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
  }

  const incomeStats = getIncomeCategoryStats()
  const expenseStats = getExpenseCategoryStats()
  const monthlyStats = getMonthlyStats()
  const cumulative = getCumulativeBalance()
  const avgDaily = getAverageDaily()
  const weekdayStats = getWeekdayStats()

  const maxIncome = incomeStats.length ? incomeStats[0][1] : 1
  const maxExpense = expenseStats.length ? expenseStats[0][1] : 1
  const maxMonthly = monthlyStats.reduce((max, [_, { income, expense }]) => Math.max(max, income, expense), 1)
  const maxWeekday = weekdayStats.reduce((max, { value }) => Math.max(max, value), 1)

  // ===== Обработчики =====
  const handleAddTransaction = async (e: React.FormEvent) => {
    e.preventDefault()
    const amount = parseFloat(formData.amount)
    if (!formData.description || isNaN(amount) || amount <= 0) {
      alert('Заполните поля корректно')
      return
    }
    if (!formData.category) {
      alert('Выберите категорию')
      return
    }
    const payload = {
      description: formData.description,
      amount: formData.type === 'income' ? amount : -amount,
      category: formData.category,
      date: formData.date,
    }
    try {
      await addTransactionApi(payload)
      setFormData({
        description: '',
        amount: '',
        category: categories.length > 0 ? categories[0].id : 0,
        type: 'expense',
        date: new Date().toISOString().split('T')[0],
      })
      setShowAddModal(false)
    } catch (err: any) {
      alert(err.message)
    }
  }

  const handleEditTransaction = (transaction: Transaction) => {
    setEditingTransaction(transaction)
    setFormData({
      description: transaction.description,
      amount: Math.abs(transaction.value).toString(),
      category: transaction.category_id,
      type: transaction.value >= 0 ? 'income' : 'expense',
      date: transaction.date,
    })
    setShowAddModal(true)
  }

  const handleUpdateTransaction = async (e: React.FormEvent) => {
    e.preventDefault()
    const amount = parseFloat(formData.amount)
    if (!formData.description || isNaN(amount) || amount <= 0) {
      alert('Заполните поля корректно')
      return
    }
    if (!formData.category) {
      alert('Выберите категорию')
      return
    }
    if (!editingTransaction) return
    const payload = {
      description: formData.description,
      amount: formData.type === 'income' ? amount : -amount,
      category: formData.category,
      date: formData.date,
    }
    try {
      await updateTransactionApi(editingTransaction.id, payload)
      setEditingTransaction(null)
      setFormData({
        description: '',
        amount: '',
        category: categories.length > 0 ? categories[0].id : 0,
        type: 'expense',
        date: new Date().toISOString().split('T')[0],
      })
      setShowAddModal(false)
    } catch (err: any) {
      alert(err.message)
    }
  }

  const handleDeleteTransaction = async (id: number) => {
    if (window.confirm('Удалить транзакцию?')) {
      try {
        await deleteTransactionApi(id)
      } catch (err: any) {
        alert(err.message)
      }
    }
  }

  const handleCreateBudget = async (payload: { category_id: number; limit: number }) => {
    await addBudget(payload)
  }

  // ===== Рендер =====
  if (!isAuthenticated) {
    return <Login onLogin={handleLogin} />
  }

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
            date: new Date().toISOString().split('T')[0],
          })
          setShowAddModal(true)
        }}
        onBudgetClick={() => setShowBudgetModal(true)}
        onLogout={handleLogout}
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
          <Budgets budgets={budgets} loading={budgetsLoading} error={budgetsError} />
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
          budgetStats={budgets.map(b => ({ ...b, spent: b.current_value, percent: Math.min((b.current_value / b.limit_value) * 100, 100) }))}
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
    </div>
  )
}

export default App