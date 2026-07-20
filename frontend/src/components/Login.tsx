import { useState } from 'react'
import { login, register, checkRegisterValidity } from '../api/auth'

interface LoginProps {
  onLogin: (username: string, accessToken: string) => void
}

export const Login = ({ onLogin }: LoginProps) => {
  const [isRegister, setIsRegister] = useState(false)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [isChecking, setIsChecking] = useState(false)

  // Проверка на потерю фокуса (onBlur)
  const validateRegistration = async () => {
    if (!isRegister) return
    if (!username.trim() || !password.trim()) {
      setError('')
      return
    }
    if (password !== confirmPassword) {
      setError('Пароли не совпадают')
      return
    }
    setIsChecking(true)
    try {
      await checkRegisterValidity(username, password)
      setError('')
    } catch (err: any) {
      setError(err.message || 'Ошибка проверки')
    } finally {
      setIsChecking(false)
    }
  }

  const handleUsernameBlur = () => {
    if (username.trim()) validateRegistration()
  }

  const handlePasswordBlur = () => {
    if (password.trim()) validateRegistration()
  }

  const handleConfirmPasswordBlur = () => {
    if (confirmPassword.trim()) validateRegistration()
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    if (!username.trim() || !password.trim()) {
      setError('Заполните все поля')
      setLoading(false)
      return
    }

    if (isRegister && password !== confirmPassword) {
      setError('Пароли не совпадают')
      setLoading(false)
      return
    }

    try {
      if (isRegister) {
        // Проверяем ещё раз перед регистрацией
        await checkRegisterValidity(username, password)
        await register(username, password)
        // Автоматический вход
        const data = await login(username, password)
        onLogin(username, data.result)
      } else {
        const data = await login(username, password)
        onLogin(username, data.result)
      }
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
              onBlur={handleUsernameBlur}
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
              onBlur={handlePasswordBlur}
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
                onBlur={handleConfirmPasswordBlur}
                placeholder="Повторите пароль"
                required
              />
            </div>
          )}
          {error && <div className="login-error">{error}</div>}
          {isChecking && <div className="login-info">Проверка данных...</div>}
          <button type="submit" className="btn btn-primary" style={{ width: '100%' }} disabled={loading || isChecking}>
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