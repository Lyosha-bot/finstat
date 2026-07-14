import { useState } from 'react'
import { login, register, checkRegisterValidity } from '../api/auth'

interface LoginProps {
  onLogin: (username: string) => void
}

export const Login = ({ onLogin }: LoginProps) => {
  const [isRegister, setIsRegister] = useState(false)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    if (!username.trim() || !password.trim()) {
      setError('Заполните все поля')
      setLoading(false)
      return
    }

    try {
      if (isRegister) {
        if (password !== confirmPassword) {
          setError('Пароли не совпадают')
          setLoading(false)
          return
        }
        // Проверяем возможность регистрации
        await checkRegisterValidity(username, password)
        // Регистрируем
        await register(username, password)
        // Автоматический вход
        await login(username, password)
      } else {
        await login(username, password)
      }

      localStorage.setItem('auth', 'true')
      localStorage.setItem('username', username)
      onLogin(username)
    } catch (err: any) {
      setError(err.message || 'Неизвестная ошибка')
    } finally {
      setLoading(false)
    }
  }

  const toggleMode = () => {
    setIsRegister(!isRegister)
    setError('')
    setPassword('')
    setConfirmPassword('')
  }

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Финансовый учёт</h1>
        <p>{isRegister ? 'Создайте аккаунт' : 'Войдите в систему'}</p>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Логин</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Введите логин"
              required
            />
          </div>
          <div className="form-group">
            <label>Пароль</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Введите пароль"
              required
            />
          </div>
          {isRegister && (
            <div className="form-group">
              <label>Повторите пароль</label>
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Повторите пароль"
                required
              />
            </div>
          )}
          {error && <div className="login-error">{error}</div>}
          <button type="submit" className="btn btn-primary" style={{ width: '100%' }} disabled={loading}>
            {loading ? 'Загрузка...' : isRegister ? 'Зарегистрироваться' : 'Войти'}
          </button>
        </form>
        <div className="login-toggle">
          <button
            type="button"
            className="btn btn-secondary"
            onClick={toggleMode}
            style={{ width: '100%', marginTop: '0.5rem' }}
          >
            {isRegister ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться'}
          </button>
        </div>
      </div>
    </div>
  )
}