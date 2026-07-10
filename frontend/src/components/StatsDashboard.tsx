import { LineChart, CumulativeChart } from './charts'
import { formatMoney } from '../utils/format'

interface StatsDashboardProps {
  incomeStats: [string, number][]
  expenseStats: [string, number][]
  budgetStats: any[]
  cumulative: { label: string; value: number }[]
  avgDaily: { avgIncome: number; avgExpense: number }
  weekdayStats: { name: string; value: number }[]
  monthlyStats: [string, { income: number; expense: number; }][]
  maxIncome: number
  maxExpense: number
  maxMonthly: number
  maxWeekday: number
  statsDateFrom: string
  statsDateTo: string
  setStatsDateFrom: (val: string) => void
  setStatsDateTo: (val: string) => void
}

export const StatsDashboard = ({
  incomeStats,
  expenseStats,
  budgetStats,
  cumulative,
  avgDaily,
  weekdayStats,
  monthlyStats,
  maxIncome,
  maxExpense,
  maxMonthly,
  maxWeekday,
  statsDateFrom,
  statsDateTo,
  setStatsDateFrom,
  setStatsDateTo,
}: StatsDashboardProps) => {
  return (
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
        {/* Доходы по категориям */}
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

        {/* Расходы по категориям */}
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

        {/* Исполнение бюджетов */}
        <div className="chart-card">
          <h3>Исполнение бюджетов</h3>
          <div className="bar-chart">
            {budgetStats.length === 0 ? <p className="empty-chart">Нет бюджетов</p> :
              budgetStats.map((b: any) => {
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

        {/* Накопленный баланс */}
        <div className="chart-card">
          <h3>Накопленный баланс</h3>
          <div className="cumulative-chart">
            {cumulative.length === 0 ? <p className="empty-chart">Нет данных</p> : <CumulativeChart data={cumulative} />}
          </div>
        </div>

        {/* Средние значения */}
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
              <span className="avg-value">{/* количество дней вычисляется вне компонента */}</span>
            </div>
          </div>
        </div>

        {/* Расходы по дням недели */}
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

        {/* Динамика доходов и расходов */}
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
  )
}