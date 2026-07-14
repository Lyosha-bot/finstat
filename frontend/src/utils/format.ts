export const formatMoney = (amount: number): string => {
  if (isNaN(amount)) return '0 ₽' // защита от NaN
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount)
}

export const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}

export const formatDateRange = (start: string, end: string): string =>
  `${formatDate(start)} — ${formatDate(end)}`