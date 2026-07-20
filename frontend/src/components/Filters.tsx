import { type Category } from '../api/categories'

interface FiltersProps {
  periodFilter: 'all' | 'today' | 'week' | 'month'
  setPeriodFilter: (val: any) => void
  filterType: 'all' | 'income' | 'expense'
  setFilterType: (val: any) => void
  filterCategory: number | 'all'   // теперь ID категории или 'all'
  setFilterCategory: (val: number | 'all') => void
  searchQuery: string
  setSearchQuery: (val: string) => void
  categories: Category[]
}

export const Filters = ({
  periodFilter,
  setPeriodFilter,
  filterType,
  setFilterType,
  filterCategory,
  setFilterCategory,
  searchQuery,
  setSearchQuery,
  categories,
}: FiltersProps) => {
  return (
    <section className="filters">
      <div className="filter-group">
        <label>Период</label>
        <select value={periodFilter} onChange={(e) => setPeriodFilter(e.target.value)}>
          <option value="all">Все</option>
          <option value="today">Сегодня</option>
          <option value="week">Неделя</option>
          <option value="month">Месяц</option>
        </select>
      </div>
      <div className="filter-group">
        <label>Тип</label>
        <select value={filterType} onChange={(e) => setFilterType(e.target.value)}>
          <option value="all">Все</option>
          <option value="income">Доходы</option>
          <option value="expense">Расходы</option>
        </select>
      </div>
      <div className="filter-group">
        <label>Категория</label>
        <select
          value={filterCategory === 'all' ? 'all' : filterCategory}
          onChange={(e) => {
            const val = e.target.value
            setFilterCategory(val === 'all' ? 'all' : Number(val))
          }}
        >
          <option value="all">Все</option>
          {categories.map(cat => (
            <option key={cat.id} value={cat.id}>{cat.name}</option>
          ))}
        </select>
      </div>
      <div className="filter-group">
        <label>Поиск</label>
        <input
          type="text"
          placeholder="Поиск по описанию..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>
    </section>
  )
}