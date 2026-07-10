const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

export const apiClient = (endpoint: string, options: RequestInit = {}) => {
  return fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    credentials: 'include', // отправлять cookies (JWT)
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })
}