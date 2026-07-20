/// <reference types="vite/client" />
import { useState, useEffect } from 'react'
import './App.css'
import { Login } from './components'
import { Dashboard } from './components/Dashboard'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
    return !!localStorage.getItem('access_token')
  })
  const [username, setUsername] = useState(() => {
    return localStorage.getItem('username') || 'Пользователь'
  })

  useEffect(() => {
    // Если access_token нет, но auth есть – удаляем
    if (!localStorage.getItem('access_token') && localStorage.getItem('auth')) {
      localStorage.removeItem('auth')
      localStorage.removeItem('username')
      setIsAuthenticated(false)
      setUsername('Пользователь')
    }
  }, [])

  const handleLogin = (name: string, accessToken: string) => {
    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('username', name)
    localStorage.setItem('auth', 'true') // для совместимости
    setIsAuthenticated(true)
    setUsername(name)
  }

  const handleLogout = async () => {
    const { logout } = await import('./api/auth')
    await logout()
    setIsAuthenticated(false)
    setUsername('Пользователь')
    window.location.href = '/'
  }

  if (!isAuthenticated) {
    return <Login onLogin={handleLogin} />
  }

  return <Dashboard username={username} onLogout={handleLogout} />
}

export default App