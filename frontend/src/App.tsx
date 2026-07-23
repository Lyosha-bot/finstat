/// <reference types="vite/client" />
import { useState, useEffect, lazy, Suspense } from 'react'
import './App.css'
import { Login } from './components'

// Lazy loading тяжёлого Dashboard
const Dashboard = lazy(() =>
  import('./components/Dashboard').then(m => ({ default: m.Dashboard }))
)

// Спиннер-заглушка на время загрузки Dashboard
function LoadingFallback() {
  return (
    <div className="app">
      <div className="header" style={{ opacity: 0.5 }}>
        <div className="header-content">
          <h1>Финансовый учёт</h1>
        </div>
      </div>
      <div className="stats">
        {[1, 2, 3].map(i => (
          <div key={i} className="stat-card">
            <div className="skeleton" style={{ height: 16, width: '60%', margin: '0 auto 8px' }} />
            <div className="skeleton" style={{ height: 32, width: '80%', margin: '0 auto' }} />
          </div>
        ))}
      </div>
    </div>
  )
}

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
    return !!localStorage.getItem('access_token')
  })
  const [username, setUsername] = useState(() => {
    return localStorage.getItem('username') || 'Пользователь'
  })

  useEffect(() => {
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
    localStorage.setItem('auth', 'true')
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

  return (
    <Suspense fallback={<LoadingFallback />}>
      <Dashboard username={username} onLogout={handleLogout} />
    </Suspense>
  )
}

export default App