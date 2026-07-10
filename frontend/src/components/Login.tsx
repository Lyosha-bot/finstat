import { useState } from 'react'

interface LoginProps {
  onLogin: (username: string) => void  // передаём имя пользователя
}

export const Login = ({ onLogin }: LoginProps) => {
  const [isRegister, setIsRegister] = useState(false)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!username.trim() || !password.trim()) {
      setError('Заполните все поля')
      return
    }

    if (isRegister) {
      // Регистрация
      if (password !== confirmPassword) {
        setError('Пароли не совпадают')
        return
      }
      // localStorage
      const users = JSON.parse(localStorage.getItem('users') || '[]')
      if (users.some((u: any) => u.username === username)) {
        setError('Пользователь уже существует')
        return
      }
      users.push({ username, password })
      localStorage.setItem('users', JSON.stringify(users))
      // После регистрации сразу логин
      localStorage.setItem('auth', 'true')
      localStorage.setItem('username', username)
      onLogin(username)
    } else {
      // Вход
      const users = JSON.parse(localStorage.getItem('users') || '[]')
      const user = users.find((u: any) => u.username === username && u.password === password)
      if (!user) {
        setError('Неверный логин или пароль')
        return
      }
      localStorage.setItem('auth', 'true')
      localStorage.setItem('username', username)
      onLogin(username)
    }
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
          <button type="submit" className="btn btn-primary" style={{ width: '100%' }}>
            {isRegister ? 'Зарегистрироваться' : 'Войти'}
          </button>
        </form>
        <div className="login-toggle">
          <button 
            type="button" 
            className="btn btn-secondary" 
            onClick={() => {
              setIsRegister(!isRegister)
              setError('')
            }}
            style={{ width: '100%', marginTop: '0.5rem' }}
          >
            {isRegister ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться'}
          </button>
        </div>
      </div>
    </div>
  )
}