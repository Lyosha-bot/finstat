/// <reference types="vite/client" />
import { useState } from 'react'
import './App.css'

import type { Transaction } from './types'
import { useTransactions, useBudgets } from './hooks'
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
    // Можно также вызвать logout API, если есть
  }

  // ===== Данные =====
  const { transactions, addTransaction, deleteTransaction, updateTransaction } = useTransactions()
  const [budgets, setBudgets] = useBudgets([
    {
      id: 1,
      category: 'Продукты',
      limit: 15000,
      spent: 0,
      period: 'monthly',
      startDate: '2026-07-01',
      endDate: '2026-07-31',
    },
    {
      id: 2,
      category: 'Транспорт',
      limit: 5000,
      spent: 0,
      period: 'monthly',
      startDate: '2026-07-01',
      endDate: '2026-07-31',
    },
  ])

  // ===== UI состояния =====
  const [activeTab, setActiveTab] = useState<'transactions' | 'stats'>('transactions')
  const [showAddModal, setShowAddModal] = useState(false)
  const [showBudgetModal, setShowBudgetModal] = useState(false)
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null)

  // ===== Формы =====
  const [formData, setFormData] = useState({
    description: '',
    amount: '',
    category: 'Продукты',
    type: 'expense' as 'income' | 'expense',
    date: new Date().toISOString().split('T')[0],
  })

  const [budgetForm, setBudgetForm] = useState({
    category: '',
    limit: '',
    period: 'monthly' as 'monthly' | 'weekly' | 'yearly',
    startDate: new Date(new Date().getFullYear(), new Date().getMonth(), 1)
      .toISOString()
      .split('T')[0],
    endDate: new Date(new Date().getFullYear(), new Date().getMonth() + 1, 0)
      .toISOString()
      .split('T')[0],
  })

  // ===== Фильтры =====
  const [filterType, setFilterType] = useState<'all' | 'income' | 'expense'>('all')
  const [filterCategory, setFilterCategory] = useState('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [periodFilter, setPeriodFilter] = useState<'all' | 'today' | 'week' | 'month'>('all')
  const [statsDateFrom, setStatsDateFrom] = useState('')
  const [statsDateTo, setStatsDateTo] = useState('')

  // ===== Вычисляемые данные =====
  const categories = ['Все', ...new Set(transactions.map(t => t.category))]

  const getFilteredTransactions = () => {
    let filtered = [...transactions]
    if (filterType !== 'all') filtered = filtered.filter(t => t.type === filterType)
    if (filterCategory !== 'all') filtered = filtered.filter(t => t.category === filterCategory)
    if (searchQuery) {
      filtered = filtered.filter(t =>
        t.description.toLowerCase().includes(searchQuery.toLowerCase())
      )
    }
    const now = new Date()
    if (periodFilter === 'today') {
      filtered = filtered.filter(t => new Date(t.date).toDateString() === now.toDateString())
    } else if (periodFilter === 'week') {
      const weekAgo = new Date(now)
      weekAgo.setDate(now.getDate() - 7)
      filtered = filtered.filter(t => new Date(t.date) >= weekAgo)
    } else if (periodFilter === 'month') {
      const monthAgo = new Date(now)
      monthAgo.setMonth(now.getMonth() - 1)
      filtered = filtered.filter(t => new Date(t.date) >= monthAgo)
    }
    filtered.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
    return filtered
  }

  const filteredTransactions = getFilteredTransactions()

  const groupedTransactions = () => {
    const groups: { [date: string]: Transaction[] } = {}
    filteredTransactions.forEach(t => {
      if (!groups[t.date]) groups[t.date] = []
      groups[t.date].push(t)
    })
    return Object.entries(groups)
      .sort(([a], [b]) => new Date(b).getTime() - new Date(a).getTime())
      .map(([date, items]) => ({
        date,
        total: items.reduce((sum, t) => sum + (t.type === 'income' ? t.amount : -t.amount), 0),
        transactions: items,
      }))
  }

  const grouped = groupedTransactions()

  const statsData = (() => {
    let filtered = [...transactions]
    if (statsDateFrom) filtered = filtered.filter(t => t.date >= statsDateFrom)
    if (statsDateTo) filtered = filtered.filter(t => t.date <= statsDateTo)
    return filtered
  })()

  const totalIncome = transactions.filter(t => t.type === 'income').reduce((s, t) => s + t.amount, 0)
  const totalExpense = transactions.filter(t => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
  const balance = totalIncome - totalExpense

  // ---- Статистические функции ----
  const getIncomeCategoryStats = () => {
    const map: Record<string, number> = {}
    statsData.filter(t => t.type === 'income').forEach(t => {
      map[t.category] = (map[t.category] || 0) + t.amount
    })
    return Object.entries(map).sort((a, b) => b[1] - a[1])
  }

  const getExpenseCategoryStats = () => {
    const map: Record<string, number> = {}
    statsData.filter(t => t.type === 'expense').forEach(t => {
      map[t.category] = (map[t.category] || 0) + t.amount
    })
    return Object.entries(map).sort((a, b) => b[1] - a[1])
  }

  const getBudgetStats = () => {
    return budgets.map(budget => {
      const spent = statsData
        .filter(t => t.category === budget.category && t.type === 'expense')
        .reduce((s, t) => s + t.amount, 0)
      return { ...budget, spent, percent: Math.min((spent / budget.limit) * 100, 100) }
    })
  }

  const getMonthlyStats = () => {
    const months: Record<string, { income: number; expense: number }> = {}
    statsData.forEach(t => {
      const m = new Date(t.date).toLocaleString('ru-RU', { month: 'short', year: 'numeric' })
      if (!months[m]) months[m] = { income: 0, expense: 0 }
      if (t.type === 'income') months[m].income += t.amount
      else months[m].expense += t.amount
    })
    return Object.entries(months).sort((a, b) => new Date(a[0]).getTime() - new Date(b[0]).getTime())
  }

  const getCumulativeBalance = () => {
    const dailyMap: Record<string, number> = {}
    statsData.forEach(t => {
      const date = t.date
      const amount = t.type === 'income' ? t.amount : -t.amount
      dailyMap[date] = (dailyMap[date] || 0) + amount
    })
    const sortedDates = Object.keys(dailyMap).sort(
      (a, b) => new Date(a).getTime() - new Date(b).getTime()
    )
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
    const totalInc = statsData.filter(t => t.type === 'income').reduce((s, t) => s + t.amount, 0)
    const totalExp = statsData.filter(t => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
    return { avgIncome: totalInc / days, avgExpense: totalExp / days }
  }

  const getWeekdayStats = () => {
    const weekdays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
    const map: Record<string, number> = { Пн: 0, Вт: 0, Ср: 0, Чт: 0, Пт: 0, Сб: 0, Вс: 0 }
    statsData.filter(t => t.type === 'expense').forEach(t => {
      const d = new Date(t.date)
      let wd = d.getDay()
      if (wd === 0) wd = 7
      const name = weekdays[wd - 1]
      map[name] = (map[name] || 0) + t.amount
    })
    return Object.entries(map).map(([name, value]) => ({ name, value }))
  }

  const incomeStats = getIncomeCategoryStats()
  const expenseStats = getExpenseCategoryStats()
  const budgetStats = getBudgetStats()
  const monthlyStats = getMonthlyStats()
  const cumulative = getCumulativeBalance()
  const avgDaily = getAverageDaily()
  const weekdayStats = getWeekdayStats()

  const maxIncome = incomeStats.length ? incomeStats[0][1] : 1
  const maxExpense = expenseStats.length ? expenseStats[0][1] : 1
  const maxMonthly = monthlyStats.reduce((max, [_, { income, expense }]) => Math.max(max, income, expense), 1)
  const maxWeekday = weekdayStats.reduce((max, { value }) => Math.max(max, value), 1)

  // ===== Обработчики =====
  const handleAddTransaction = (e: React.FormEvent) => {
    e.preventDefault()
    const amount = parseFloat(formData.amount)
    if (!formData.description || isNaN(amount) || amount <= 0) {
      alert('Заполните поля корректно')
      return
    }
    const newTransaction: Transaction = {
      id: Date.now(),
      date: formData.date,
      description: formData.description,
      amount,
      category: formData.category,
      type: formData.type,
    }
    addTransaction(newTransaction)
    setFormData({
      description: '',
      amount: '',
      category: 'Продукты',
      type: 'expense',
      date: new Date().toISOString().split('T')[0],
    })
    setShowAddModal(false)
  }

  const handleEditTransaction = (transaction: Transaction) => {
    setEditingTransaction(transaction)
    setFormData({
      description: transaction.description,
      amount: transaction.amount.toString(),
      category: transaction.category,
      type: transaction.type,
      date: transaction.date,
    })
    setShowAddModal(true)
  }

  const handleUpdateTransaction = (e: React.FormEvent) => {
    e.preventDefault()
    const amount = parseFloat(formData.amount)
    if (!formData.description || isNaN(amount) || amount <= 0) {
      alert('Заполните поля корректно')
      return
    }
    if (!editingTransaction) return
    const updated: Transaction = {
      ...editingTransaction,
      date: formData.date,
      description: formData.description,
      amount,
      category: formData.category,
      type: formData.type,
    }
    updateTransaction(updated)
    setEditingTransaction(null)
    setFormData({
      description: '',
      amount: '',
      category: 'Продукты',
      type: 'expense',
      date: new Date().toISOString().split('T')[0],
    })
    setShowAddModal(false)
  }

  const handleDeleteTransaction = (id: number) => {
    if (window.confirm('Удалить транзакцию?')) {
      deleteTransaction(id)
    }
  }

  const handleAddBudget = (e: React.FormEvent) => {
    e.preventDefault()
    const limit = parseFloat(budgetForm.limit)
    if (!budgetForm.category || isNaN(limit) || limit <= 0) {
      alert('Заполните поля корректно')
      return
    }
    if (budgets.some(b => b.category === budgetForm.category)) {
      alert('Бюджет для этой категории уже существует')
      return
    }
    const newBudget = {
      id: Date.now(),
      category: budgetForm.category,
      limit,
      spent: 0,
      period: budgetForm.period,
      startDate: budgetForm.startDate,
      endDate: budgetForm.endDate,
    }
    setBudgets([...budgets, newBudget])
    setBudgetForm({
      category: '',
      limit: '',
      period: 'monthly',
      startDate: new Date(new Date().getFullYear(), new Date().getMonth(), 1)
        .toISOString()
        .split('T')[0],
      endDate: new Date(new Date().getFullYear(), new Date().getMonth() + 1, 0)
        .toISOString()
        .split('T')[0],
    })
    setShowBudgetModal(false)
  }

  const handleDeleteBudget = (id: number) => {
    if (window.confirm('Удалить бюджет?')) {
      setBudgets(budgets.filter(b => b.id !== id))
    }
  }

  const getCategorySpent = (category: string) => {
    return transactions
      .filter(t => t.category === category && t.type === 'expense')
      .reduce((s, t) => s + t.amount, 0)
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
            category: 'Продукты',
            type: 'expense',
            date: new Date().toISOString().split('T')[0],
          })
          setShowAddModal(true)
        }}
        onBudgetClick={() => setShowBudgetModal(true)}
        onLogout={handleLogout}
        username={username}
      />

      <Stats balance={balance} totalIncome={totalIncome} totalExpense={totalExpense} />

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
            onDeleteBudget={handleDeleteBudget}
            getCategorySpent={getCategorySpent}
          />
          <section className="transactions">
            <h2>Транзакции</h2>
            <TransactionList
              grouped={grouped}
              onEdit={handleEditTransaction}
              onDelete={handleDeleteTransaction}
            />
          </section>
        </>
      ) : (
        <StatsDashboard
          incomeStats={incomeStats}
          expenseStats={expenseStats}
          budgetStats={budgetStats}
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
      />

      <BudgetModal
        isOpen={showBudgetModal}
        onClose={() => setShowBudgetModal(false)}
        budgetForm={budgetForm}
        setBudgetForm={setBudgetForm}
        onSubmit={handleAddBudget}
      />
    </div>
  )
}

export default App