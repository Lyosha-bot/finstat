const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

export const apiClient = async (endpoint: string, options: RequestInit = {}) => {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })

  if (response.status === 401) {
    localStorage.removeItem('auth')
    localStorage.removeItem('username')
    if (!window.location.pathname.includes('/')) {
      window.location.href = '/'
    }
    throw new Error('Сессия истекла')
  }

  return response
}