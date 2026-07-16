/// <reference types="vite/client" />
import { useState } from 'react'
import './App.css'
import { Login } from './components'
import { Dashboard } from './components/Dashboard'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    return localStorage.getItem('auth') === 'true'
  })
  const [username, setUsername] = useState(() => {
    return localStorage.getItem('username') || 'Пользователь'
  })

  const handleLogin = (name: string) => {
    setIsAuthenticated(true)
    setUsername(name)
  }

  const handleLogout = () => {
    localStorage.removeItem('auth')
    localStorage.removeItem('username')
    setIsAuthenticated(false)
    window.location.href = '/'
  }

  if (!isAuthenticated) {
    return <Login onLogin={handleLogin} />
  }

  return <Dashboard username={username} onLogout={handleLogout} />
}

export default App