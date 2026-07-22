// Кэшированный форматтер
const moneyFormatter = new Intl.NumberFormat('ru-RU', {
  style: 'currency',
  currency: 'RUB',
  minimumFractionDigits: 0,
  maximumFractionDigits: 0,
})

const dateFormatter = new Intl.DateTimeFormat('ru-RU', {
  day: '2-digit',
  month: '2-digit',
  year: 'numeric',
})

export const formatMoney = (amount: number): string => {
  if (isNaN(amount)) return '0 ₽'
  return moneyFormatter.format(amount)
}

export const formatDate = (dateStr: string): string => {
  return dateFormatter.format(new Date(dateStr))
}

export const formatDateRange = (start: string, end: string): string =>
  `${formatDate(start)} — ${formatDate(end)}`

/* Форматирует Date в YYYY-MM-DD в ЛОКАЛЬНОЙ таймзоне */
export const toLocalDateStr = (d: Date): string => {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}
