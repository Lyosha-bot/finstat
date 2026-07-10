interface FiltersProps {
  periodFilter: 'all' | 'today' | 'week' | 'month'
  setPeriodFilter: (val: any) => void
  filterType: 'all' | 'income' | 'expense'
  setFilterType: (val: any) => void
  filterCategory: string
  setFilterCategory: (val: string) => void
  searchQuery: string
  setSearchQuery: (val: string) => void
  categories: string[]
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
        <select value={filterCategory} onChange={(e) => setFilterCategory(e.target.value)}>
          <option value="all">Все</option>
          {categories.filter(c => c !== 'Все').map(cat => (
            <option key={cat} value={cat}>{cat}</option>
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