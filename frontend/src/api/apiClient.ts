const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

let isRefreshing = false
let failedQueue: Array<{
  resolve: (value: any) => void
  reject: (reason?: any) => void
}> = []

const processQueue = (error: Error | null, token: string | null = null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error)
    } else {
      resolve(token)
    }
  })
  failedQueue = []
}

export const apiClient = async (endpoint: string, options: RequestInit = {}, retry = true): Promise<Response> => {
  const accessToken = localStorage.getItem('access_token')

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  if (accessToken) {
    ;(headers as Record<string, string>)['Authorization'] = accessToken
  }

  let response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    credentials: 'include',
    headers,
  })

  // Если 401 и у нас есть токен (значит мы авторизованы) и это не запрос на обновление
  const isRefreshEndpoint = endpoint.includes('/auth/refresh')

  if (response.status === 401 && retry && !isRefreshEndpoint && accessToken) {
    console.log(' 401 Unauthorized, attempting token refresh...')

    if (isRefreshing) {
      console.log('⏳ Refresh already in progress, waiting...')
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject })
      }).then(() => {
        return apiClient(endpoint, options, false)
      }).catch((err) => {
        throw err
      })
    }

    isRefreshing = true

    try {
      console.log(' Calling /auth/refresh...')
      const refreshResponse = await fetch(`${API_BASE_URL}/auth/refresh`, {
        method: 'POST',
        credentials: 'include',
      })

      if (!refreshResponse.ok) {
        const text = await refreshResponse.text()
        console.error(' Refresh failed:', refreshResponse.status, text)
        throw new Error(`Refresh failed with status ${refreshResponse.status}`)
      }

      const data = await refreshResponse.json()
      const newAccessToken = data.result

      if (!newAccessToken) {
        console.error(' No access token in refresh response')
        throw new Error('No access token in refresh response')
      }

      console.log(' Refresh successful, new token obtained')
      localStorage.setItem('access_token', newAccessToken)
      processQueue(null, newAccessToken)
      isRefreshing = false

      // Повторяем исходный запрос с новым токеном
      return apiClient(endpoint, options, false)
    } catch (error) {
      console.error(' Refresh failed, logging out...', error)
      localStorage.removeItem('access_token')
      localStorage.removeItem('username')
      localStorage.removeItem('auth')
      processQueue(error as Error, null)
      isRefreshing = false

      window.location.href = '/'
      throw new Error('Сессия истекла')
    }
  }

  return response
}