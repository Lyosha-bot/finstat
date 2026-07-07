import { useState, useEffect } from 'react'
import './App.css'

interface Transaction {
  id: number
  date: string
  description: string
  amount: number
  category: string
  type: 'income' | 'expense'
}

interface Budget {
  id: number
  category: string
  limit: number
  spent: number
  period: 'monthly' | 'weekly' | 'yearly'
  startDate: string
  endDate: string
}

function App() {
  const [transactions, setTransactions] = useState<Transaction[]>(() => {
    const saved = localStorage.getItem('transactions')
    return saved ? JSON.parse(saved) : []
  })

  const today = new Date()
  const firstDay = new Date(today.getFullYear(), today.getMonth(), 1).toISOString().split('T')[0]
  const lastDay = new Date(today.getFullYear(), today.getMonth() + 1, 0).toISOString().split('T')[0]

  const [budgets, setBudgets] = useState<Budget[]>(() => {
    const saved = localStorage.getItem('budgets')
    return saved ? JSON.parse(saved) : [
      { id: 1, category: 'Продукты', limit: 15000, spent: 0, period: 'monthly', startDate: firstDay, endDate: lastDay },
      { id: 2, category: 'Транспорт', limit: 5000, spent: 0, period: 'monthly', startDate: firstDay, endDate: lastDay },
    ]
  })

  const [activeTab, setActiveTab] = useState<'transactions' | 'stats'>('transactions')
  const [showAddModal, setShowAddModal] = useState(false)
  const [showBudgetModal, setShowBudgetModal] = useState(false)
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null)

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
    startDate: firstDay,
    endDate: lastDay,
  })

  const [filterType, setFilterType] = useState<'all' | 'income' | 'expense'>('all')
  const [filterCategory, setFilterCategory] = useState('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [periodFilter, setPeriodFilter] = useState<'all' | 'today' | 'week' | 'month'>('all')
  const [statsDateFrom, setStatsDateFrom] = useState('')
  const [statsDateTo, setStatsDateTo] = useState('')

  const getFilteredForStats = () => {
    let filtered = [...transactions]
    if (statsDateFrom) filtered = filtered.filter(t => t.date >= statsDateFrom)
    if (statsDateTo) filtered = filtered.filter(t => t.date <= statsDateTo)
    return filtered
  }

  useEffect(() => {
    localStorage.setItem('transactions', JSON.stringify(transactions))
  }, [transactions])

  useEffect(() => {
    localStorage.setItem('budgets', JSON.stringify(budgets))
  }, [budgets])

  const categories = ['Все', ...new Set(transactions.map(t => t.category))]

  const getFilteredTransactions = () => {
    let filtered = [...transactions]
    if (filterType !== 'all') filtered = filtered.filter(t => t.type === filterType)
    if (filterCategory !== 'all') filtered = filtered.filter(t => t.category === filterCategory)
    if (searchQuery) filtered = filtered.filter(t => t.description.toLowerCase().includes(searchQuery.toLowerCase()))

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

  const statsTransactions = getFilteredForStats()

  const totalIncome = transactions.filter(t => t.type === 'income').reduce((s, t) => s + t.amount, 0)
  const totalExpense = transactions.filter(t => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
  const balance = totalIncome - totalExpense

  // ---- Статистические функции ----
  const getIncomeCategoryStats = () => {
    const map: Record<string, number> = {}
    statsTransactions.filter(t => t.type === 'income').forEach(t => {
      map[t.category] = (map[t.category] || 0) + t.amount
    })
    return Object.entries(map).sort((a, b) => b[1] - a[1])
  }

  const getExpenseCategoryStats = () => {
    const map: Record<string, number> = {}
    statsTransactions.filter(t => t.type === 'expense').forEach(t => {
      map[t.category] = (map[t.category] || 0) + t.amount
    })
    return Object.entries(map).sort((a, b) => b[1] - a[1])
  }

  const getBudgetStats = () => {
    return budgets.map(budget => {
      const spent = statsTransactions
        .filter(t => t.category === budget.category && t.type === 'expense')
        .reduce((s, t) => s + t.amount, 0)
      return { ...budget, spent, percent: Math.min((spent / budget.limit) * 100, 100) }
    })
  }

  const getMonthlyStats = () => {
    const months: Record<string, { income: number; expense: number }> = {}
    statsTransactions.forEach(t => {
      const m = new Date(t.date).toLocaleString('ru-RU', { month: 'short', year: 'numeric' })
      if (!months[m]) months[m] = { income: 0, expense: 0 }
      if (t.type === 'income') months[m].income += t.amount
      else months[m].expense += t.amount
    })
    return Object.entries(months).sort((a, b) => new Date(a[0]).getTime() - new Date(b[0]).getTime())
  }

  // ---- ИСПРАВЛЕННАЯ ФУНКЦИЯ НАКОПЛЕННОГО БАЛАНСА (группировка по дням) ----
  const getCumulativeBalance = () => {
    // Группируем по дате (суммируем за день)
    const dailyMap: Record<string, number> = {}
    statsTransactions.forEach(t => {
      const date = t.date
      const amount = t.type === 'income' ? t.amount : -t.amount
      dailyMap[date] = (dailyMap[date] || 0) + amount
    })
    // Сортируем даты по возрастанию
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
    const days = new Set(statsTransactions.map(t => t.date)).size || 1
    const totalInc = statsTransactions.filter(t => t.type === 'income').reduce((s, t) => s + t.amount, 0)
    const totalExp = statsTransactions.filter(t => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
    return { avgIncome: totalInc / days, avgExpense: totalExp / days }
  }

  const getWeekdayStats = () => {
    const weekdays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
    const map: Record<string, number> = { Пн: 0, Вт: 0, Ср: 0, Чт: 0, Пт: 0, Сб: 0, Вс: 0 }
    statsTransactions.filter(t => t.type === 'expense').forEach(t => {
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
  const cumulative = getCumulativeBalance()  // теперь одна точка на день
  const avgDaily = getAverageDaily()
  const weekdayStats = getWeekdayStats()

  const maxIncome = incomeStats.length ? incomeStats[0][1] : 1
  const maxExpense = expenseStats.length ? expenseStats[0][1] : 1
  const maxMonthly = monthlyStats.reduce((max, [_, { income, expense }]) => Math.max(max, income, expense), 1)
  const maxWeekday = weekdayStats.reduce((max, { value }) => Math.max(max, value), 1)

  // ---- Обработчики ----
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
    setTransactions([...transactions, newTransaction])
    setFormData({ description: '', amount: '', category: 'Продукты', type: 'expense', date: new Date().toISOString().split('T')[0] })
    setShowAddModal(false)
  }

  const handleDeleteTransaction = (id: number) => {
    if (window.confirm('Удалить транзакцию?')) {
      setTransactions(transactions.filter(t => t.id !== id))
    }
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
    const updated: Transaction = { ...editingTransaction, date: formData.date, description: formData.description, amount, category: formData.category, type: formData.type }
    setTransactions(transactions.map(t => (t.id === editingTransaction.id ? updated : t)))
    setEditingTransaction(null)
    setFormData({ description: '', amount: '', category: 'Продукты', type: 'expense', date: new Date().toISOString().split('T')[0] })
    setShowAddModal(false)
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
    const newBudget: Budget = {
      id: Date.now(),
      category: budgetForm.category,
      limit,
      spent: 0,
      period: budgetForm.period,
      startDate: budgetForm.startDate,
      endDate: budgetForm.endDate,
    }
    setBudgets([...budgets, newBudget])
    setBudgetForm({ category: '', limit: '', period: 'monthly', startDate: firstDay, endDate: lastDay })
    setShowBudgetModal(false)
  }

  const handleDeleteBudget = (id: number) => {
    if (window.confirm('Удалить бюджет?')) {
      setBudgets(budgets.filter(b => b.id !== id))
    }
  }

  const getCategorySpent = (category: string) => {
    return transactions.filter(t => t.category === category && t.type === 'expense').reduce((s, t) => s + t.amount, 0)
  }

  const formatMoney = (amount: number) => {
    return new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(amount)
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
  }

  const formatDateRange = (start: string, end: string) => `${formatDate(start)} — ${formatDate(end)}`

  // ---- Компоненты графиков ----
  const LineChart = ({ data, labels, colors, max }: { data: number[][]; labels: string[]; colors: string[]; max: number }) => {
    const padding = { top: 30, bottom: 30, left: 60, right: 20 }
    const width = 600
    const height = 250
    const innerWidth = width - padding.left - padding.right
    const innerHeight = height - padding.top - padding.bottom

    const steps = 5
    const stepValue = Math.ceil(max / steps)
    const maxScaled = stepValue * steps

    const points = data.map((dataset) =>
      dataset.map((value, i) => ({
        x: padding.left + (i / (labels.length - 1 || 1)) * innerWidth,
        y: padding.top + innerHeight - (value / maxScaled) * innerHeight,
      }))
    )

    const pathD = points.map((p) =>
      p.map((pt, i) => (i === 0 ? `M ${pt.x} ${pt.y}` : `L ${pt.x} ${pt.y}`)).join(' ')
    )

    return (
      <svg viewBox={`0 0 ${width} ${height}`} style={{ width: '100%', height: 'auto' }}>
        {Array.from({ length: steps + 1 }, (_, i) => {
          const y = padding.top + innerHeight - (i / steps) * innerHeight
          const val = (i / steps) * maxScaled
          return (
            <g key={i}>
              <line x1={padding.left} y1={y} x2={width - padding.right} y2={y} stroke="#444" strokeWidth="1" strokeDasharray="4" />
              <text x={padding.left - 8} y={y + 4} textAnchor="end" fontSize="10" fill="#888">
                {formatMoney(val)}
              </text>
            </g>
          )
        })}
        {labels.map((label, i) => {
          const x = padding.left + (i / (labels.length - 1 || 1)) * innerWidth
          return (
            <text key={i} x={x} y={height - 4} textAnchor="middle" fontSize="10" fill="#888">
              {label}
            </text>
          )
        })}
        {pathD.map((d, idx) => (
          <path key={idx} d={d} fill="none" stroke={colors[idx]} strokeWidth="2.5" strokeLinejoin="round" />
        ))}
        {points.map((dataset, idx) =>
          dataset.map((pt, i) => (
            <circle key={`${idx}-${i}`} cx={pt.x} cy={pt.y} r="5" fill={colors[idx]}>
              <title>
                {labels[i]}: {formatMoney(data[idx][i])}
              </title>
            </circle>
          ))
        )}
      </svg>
    )
  }

  // ---- НАКОПЛЕННЫЙ БАЛАНС (шкала симметричная) ----
  const CumulativeChart = ({ data }: { data: { label: string; value: number }[] }) => {
    const padding = { top: 30, bottom: 30, left: 60, right: 20 }
    const width = 600
    const height = 180
    const innerWidth = width - padding.left - padding.right
    const innerHeight = height - padding.top - padding.bottom

    const maxAbs = data.reduce((m, d) => Math.max(m, Math.abs(d.value)), 1)
    const magnitude = Math.pow(10, Math.floor(Math.log10(maxAbs)))
    let maxScaled = Math.ceil(maxAbs / magnitude) * magnitude
    if (maxScaled === 0) maxScaled = 1000

    const centerY = padding.top + innerHeight / 2

    const points = data.map((d, i) => ({
      x: padding.left + (i / (data.length - 1 || 1)) * innerWidth,
      y: centerY - (d.value / maxScaled) * (innerHeight / 2),
    }))

    const path = points.map((p, i) => (i === 0 ? `M ${p.x} ${p.y}` : `L ${p.x} ${p.y}`)).join(' ')

    // Шкала от -maxScaled до +maxScaled с шагом maxScaled/4
    const steps = 4
    const stepVal = maxScaled / steps
    const yValues = []
    for (let i = -steps; i <= steps; i++) {
      yValues.push(i * stepVal)
    }

    return (
      <svg viewBox={`0 0 ${width} ${height}`} style={{ width: '100%', height: 'auto' }}>
        {yValues.map((val) => {
          const y = centerY - (val / maxScaled) * (innerHeight / 2)
          if (Math.abs(val) < 0.001) return null
          return (
            <g key={val}>
              <line x1={padding.left} y1={y} x2={width - padding.right} y2={y} stroke="#444" strokeWidth="1" strokeDasharray="4" />
              <text x={padding.left - 8} y={y + 4} textAnchor="end" fontSize="10" fill="#888">
                {formatMoney(val)}
              </text>
            </g>
          )
        })}
        <line x1={padding.left} y1={centerY} x2={width - padding.right} y2={centerY} stroke="#888" strokeWidth="2" strokeDasharray="6" />
        <text x={padding.left - 8} y={centerY + 4} textAnchor="end" fontSize="10" fill="#888">
          {formatMoney(0)}
        </text>

        <path d={path} fill="none" stroke="#4ade80" strokeWidth="2.5" strokeLinejoin="round" />
        {points.map((p, i) => (
          <circle key={i} cx={p.x} cy={p.y} r="5" fill="#4ade80">
            <title>
              {formatDate(data[i].label)}: {formatMoney(data[i].value)}
            </title>
          </circle>
        ))}
        {data.map((d, i) => {
          const x = padding.left + (i / (data.length - 1 || 1)) * innerWidth
          return (
            <text key={i} x={x} y={height - 4} textAnchor="middle" fontSize="9" fill="#888">
              {formatDate(d.label)}
            </text>
          )
        })}
      </svg>
    )
  }

  // ---- JSX ----
  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <h1>Финансовый учёт</h1>
          <div className="header-actions">
            <button className={`tab-btn ${activeTab === 'transactions' ? 'active' : ''}`} onClick={() => setActiveTab('transactions')}>
              Транзакции
            </button>
            <button className={`tab-btn ${activeTab === 'stats' ? 'active' : ''}`} onClick={() => setActiveTab('stats')}>
              Инфографика
            </button>
            <button className="btn btn-primary" onClick={() => { setEditingTransaction(null); setFormData({ description: '', amount: '', category: 'Продукты', type: 'expense', date: new Date().toISOString().split('T')[0] }); setShowAddModal(true) }}>
              + Добавить
            </button>
            <button className="btn btn-secondary" onClick={() => setShowBudgetModal(true)}>Бюджет</button>
          </div>
        </div>
      </header>

      <section className="stats">
        <div className="stat-card"><div className="stat-label">Баланс</div><div className="stat-value balance">{formatMoney(balance)}</div></div>
        <div className="stat-card"><div className="stat-label">Доходы</div><div className="stat-value income">{formatMoney(totalIncome)}</div></div>
        <div className="stat-card"><div className="stat-label">Расходы</div><div className="stat-value expense">{formatMoney(totalExpense)}</div></div>
      </section>

      {activeTab === 'transactions' && (
        <section className="filters">
          <div className="filter-group">
            <label>Период</label>
            <select value={periodFilter} onChange={(e) => setPeriodFilter(e.target.value as any)}>
              <option value="all">Все</option><option value="today">Сегодня</option><option value="week">Неделя</option><option value="month">Месяц</option>
            </select>
          </div>
          <div className="filter-group">
            <label>Тип</label>
            <select value={filterType} onChange={(e) => setFilterType(e.target.value as any)}>
              <option value="all">Все</option><option value="income">Доходы</option><option value="expense">Расходы</option>
            </select>
          </div>
          <div className="filter-group">
            <label>Категория</label>
            <select value={filterCategory} onChange={(e) => setFilterCategory(e.target.value)}>
              <option value="all">Все</option>
              {categories.filter(c => c !== 'Все').map(cat => <option key={cat} value={cat}>{cat}</option>)}
            </select>
          </div>
          <div className="filter-group">
            <label>Поиск</label>
            <input type="text" placeholder="Поиск по описанию..." value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} />
          </div>
        </section>
      )}

      {activeTab === 'transactions' ? (
        <>
          {budgets.length > 0 && (
            <section className="budgets">
              <h2>Бюджеты</h2>
              <div className="budget-grid">
                {budgets.map(budget => {
                  const spent = getCategorySpent(budget.category)
                  const percent = Math.min((spent / budget.limit) * 100, 100)
                  const color = percent > 90 ? '#ef4444' : percent > 70 ? '#f59e0b' : '#4ade80'
                  return (
                    <div key={budget.id} className="budget-card">
                      <div className="budget-header">
                        <span className="budget-category">{budget.category}</span>
                        <button className="btn-delete-small" onClick={() => handleDeleteBudget(budget.id)}>✕</button>
                      </div>
                      <div className="budget-dates">{formatDateRange(budget.startDate, budget.endDate)}</div>
                      <div className="budget-amounts"><span>{formatMoney(spent)}</span><span>/ {formatMoney(budget.limit)}</span></div>
                      <div className="budget-bar"><div className="budget-bar-fill" style={{ width: `${percent}%`, backgroundColor: color }} /></div>
                      <div className="budget-percentage">{Math.round(percent)}%</div>
                    </div>
                  )
                })}
              </div>
            </section>
          )}

          <section className="transactions">
            <h2>Транзакции</h2>
            {filteredTransactions.length === 0 ? (
              <div className="empty-state"><p>Нет транзакций</p><p className="empty-sub">Добавьте первую транзакцию</p></div>
            ) : (
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
                            <div className="transaction-category"><span className="category-badge">{t.category}</span></div>
                          </div>
                          <div className="transaction-amounts">
                            <span className={t.type === 'income' ? 'income' : 'expense'}>
                              {t.type === 'income' ? '+' : '-'}{formatMoney(t.amount)}
                            </span>
                            <div className="transaction-actions">
                              <button className="btn-edit" onClick={() => handleEditTransaction(t)}>✏️</button>
                              <button className="btn-delete" onClick={() => handleDeleteTransaction(t.id)}>🗑️</button>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </section>
        </>
      ) : (
        <section className="stats-dashboard">
          <div className="stats-date-filter">
            <div className="filter-group">
              <label>С даты</label>
              <input type="date" value={statsDateFrom} onChange={(e) => setStatsDateFrom(e.target.value)} />
            </div>
            <div className="filter-group">
              <label>По дату</label>
              <input type="date" value={statsDateTo} onChange={(e) => setStatsDateTo(e.target.value)} />
            </div>
            <button className="btn btn-secondary" onClick={() => { setStatsDateFrom(''); setStatsDateTo(''); }}>Сбросить</button>
          </div>

          <div className="stats-dashboard-grid">
            <div className="chart-card">
              <h3>Доходы по категориям</h3>
              <div className="bar-chart">
                {incomeStats.length === 0 ? <p className="empty-chart">Нет данных</p> :
                  incomeStats.map(([cat, amt]) => (
                    <div key={cat} className="bar-item">
                      <span className="bar-label">{cat}</span>
                      <div className="bar-track"><div className="bar-fill income-fill" style={{ width: `${(amt / maxIncome) * 100}%` }} /></div>
                      <span className="bar-value">{formatMoney(amt)}</span>
                    </div>
                  ))}
              </div>
            </div>

            <div className="chart-card">
              <h3>Расходы по категориям</h3>
              <div className="bar-chart">
                {expenseStats.length === 0 ? <p className="empty-chart">Нет данных</p> :
                  expenseStats.map(([cat, amt]) => (
                    <div key={cat} className="bar-item">
                      <span className="bar-label">{cat}</span>
                      <div className="bar-track"><div className="bar-fill expense-fill" style={{ width: `${(amt / maxExpense) * 100}%` }} /></div>
                      <span className="bar-value">{formatMoney(amt)}</span>
                    </div>
                  ))}
              </div>
            </div>

            <div className="chart-card">
              <h3>Исполнение бюджетов</h3>
              <div className="bar-chart">
                {budgetStats.length === 0 ? <p className="empty-chart">Нет бюджетов</p> :
                  budgetStats.map(b => {
                    const color = b.percent > 90 ? '#ef4444' : b.percent > 70 ? '#f59e0b' : '#4ade80'
                    return (
                      <div key={b.id} className="bar-item">
                        <span className="bar-label">{b.category}</span>
                        <div className="bar-track"><div className="bar-fill" style={{ width: `${b.percent}%`, backgroundColor: color }} /></div>
                        <span className="bar-value">{formatMoney(b.spent)} / {formatMoney(b.limit)}</span>
                      </div>
                    )
                  })}
              </div>
            </div>

            <div className="chart-card">
              <h3>Накопленный баланс</h3>
              <div className="cumulative-chart">
                {cumulative.length === 0 ? <p className="empty-chart">Нет данных</p> : <CumulativeChart data={cumulative} />}
              </div>
            </div>

            <div className="chart-card">
              <h3>Средние значения</h3>
              <div className="avg-stats">
                <div className="avg-item">
                  <span className="avg-label">Средний доход в день</span>
                  <span className="avg-value income">{formatMoney(avgDaily.avgIncome)}</span>
                </div>
                <div className="avg-item">
                  <span className="avg-label">Средний расход в день</span>
                  <span className="avg-value expense">{formatMoney(avgDaily.avgExpense)}</span>
                </div>
                <div className="avg-item">
                  <span className="avg-label">Количество дней в периоде</span>
                  <span className="avg-value">{new Set(statsTransactions.map(t => t.date)).size}</span>
                </div>
              </div>
            </div>

            <div className="chart-card">
              <h3>Расходы по дням недели</h3>
              <div className="bar-chart">
                {weekdayStats.every(w => w.value === 0) ? <p className="empty-chart">Нет данных</p> :
                  weekdayStats.map(({ name, value }) => (
                    <div key={name} className="bar-item">
                      <span className="bar-label">{name}</span>
                      <div className="bar-track"><div className="bar-fill expense-fill" style={{ width: `${(value / maxWeekday) * 100}%` }} /></div>
                      <span className="bar-value">{formatMoney(value)}</span>
                    </div>
                  ))}
              </div>
            </div>

            <div className="chart-card">
              <h3>Динамика доходов и расходов</h3>
              <div className="monthly-balance">
                {monthlyStats.length === 0 ? <p className="empty-chart">Нет данных</p> : (
                  <div className="line-chart-container">
                    <LineChart
                      data={[monthlyStats.map(([_, d]) => d.income), monthlyStats.map(([_, d]) => d.expense)]}
                      labels={monthlyStats.map(([m]) => m)}
                      colors={['#4ade80', '#f87171']}
                      max={maxMonthly}
                    />
                    <div className="line-legend">
                      <span><span className="dot income-dot"></span> Доходы</span>
                      <span><span className="dot expense-dot"></span> Расходы</span>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </section>
      )}

      {showAddModal && (
        <div className="modal-overlay" onClick={() => { setShowAddModal(false); setEditingTransaction(null) }}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2>{editingTransaction ? 'Редактировать' : 'Новая транзакция'}</h2>
            <form onSubmit={editingTransaction ? handleUpdateTransaction : handleAddTransaction}>
              <div className="form-group">
                <label>Тип</label>
                <select value={formData.type} onChange={(e) => setFormData({...formData, type: e.target.value as 'income' | 'expense'})}>
                  <option value="expense">Расход</option>
                  <option value="income">Доход</option>
                </select>
              </div>
              <div className="form-group">
                <label>Описание</label>
                <input type="text" value={formData.description} onChange={(e) => setFormData({...formData, description: e.target.value})} placeholder="Например: Продукты" required />
              </div>
              <div className="form-group">
                <label>Сумма (₽)</label>
                <input type="number" step="0.01" min="0.01" value={formData.amount} onChange={(e) => setFormData({...formData, amount: e.target.value})} placeholder="1000" required />
              </div>
              <div className="form-group">
                <label>Категория</label>
                <select value={formData.category} onChange={(e) => setFormData({...formData, category: e.target.value})}>
                  <option value="Продукты">Продукты</option>
                  <option value="Транспорт">Транспорт</option>
                  <option value="Жильё">Жильё</option>
                  <option value="Развлечения">Развлечения</option>
                  <option value="Здоровье">Здоровье</option>
                  <option value="Образование">Образование</option>
                  <option value="Зарплата">Зарплата</option>
                  <option value="Прочее">Прочее</option>
                </select>
              </div>
              <div className="form-group">
                <label>Дата</label>
                <input type="date" value={formData.date} onChange={(e) => setFormData({...formData, date: e.target.value})} required />
              </div>
              <div className="modal-actions">
                <button type="submit" className="btn btn-primary">{editingTransaction ? 'Сохранить' : 'Добавить'}</button>
                <button type="button" className="btn btn-secondary" onClick={() => { setShowAddModal(false); setEditingTransaction(null) }}>Отмена</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showBudgetModal && (
        <div className="modal-overlay" onClick={() => setShowBudgetModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2>Новый бюджет</h2>
            <form onSubmit={handleAddBudget}>
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
                <button type="button" className="btn btn-secondary" onClick={() => setShowBudgetModal(false)}>Отмена</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default App