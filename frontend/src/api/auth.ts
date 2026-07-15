import { apiClient } from './apiClient'

export const login = async (username: string, password: string): Promise<{ message: string }> => {
  const response = await apiClient('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })

  if (!response.ok) {
    let errorMessage = 'Ошибка авторизации'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }

  return response.json()
}

// Функция регистрации
export const register = async (username: string, password: string): Promise<{ message: string }> => {
  const response = await apiClient('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })

  if (!response.ok) {
    let errorMessage = 'Ошибка регистрации'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }

  return response.json()
}

//проверка
export const checkRegisterValidity = async (username: string, password: string): Promise<{ message: string }> => {
  const response = await apiClient('/auth/register/is-valid', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })
  if (!response.ok) {
    let errorMessage = 'Ошибка проверки данных'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
  return response.json()
}