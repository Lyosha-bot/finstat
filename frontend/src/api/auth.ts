import { apiClient } from './apiClient'

export const login = async (username: string, password: string): Promise<{ result: string }> => {
  const response = await apiClient('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  }, false) // не пытаемся обновлять токен при логине

  if (!response.ok) {
    let errorMessage = 'Ошибка авторизации'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }

  const data = await response.json()
  // data.result содержит access токен
  if (!data.result) {
    throw new Error('No access token received')
  }
  return data
}

export const register = async (username: string, password: string): Promise<void> => {
  const response = await apiClient('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  }, false)

  if (!response.ok) {
    let errorMessage = 'Ошибка регистрации'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
}

export const checkRegisterValidity = async (username: string, password: string): Promise<void> => {
  const response = await apiClient('/auth/register/is-valid', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  }, false)

  if (!response.ok) {
    let errorMessage = 'Ошибка проверки'
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch (_) {}
    throw new Error(errorMessage)
  }
}

export const logout = async (): Promise<void> => {
  try {
    await apiClient('/auth/logout', { method: 'POST' }, false)
  } catch (_) {
    // игнорируем ошибки при логауте
  } finally {
    localStorage.removeItem('access_token')
    localStorage.removeItem('username')
    localStorage.removeItem('auth')
  }
}